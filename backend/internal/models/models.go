package models

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "modernc.org/sqlite"
)

type Address struct {
	ID        int64     `json:"id"`
	Address   string    `json:"address"`
	Label     string    `json:"label"`
	Balance   float64   `json:"balance,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Rule struct {
	ID        int64     `json:"id"`
	AddressID int64     `json:"address_id"`
	RuleType  string    `json:"rule_type"`
	Threshold float64   `json:"threshold"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Config    sql.NullString `json:"config"`
	ConfigMap map[string]string `json:"config_map,omitempty"`
	Enabled   bool           `json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
}

type Alert struct {
	ID          int64     `json:"id"`
	AddressID   int64     `json:"address_id"`
	RuleID      *int64    `json:"rule_id,omitempty"`
	AlertType   string    `json:"alert_type"`
	OldValue    float64   `json:"old_value"`
	NewValue    float64   `json:"new_value"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"sent_at"`
	Address     *Address  `json:"address,omitempty"`
	AddressStr  string    `json:"address"` // Plain address string for frontend
}

type DB struct {
	*sql.DB
}

func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS addresses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address VARCHAR(44) NOT NULL UNIQUE,
		label VARCHAR(255),
		balance REAL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS rules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address_id INTEGER NOT NULL,
		rule_type VARCHAR(50) NOT NULL,
		threshold REAL NOT NULL,
		enabled BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		config TEXT,
		enabled BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS alerts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address_id INTEGER NOT NULL,
		rule_id INTEGER,
		alert_type VARCHAR(50) NOT NULL,
		old_value REAL NOT NULL,
		new_value REAL NOT NULL,
		message TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,
		FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE SET NULL
	);
	`
	_, err := db.Exec(schema)
	return err
}

func (n *Notification) SetConfig(configMap map[string]string) error {
	data, err := json.Marshal(configMap)
	if err != nil {
		return err
	}
	n.Config = sql.NullString{String: string(data), Valid: true}
	n.ConfigMap = configMap
	return nil
}

func (n *Notification) GetConfig() map[string]string {
	if n.Config.Valid {
		var config map[string]string
		json.Unmarshal([]byte(n.Config.String), &config)
		return config
	}
	return nil
}