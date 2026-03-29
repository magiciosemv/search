package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"net/url"
	"time"

	"solana-monitor/internal/models"
)

type NotificationService struct {
	db            *models.DB
	telegramToken string
	smtpConfig    SMTPConfig
	proxyURL      string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type AlertMessage struct {
	Address   string
	Label     string
	AlertType string
	OldValue  float64
	NewValue  float64
	Time      time.Time
}

func NewNotificationService(db *models.DB, telegramToken string, smtpConfig SMTPConfig, proxyURL string) *NotificationService {
	return &NotificationService{
		db:            db,
		telegramToken: telegramToken,
		smtpConfig:    smtpConfig,
		proxyURL:      proxyURL,
	}
}

func (n *NotificationService) SendTelegram(chatID, message string) error {
	if n.telegramToken == "" || chatID == "" {
		return fmt.Errorf("telegram token or chat ID not configured")
	}

	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.telegramToken)
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var client *http.Client
	if n.proxyURL != "" {
		proxyURLParsed, err := url.Parse(n.proxyURL)
		if err == nil {
			client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURLParsed)}}
		} else {
			client = http.DefaultClient
		}
	} else {
		client = http.DefaultClient
	}

	resp, err := client.Post(telegramURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}

func (n *NotificationService) SendEmail(to, subject, body string) error {
	if n.smtpConfig.Username == "" || n.smtpConfig.Password == "" {
		return fmt.Errorf("SMTP not configured")
	}

	msg := fmt.Sprintf("From: %s\r\n", n.smtpConfig.From) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		body

	auth := smtp.PlainAuth("", n.smtpConfig.Username, n.smtpConfig.Password, n.smtpConfig.Host)
	err := smtp.SendMail(
		n.smtpConfig.Host+":"+n.smtpConfig.Port,
		auth,
		n.smtpConfig.From,
		[]string{to},
		[]byte(msg),
	)
	return err
}

func (n *NotificationService) SendAlert(msg AlertMessage) error {
	if n.db == nil {
		fmt.Printf("Alert: %s balance changed from %.4f to %.4f\n", msg.Address, msg.OldValue, msg.NewValue)
		return nil
	}

	// Get enabled notifications from database
	notifications, err := n.db.GetAllNotifications()
	if err != nil {
		fmt.Printf("Failed to get notifications: %v\n", err)
		return err
	}

	formatted := n.FormatTelegramAlert(msg)
	emailBody := n.FormatEmailAlert(msg)

	for _, notif := range notifications {
		if !notif.Enabled {
			continue
		}

		config := notif.GetConfig()

		if notif.Type == "telegram" {
			if chatID, ok := config["chat_id"]; ok {
				if err := n.SendTelegram(chatID, formatted); err != nil {
					fmt.Printf("Failed to send Telegram notification: %v\n", err)
				} else {
					fmt.Printf("Telegram notification sent to %s\n", chatID)
				}
			}
		} else if notif.Type == "email" {
			if to, ok := config["email"]; ok {
				if err := n.SendEmail(to, "Solana Wallet Alert", emailBody); err != nil {
					fmt.Printf("Failed to send email notification: %v\n", err)
				} else {
					fmt.Printf("Email notification sent to %s\n", to)
				}
			}
		}
	}

	return nil
}

func (n *NotificationService) FormatTelegramAlert(msg AlertMessage) string {
	var emoji string

	change := msg.NewValue - msg.OldValue

	if change > 0 {
		emoji = "📈"
	} else if change < 0 {
		emoji = "📉"
	} else {
		emoji = "➡️"
	}

	address := msg.Address
	if len(address) > 8 {
		address = address[:4] + "..." + address[len(address)-4:]
	}

	return fmt.Sprintf(`%s <b>Solana Wallet Alert</b>

📍 <b>Address:</b> <code>%s</code>
💰 <b>Balance:</b> %.4f → %.4f SOL
📊 <b>Change:</b> %s%.4f SOL

⏰ %s`, emoji, address, msg.OldValue, msg.NewValue, emoji, change, msg.Time.Format("2006-01-02 15:04:05"))
}

func (n *NotificationService) FormatEmailAlert(msg AlertMessage) string {
	change := msg.NewValue - msg.OldValue

	return fmt.Sprintf(`Solana Wallet Alert

Address: %s
Label: %s
Old Balance: %.4f SOL
New Balance: %.4f SOL
Change: %.4f SOL
Time: %s`, msg.Address, msg.Label, msg.OldValue, msg.NewValue, change, msg.Time.Format("2006-01-02 15:04:05"))
}