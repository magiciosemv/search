package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"solana-monitor/internal/models"
	"solana-monitor/internal/services"
)

type Handler struct {
	db       *models.DB
	solana   *services.SolanaService
	monitor  *services.MonitorService
	notifier *services.NotificationService
	bus      *services.EventBus
}

func NewHandler(db *models.DB, solana *services.SolanaService, monitor *services.MonitorService, notifier *services.NotificationService, bus *services.EventBus) *Handler {
	return &Handler{
		db:       db,
		solana:   solana,
		monitor:  monitor,
		notifier: notifier,
		bus:      bus,
	}
}

// respondError sends a structured JSON error response.
func respondError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{"error": message, "code": code})
}

// respondOK sends a structured success response.
func respondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// parseID parses a URL path parameter as a positive int64 ID.
// Returns 0 and sends a 400 error if invalid.
func parseID(c *gin.Context) int64 {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return 0
	}
	return id
}

// solanaAddressRegex matches base58-encoded Solana addresses (32-44 chars).
var solanaAddressRegex = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`)

// isValidSolanaAddress checks if the string looks like a valid Solana address.
func isValidSolanaAddress(addr string) bool {
	return solanaAddressRegex.MatchString(addr)
}

func (h *Handler) GetAddresses(c *gin.Context) {
	addresses, err := h.db.GetAllAddresses()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch addresses")
		return
	}
	respondOK(c, addresses)
}

func (h *Handler) CreateAddress(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
		Label   string `json:"label"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "address is required")
		return
	}

	if !isValidSolanaAddress(req.Address) {
		respondError(c, http.StatusBadRequest, "INVALID_ADDRESS", "invalid Solana address format (base58, 32-44 characters)")
		return
	}

	addr, err := h.db.CreateAddress(req.Address, req.Label)
	if err != nil {
		respondError(c, http.StatusConflict, "DUPLICATE_ADDRESS", "this address already exists")
		return
	}

	// Get initial balance
	balance, err := h.solana.GetBalance(c.Request.Context(), req.Address)
	if err == nil {
		h.db.UpdateAddressBalance(addr.ID, balance)
		addr.Balance = balance
	}

	// Send notification about new wallet
	h.notifier.SendAlert(services.AlertMessage{
		Address:   req.Address,
		Label:     req.Label,
		AlertType: "wallet_added",
		OldValue:  0,
		NewValue:  balance,
		Time:      time.Now(),
	})

	c.JSON(http.StatusCreated, addr)
}

func (h *Handler) GetAddress(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	addr, err := h.db.GetAddressByID(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "address not found")
		return
	}
	respondOK(c, addr)
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	var req struct {
		Label string `json:"label"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if err := h.db.UpdateAddress(id, req.Label); err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update address")
		return
	}
	respondOK(c, gin.H{"success": true})
}

func (h *Handler) DeleteAddress(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	if err := h.db.DeleteAddress(id); err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete address")
		return
	}
	respondOK(c, gin.H{"success": true})
}

func (h *Handler) RefreshAddressBalance(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	addr, err := h.db.GetAddressByID(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "address not found")
		return
	}

	balance, err := h.solana.GetBalance(c.Request.Context(), addr.Address)
	if err != nil {
		respondError(c, http.StatusBadGateway, "RPC_ERROR", "failed to fetch balance from Solana RPC")
		return
	}

	h.db.UpdateAddressBalance(id, balance)
	addr.Balance = balance

	respondOK(c, addr)
}

func (h *Handler) GetRules(c *gin.Context) {
	rules, err := h.db.GetAllRules()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch rules")
		return
	}
	respondOK(c, rules)
}

func (h *Handler) CreateRule(c *gin.Context) {
	var req struct {
		AddressID int64   `json:"address_id" binding:"required"`
		RuleType  string  `json:"rule_type" binding:"required"`
		Threshold float64 `json:"threshold" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "address_id, rule_type, and threshold are required")
		return
	}

	if req.Threshold < 0 {
		respondError(c, http.StatusBadRequest, "INVALID_THRESHOLD", "threshold must be >= 0")
		return
	}

	if req.RuleType != "balance_change" {
		respondError(c, http.StatusBadRequest, "INVALID_RULE_TYPE", "rule_type must be 'balance_change'")
		return
	}

	rule, err := h.db.CreateRule(req.AddressID, req.RuleType, req.Threshold)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create rule")
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	var req struct {
		RuleType  string  `json:"rule_type"`
		Threshold float64 `json:"threshold"`
		Enabled   bool    `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if err := h.db.UpdateRule(id, req.RuleType, req.Threshold, req.Enabled); err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update rule")
		return
	}
	respondOK(c, gin.H{"success": true})
}

func (h *Handler) DeleteRule(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	if err := h.db.DeleteRule(id); err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete rule")
		return
	}
	respondOK(c, gin.H{"success": true})
}

func (h *Handler) GetNotifications(c *gin.Context) {
	notifications, err := h.db.GetAllNotifications()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch notifications")
		return
	}
	respondOK(c, notifications)
}

func (h *Handler) CreateNotification(c *gin.Context) {
	var req struct {
		Name   string            `json:"name" binding:"required"`
		Type   string            `json:"type" binding:"required"`
		Config map[string]string `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "name and type are required")
		return
	}

	if req.Type != "telegram" && req.Type != "email" {
		respondError(c, http.StatusBadRequest, "INVALID_TYPE", "type must be 'telegram' or 'email'")
		return
	}

	if req.Type == "telegram" {
		if req.Config == nil || req.Config["chat_id"] == "" {
			respondError(c, http.StatusBadRequest, "MISSING_CONFIG", "chat_id is required for telegram notifications")
			return
		}
	}

	if req.Type == "email" {
		if req.Config == nil || req.Config["email"] == "" {
			respondError(c, http.StatusBadRequest, "MISSING_CONFIG", "email is required for email notifications")
			return
		}
	}

	configJSON := ""
	if req.Config != nil {
		notif := &models.Notification{}
		notif.SetConfig(req.Config)
		configJSON = notif.Config.String
	}

	notif, err := h.db.CreateNotification(req.Name, req.Type, configJSON)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create notification channel")
		return
	}
	c.JSON(http.StatusCreated, notif)
}

func (h *Handler) UpdateNotification(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	var req struct {
		Name    string            `json:"name"`
		Type    string            `json:"type"`
		Config  map[string]string `json:"config"`
		Enabled bool              `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	configJSON := ""
	if req.Config != nil {
		notif := &models.Notification{}
		notif.SetConfig(req.Config)
		configJSON = notif.Config.String
	}

	if err := h.db.UpdateNotification(id, req.Name, req.Type, configJSON, req.Enabled); err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update notification channel")
		return
	}
	respondOK(c, gin.H{"success": true})
}

func (h *Handler) DeleteNotification(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	if err := h.db.DeleteNotification(id); err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete notification channel")
		return
	}
	respondOK(c, gin.H{"success": true})
}

func (h *Handler) TestNotification(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		return
	}

	notif, err := h.db.GetNotificationByID(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "notification channel not found")
		return
	}

	config := notif.GetConfig()
	msg := "Test notification from Solana Monitor"

	if notif.Type == "telegram" {
		if chatID, ok := config["chat_id"]; ok {
			err = h.notifier.SendTelegram(chatID, msg)
		}
	} else if notif.Type == "email" {
		if to, ok := config["email"]; ok {
			err = h.notifier.SendEmail(to, "Test Alert", msg)
		}
	}

	if err != nil {
		respondError(c, http.StatusBadGateway, "SEND_FAILED", "failed to send test notification: "+err.Error())
		return
	}
	respondOK(c, gin.H{"success": true})
}

func (h *Handler) GetAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit < 1 || limit > 200 {
		respondError(c, http.StatusBadRequest, "INVALID_LIMIT", "limit must be between 1 and 200")
		return
	}

	alerts, err := h.db.GetAllAlerts(limit, offset)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch alerts")
		return
	}
	respondOK(c, alerts)
}

func (h *Handler) GetAlertStats(c *gin.Context) {
	stats, err := h.db.GetAlertStats()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch alert stats")
		return
	}
	respondOK(c, stats)
}

func (h *Handler) Health(c *gin.Context) {
	respondOK(c, gin.H{
		"status":         "ok",
		"monitor_running": true,
	})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.db.GetStats()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch stats")
		return
	}

	alertStats, _ := h.db.GetAlertStats()

	respondOK(c, gin.H{
		"total_addresses":     stats["addresses"],
		"total_rules":         stats["rules"],
		"total_notifications": stats["notifications"],
		"total_alerts":        alertStats["total"],
		"today_alerts":        alertStats["today"],
		"monitor_running":     true,
	})
}

// StreamEvents is a SSE endpoint that pushes real-time balance updates and alerts.
func (h *Handler) StreamEvents(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ch := h.bus.Subscribe()
	defer h.bus.Unsubscribe(ch)

	ctx := c.Request.Context()

	// Send initial connection event
	c.Writer.Write([]byte("event: connected\ndata: {}\n\n"))
	c.Writer.Flush()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(event.Data)
			if err != nil {
				continue
			}
			fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event.Type, data)
			c.Writer.Flush()
		case <-ticker.C:
			// Keepalive
			c.Writer.Write([]byte(": keepalive\n\n"))
			c.Writer.Flush()
		}
	}
}

func (h *Handler) BackupDatabase(c *gin.Context) {
	dbPath := h.db.DBPath()
	if dbPath == "" {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "database path not available")
		return
	}

	backupPath := dbPath + ".backup.tmp"
	_, err := h.db.Exec(fmt.Sprintf("VACUUM INTO '%s'", backupPath))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "backup failed")
		return
	}

	c.Header("Content-Disposition", "attachment; filename=solana-monitor-backup.db")
	c.Header("Content-Type", "application/octet-stream")
	c.File(backupPath)

	// Clean up temp file after sending
	go func() {
		time.Sleep(2 * time.Second)
		os.Remove(backupPath)
	}()
}
