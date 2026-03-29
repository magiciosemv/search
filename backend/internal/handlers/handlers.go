package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"solana-monitor/internal/models"
	"solana-monitor/internal/services"
)

type Handler struct {
	db       *models.DB
	solana   *services.SolanaService
	monitor  *services.MonitorService
	notifier *services.NotificationService
}

func NewHandler(db *models.DB, solana *services.SolanaService, monitor *services.MonitorService, notifier *services.NotificationService) *Handler {
	return &Handler{
		db:       db,
		solana:   solana,
		monitor:  monitor,
		notifier: notifier,
	}
}

func (h *Handler) GetAddresses(c *gin.Context) {
	addresses, err := h.db.GetAllAddresses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, addresses)
}

func (h *Handler) CreateAddress(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
		Label   string `json:"label"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr, err := h.db.CreateAddress(req.Address, req.Label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get initial balance
	balance, err := h.solana.GetBalance(c.Request.Context(), req.Address)
	if err == nil {
		h.db.UpdateAddressBalance(addr.ID, balance)
		addr.Balance = balance
	}

	c.JSON(http.StatusCreated, addr)
}

func (h *Handler) GetAddress(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	addr, err := h.db.GetAddressByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, addr)
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Label string `json:"label"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateAddress(id, req.Label); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteAddress(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.db.DeleteAddress(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) RefreshAddressBalance(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	addr, err := h.db.GetAddressByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	balance, err := h.solana.GetBalance(c.Request.Context(), addr.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.db.UpdateAddressBalance(id, balance)
	addr.Balance = balance

	c.JSON(http.StatusOK, addr)
}

func (h *Handler) GetRules(c *gin.Context) {
	rules, err := h.db.GetAllRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

func (h *Handler) CreateRule(c *gin.Context) {
	var req struct {
		AddressID int64   `json:"address_id" binding:"required"`
		RuleType  string  `json:"rule_type" binding:"required"`
		Threshold float64 `json:"threshold" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.db.CreateRule(req.AddressID, req.RuleType, req.Threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		RuleType  string  `json:"rule_type"`
		Threshold float64 `json:"threshold"`
		Enabled   bool    `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateRule(id, req.RuleType, req.Threshold, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.db.DeleteRule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetNotifications(c *gin.Context) {
	notifications, err := h.db.GetAllNotifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func (h *Handler) CreateNotification(c *gin.Context) {
	var req struct {
		Name   string                 `json:"name" binding:"required"`
		Type   string                 `json:"type" binding:"required"`
		Config map[string]string      `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configJSON := ""
	if req.Config != nil {
		notif := &models.Notification{}
		notif.SetConfig(req.Config)
		configJSON = notif.Config.String
	}

	notif, err := h.db.CreateNotification(req.Name, req.Type, configJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, notif)
}

func (h *Handler) UpdateNotification(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Name    string            `json:"name"`
		Type    string            `json:"type"`
		Config  map[string]string `json:"config"`
		Enabled bool              `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configJSON := ""
	if req.Config != nil {
		notif := &models.Notification{}
		notif.SetConfig(req.Config)
		configJSON = notif.Config.String
	}

	if err := h.db.UpdateNotification(id, req.Name, req.Type, configJSON, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteNotification(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.db.DeleteNotification(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) TestNotification(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	notif, err := h.db.GetNotificationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	alerts, err := h.db.GetAllAlerts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *Handler) GetAlertStats(c *gin.Context) {
	stats, err := h.db.GetAlertStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":         "ok",
		"monitor_running": true,
	})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.db.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	alertStats, _ := h.db.GetAlertStats()

	c.JSON(http.StatusOK, gin.H{
		"total_addresses":    stats["addresses"],
		"total_rules":        stats["rules"],
		"total_notifications": stats["notifications"],
		"total_alerts":       alertStats["total"],
		"today_alerts":       alertStats["today"],
		"monitor_running":    true,
	})
}