package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"solana-monitor/internal/models"
)

type MonitorService struct {
	db       *models.DB
	solana   *SolanaService
	notifier *NotificationService
	bus      *EventBus
}

func NewMonitorService(db *models.DB, solana *SolanaService, notifier *NotificationService, bus *EventBus) *MonitorService {
	return &MonitorService{
		db:       db,
		solana:   solana,
		notifier: notifier,
		bus:      bus,
	}
}

func (m *MonitorService) Start(ctx context.Context) {
	log.Println("Monitor service started")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Monitor service stopped")
			return
		case <-ticker.C:
			m.checkAllAddresses(ctx)
		}
	}
}

func (m *MonitorService) checkAllAddresses(ctx context.Context) {
	addresses, err := m.db.GetAllAddresses()
	if err != nil {
		log.Printf("Failed to get addresses: %v", err)
		return
	}

	for _, addr := range addresses {
		m.checkAddress(ctx, addr)
	}
}

func (m *MonitorService) checkAddress(ctx context.Context, addr models.Address) {
	balance, err := m.solana.GetBalance(ctx, addr.Address)
	if err != nil {
		log.Printf("Failed to get balance for %s: %v", addr.Address, err)
		return
	}

	oldBalance := addr.Balance
	if balance == oldBalance {
		return
	}

	log.Printf("Balance changed for %s: %.4f -> %.4f SOL", addr.Address, oldBalance, balance)

	// Update stored balance
	if err := m.db.UpdateAddressBalance(addr.ID, balance); err != nil {
		log.Printf("Failed to update balance for %s: %v", addr.Address, err)
	}

	// Publish balance update event
	if m.bus != nil {
		m.bus.Publish("balance_update", map[string]interface{}{
			"id":      addr.ID,
			"address": addr.Address,
			"balance": balance,
		})
	}

	// Get rules for this address
	rules, err := m.db.GetRulesByAddressID(addr.ID)
	if err != nil {
		log.Printf("Failed to get rules for address %d: %v", addr.ID, err)
		return
	}

	// Check each rule
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		shouldAlert, alertType := m.evaluateRule(rule, addr.Balance, balance)
		if shouldAlert {
			m.triggerAlert(ctx, addr, rule, alertType, addr.Balance, balance)
		}
	}
}

func (m *MonitorService) evaluateRule(rule models.Rule, oldBalance, newBalance float64) (bool, string) {
	change := newBalance - oldBalance
	changePercent := 0.0
	if oldBalance > 0 {
		changePercent = (change / oldBalance) * 100
	}

	switch rule.RuleType {
	case "balance_change":
		if change > rule.Threshold {
			return true, "balance_increase"
		}
		if change < -rule.Threshold {
			return true, "balance_decrease"
		}
	case "balance_change_percent":
		if changePercent > rule.Threshold {
			return true, "balance_increase_percent"
		}
		if changePercent < -rule.Threshold {
			return true, "balance_decrease_percent"
		}
	case "threshold_above":
		if newBalance > rule.Threshold {
			return true, "threshold_above"
		}
	case "threshold_below":
		if newBalance < rule.Threshold {
			return true, "threshold_below"
		}
	case "large_incoming":
		if change > rule.Threshold {
			return true, "large_incoming"
		}
	case "large_outgoing":
		if change < -rule.Threshold {
			return true, "large_outgoing"
		}
	}
	return false, ""
}

func (m *MonitorService) triggerAlert(ctx context.Context, addr models.Address, rule models.Rule, alertType string, oldValue, newValue float64) {
	// Save alert to database
	ruleID := rule.ID
	_, err := m.db.CreateAlert(addr.ID, &ruleID, alertType, oldValue, newValue, fmt.Sprintf("Balance changed from %.4f to %.4f", oldValue, newValue))
	if err != nil {
		log.Printf("Failed to save alert: %v", err)
	}

	// Send notifications
	msg := AlertMessage{
		Address:   addr.Address,
		Label:     addr.Label,
		AlertType: alertType,
		OldValue:  oldValue,
		NewValue:  newValue,
		Time:      time.Now(),
	}

	if err := m.notifier.SendAlert(msg); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}

	// Publish alert event
	if m.bus != nil {
		m.bus.Publish("new_alert", map[string]interface{}{
			"address":    msg.Address,
			"label":      msg.Label,
			"alert_type": msg.AlertType,
			"old_value":  msg.OldValue,
			"new_value":  msg.NewValue,
			"time":       msg.Time.Format(time.RFC3339),
		})
	}
}

func (m *MonitorService) ForceCheckAddress(ctx context.Context, addressID int64) error {
	addr, err := m.db.GetAddressByID(addressID)
	if err != nil {
		return err
	}

	m.checkAddress(ctx, *addr)
	return nil
}