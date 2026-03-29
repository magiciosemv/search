package models

import (
	"time"
)

func (db *DB) GetAllAddresses() ([]Address, error) {
	rows, err := db.Query("SELECT id, address, label, balance, created_at, updated_at FROM addresses ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var a Address
		if err := rows.Scan(&a.ID, &a.Address, &a.Label, &a.Balance, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func (db *DB) GetAddressByID(id int64) (*Address, error) {
	var a Address
	err := db.QueryRow("SELECT id, address, label, balance, created_at, updated_at FROM addresses WHERE id = ?", id).
		Scan(&a.ID, &a.Address, &a.Label, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (db *DB) GetAddressByAddress(addr string) (*Address, error) {
	var a Address
	err := db.QueryRow("SELECT id, address, label, balance, created_at, updated_at FROM addresses WHERE address = ?", addr).
		Scan(&a.ID, &a.Address, &a.Label, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (db *DB) CreateAddress(address, label string) (*Address, error) {
	result, err := db.Exec("INSERT INTO addresses (address, label) VALUES (?, ?)", address, label)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return db.GetAddressByID(id)
}

func (db *DB) UpdateAddress(id int64, label string) error {
	_, err := db.Exec("UPDATE addresses SET label = ?, updated_at = ? WHERE id = ?", label, time.Now(), id)
	return err
}

func (db *DB) DeleteAddress(id int64) error {
	_, err := db.Exec("DELETE FROM addresses WHERE id = ?", id)
	return err
}

func (db *DB) UpdateAddressBalance(id int64, balance float64) error {
	_, err := db.Exec("UPDATE addresses SET balance = ?, updated_at = ? WHERE id = ?", balance, time.Now(), id)
	return err
}

func (db *DB) GetAllRules() ([]Rule, error) {
	rows, err := db.Query("SELECT id, address_id, rule_type, threshold, enabled, created_at FROM rules ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var r Rule
		if err := rows.Scan(&r.ID, &r.AddressID, &r.RuleType, &r.Threshold, &r.Enabled, &r.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func (db *DB) GetRulesByAddressID(addressID int64) ([]Rule, error) {
	rows, err := db.Query("SELECT id, address_id, rule_type, threshold, enabled, created_at FROM rules WHERE address_id = ?", addressID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var r Rule
		if err := rows.Scan(&r.ID, &r.AddressID, &r.RuleType, &r.Threshold, &r.Enabled, &r.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func (db *DB) CreateRule(addressID int64, ruleType string, threshold float64) (*Rule, error) {
	result, err := db.Exec("INSERT INTO rules (address_id, rule_type, threshold) VALUES (?, ?, ?)", addressID, ruleType, threshold)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	err = db.QueryRow("SELECT id, address_id, rule_type, threshold, enabled, created_at FROM rules WHERE id = ?", id).
		Scan(&id, &addressID, &ruleType, &threshold, new(bool), new(time.Time))
	if err != nil {
		return nil, err
	}
	return &Rule{ID: id, AddressID: addressID, RuleType: ruleType, Threshold: threshold, Enabled: true}, nil
}

func (db *DB) UpdateRule(id int64, ruleType string, threshold float64, enabled bool) error {
	_, err := db.Exec("UPDATE rules SET rule_type = ?, threshold = ?, enabled = ? WHERE id = ?", ruleType, threshold, enabled, id)
	return err
}

func (db *DB) DeleteRule(id int64) error {
	_, err := db.Exec("DELETE FROM rules WHERE id = ?", id)
	return err
}

func (db *DB) GetAllNotifications() ([]Notification, error) {
	rows, err := db.Query("SELECT id, name, type, config, enabled, created_at FROM notifications ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.Name, &n.Type, &n.Config, &n.Enabled, &n.CreatedAt); err != nil {
			return nil, err
		}
		n.ConfigMap = n.GetConfig()
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (db *DB) CreateNotification(name, ntype string, config string) (*Notification, error) {
	result, err := db.Exec("INSERT INTO notifications (name, type, config) VALUES (?, ?, ?)", name, ntype, config)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return db.GetNotificationByID(id)
}

func (db *DB) GetNotificationByID(id int64) (*Notification, error) {
	var n Notification
	err := db.QueryRow("SELECT id, name, type, config, enabled, created_at FROM notifications WHERE id = ?", id).
		Scan(&n.ID, &n.Name, &n.Type, &n.Config, &n.Enabled, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	n.ConfigMap = n.GetConfig()
	return &n, nil
}

func (db *DB) UpdateNotification(id int64, name, ntype, config string, enabled bool) error {
	_, err := db.Exec("UPDATE notifications SET name = ?, type = ?, config = ?, enabled = ? WHERE id = ?", name, ntype, config, enabled, id)
	return err
}

func (db *DB) DeleteNotification(id int64) error {
	_, err := db.Exec("DELETE FROM notifications WHERE id = ?", id)
	return err
}

func (db *DB) GetAllAlerts(limit, offset int) ([]Alert, error) {
	rows, err := db.Query(`
		SELECT a.id, a.address_id, a.rule_id, a.alert_type, a.old_value, a.new_value, a.message, a.created_at,
		       ad.address, ad.label
		FROM alerts a
		LEFT JOIN addresses ad ON a.address_id = ad.id
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var al Alert
		var addr Address
		if err := rows.Scan(&al.ID, &al.AddressID, &al.RuleID, &al.AlertType, &al.OldValue, &al.NewValue, &al.Message, &al.CreatedAt, &addr.Address, &addr.Label); err != nil {
			return nil, err
		}
		if addr.Address != "" {
			al.Address = &addr
			al.AddressStr = addr.Address
		}
		alerts = append(alerts, al)
	}
	return alerts, nil
}

func (db *DB) CreateAlert(addressID int64, ruleID *int64, alertType string, oldValue, newValue float64, message string) (*Alert, error) {
	result, err := db.Exec(
		"INSERT INTO alerts (address_id, rule_id, alert_type, old_value, new_value, message) VALUES (?, ?, ?, ?, ?, ?)",
		addressID, ruleID, alertType, oldValue, newValue, message)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &Alert{
		ID:        id,
		AddressID: addressID,
		RuleID:    ruleID,
		AlertType: alertType,
		OldValue:  oldValue,
		NewValue:  newValue,
		Message:   message,
	}, nil
}

func (db *DB) GetAlertStats() (map[string]int, error) {
	var total int
	db.QueryRow("SELECT COUNT(*) FROM alerts").Scan(&total)

	var today int
	db.QueryRow("SELECT COUNT(*) FROM alerts WHERE date(created_at) = date('now')").Scan(&today)

	return map[string]int{
		"total": total,
		"today": today,
	}, nil
}

func (db *DB) GetStats() (map[string]int, error) {
	var addresses, rules, notifications, alerts int
	db.QueryRow("SELECT COUNT(*) FROM addresses").Scan(&addresses)
	db.QueryRow("SELECT COUNT(*) FROM rules WHERE enabled = 1").Scan(&rules)
	db.QueryRow("SELECT COUNT(*) FROM notifications WHERE enabled = 1").Scan(&notifications)
	db.QueryRow("SELECT COUNT(*) FROM alerts").Scan(&alerts)

	return map[string]int{
		"addresses":    addresses,
		"rules":        rules,
		"notifications": notifications,
		"alerts":       alerts,
	}, nil
}