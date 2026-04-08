package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"solana-monitor/internal/models"

	"golang.org/x/net/proxy"
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

	fromAddr, err := mail.ParseAddress(n.smtpConfig.From)
	if err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}

	toAddr, err := mail.ParseAddress(to)
	if err != nil {
		return fmt.Errorf("invalid to address: %w", err)
	}

	msg := fmt.Sprintf("From: %s\r\n", fromAddr.String()) +
		fmt.Sprintf("To: %s\r\n", toAddr.String()) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
		"\r\n" +
		body

	addr := n.smtpConfig.Host + ":" + n.smtpConfig.Port

	var conn net.Conn

	// Connect through SOCKS5 proxy if configured
	if n.proxyURL != "" {
		proxyURL, parseErr := url.Parse(n.proxyURL)
		if parseErr != nil {
			return fmt.Errorf("failed to parse proxy URL: %w", parseErr)
		}

		// Extract SOCKS5 proxy address (strip http:// prefix)
		proxyAddr := proxyURL.Host
		if proxyAddr == "" {
			proxyAddr = strings.TrimPrefix(n.proxyURL, "http://")
			proxyAddr = strings.TrimPrefix(proxyAddr, "socks5://")
			proxyAddr = strings.TrimPrefix(proxyAddr, "socks5h://")
		}

		dialer, dialErr := proxy.SOCKS5("tcp", proxyAddr, nil, &net.Dialer{Timeout: 10 * time.Second})
		if dialErr != nil {
			return fmt.Errorf("failed to create SOCKS5 dialer: %w", dialErr)
		}

		conn, err = dialer.Dial("tcp", addr)
	} else {
		conn, err = net.DialTimeout("tcp", addr, 10*time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	// Read SMTP server greeting (with timeout)
	// smtp.NewClient will read the greeting, but we need to set a deadline first
	conn.SetDeadline(time.Now().Add(15 * time.Second))

	tlsConfig := &tls.Config{
		ServerName: n.smtpConfig.Host,
	}

	client, err := smtp.NewClient(conn, n.smtpConfig.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Upgrade to TLS if server supports it
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	auth := smtp.PlainAuth("", n.smtpConfig.Username, n.smtpConfig.Password, n.smtpConfig.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err := client.Mail(fromAddr.Address); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	if err := client.Rcpt(toAddr.Address); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA failed: %w", err)
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close email body: %w", err)
	}

	return client.Quit()
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

	if msg.AlertType == "wallet_added" {
		emoji = "➕"
	} else {
		change := msg.NewValue - msg.OldValue
		if change > 0 {
			emoji = "📈"
		} else if change < 0 {
			emoji = "📉"
		} else {
			emoji = "➡️"
		}
	}

	address := msg.Address
	if len(address) > 8 {
		address = address[:4] + "..." + address[len(address)-4:]
	}

	if msg.AlertType == "wallet_added" {
		return fmt.Sprintf(`%s <b>New Wallet Added</b>

📍 <b>Address:</b> <code>%s</code>
🏷️ <b>Label:</b> %s
💰 <b>Balance:</b> %.4f SOL

⏰ %s`, emoji, address, msg.Label, msg.NewValue, msg.Time.Format("2006-01-02 15:04:05"))
	}

	return fmt.Sprintf(`%s <b>Solana Wallet Alert</b>

📍 <b>Address:</b> <code>%s</code>
💰 <b>Balance:</b> %.4f → %.4f SOL
📊 <b>Change:</b> %s%.4f SOL

⏰ %s`, emoji, address, msg.OldValue, msg.NewValue, emoji, msg.NewValue - msg.OldValue, msg.Time.Format("2006-01-02 15:04:05"))
}

func (n *NotificationService) FormatEmailAlert(msg AlertMessage) string {
	if msg.AlertType == "wallet_added" {
		return fmt.Sprintf(`New Wallet Added

Address: %s
Label: %s
Balance: %.4f SOL
Time: %s`, msg.Address, msg.Label, msg.NewValue, msg.Time.Format("2006-01-02 15:04:05"))
	}

	change := msg.NewValue - msg.OldValue
	return fmt.Sprintf(`Solana Wallet Alert

Address: %s
Label: %s
Old Balance: %.4f SOL
New Balance: %.4f SOL
Change: %.4f SOL
Time: %s`, msg.Address, msg.Label, msg.OldValue, msg.NewValue, change, msg.Time.Format("2006-01-02 15:04:05"))
}