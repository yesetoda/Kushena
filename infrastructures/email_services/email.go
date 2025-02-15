package email_services

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SendEmail sends an HTML email with enhanced compatibility and animations.
func SendEmail(toEmail, title, body, link string) error {
	log.Printf("📧 Sending '%s' email to %s", title, toEmail)

	message := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>%s</title>
		<style>
			body {
				margin: 0;
				padding: 0;
				min-width: 100%%;
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen-Sans, Ubuntu, Cantarell, sans-serif;
			}
			.container {
				width: 100%%;
				max-width: 600px;
				margin: 20px auto;
				background: #ffffff;
				border-radius: 12px;
				position: relative;
				overflow: hidden;
				box-shadow: 0 4px 20px rgba(0,0,0,0.1);
			}
			.container::before {
				content: '';
				position: absolute;
				top: -2px;
				left: -2px;
				right: -2px;
				bottom: -2px;
				background: linear-gradient(45deg, #6366f1, #10b981, #3b82f6);
				z-index: -1;
				animation: borderAnimation 3s linear infinite;
				background-size: 400%%;
			}
			@keyframes borderAnimation {
				0%% { background-position: 0%% 50%%; }
				50%% { background-position: 100%% 50%%; }
				100%% { background-position: 0%% 50%%; }
			}
			.content {
				padding: 32px 24px;
				color: #1f2937;
				text-align: center;
			}
			.button {
				display: inline-block;
				padding: 14px 32px;
				background: #3b82f6;
				color: white !important;
				text-decoration: none;
				border-radius: 8px;
				font-weight: 600;
				margin: 20px 0;
				transition: all 0.3s ease;
			}
			.button:hover {
				background: #2563eb;
				transform: translateY(-2px);
				box-shadow: 0 8px 20px rgba(59, 130, 246, 0.3);
			}
			.footer {
				background: #f8fafc;
				padding: 24px;
				text-align: center;
				color: #64748b;
				font-size: 12px;
			}
			@media only screen and (max-width: 600px) {
				.container {
					width: 95%% !important;
				}
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="content">
				<h1 style="color: #1e40af; margin: 0 0 20px 0;">%s</h1>
				<div style="font-size: 16px; line-height: 1.6; color: #374151;">
					%s
				</div>
				<a href="%s" class="button" target="_blank" style="color: white;">Take Action →</a>
				<div style="margin-top: 24px; font-size: 14px; color: #6b7280;">
					If the button doesn't work, copy and paste this link:<br>
					<span style="word-break: break-all;">%s</span>
				</div>
			</div>
			<div class="footer">
				© 2024 Kushena. All rights reserved.
			</div>
		</div>
	</body>
	</html>
	`, title, title, body, link, link)

	// Build email headers.
	headers := make(map[string]string)
	headers["From"] = os.Getenv("EMAIL_SENDER")
	headers["To"] = toEmail
	headers["Subject"] = title
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var emailMessage string
	for k, v := range headers {
		emailMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	emailMessage += "\r\n" + message

	// SMTP configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	addressStr := smtpHost + ":" + smtpPort
	sender := os.Getenv("EMAIL_SENDER")
	emailKey := os.Getenv("EMAIL_KEY")

	if sender == "" || emailKey == "" {
		log.Println("❌ SMTP credentials are missing")
		return fmt.Errorf("SMTP credentials not set")
	}

	auth := smtp.PlainAuth("", sender, emailKey, smtpHost)
	err := smtp.SendMail(addressStr, auth, sender, []string{toEmail}, []byte(emailMessage))
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", toEmail, err)
		return err
	}

	log.Printf("✅ Email successfully sent to %s", toEmail)
	return nil
}

// SendLoginAlertEmail sends a login alert email containing the client's IP, location, and device info.
func SendLoginAlertEmail(c *gin.Context, alertRecipient, link string) error {
	ip := GetClientIP(c)
	device := GetDevice(c)

	// Handle loopback addresses.
	location := "Localhost"
	if ip != "::1" && ip != "127.0.0.1" {
		var err error
		location, err = GetGeoLocation(ip)
		if err != nil {
			log.Printf("❌ Failed to get location for IP %s: %v", ip, err)
			location = "Unknown location"
		}
	}

	title := "Login Attempt Alert"
	body := fmt.Sprintf("A login attempt was made from location: %s, IP address: %s, using device: %s.", location, ip, device)
	// No specific action link is needed for this alert.
	return SendEmail(alertRecipient, title, body, link)
}

// GetClientIP extracts the client's IP address from the gin context.
func GetClientIP(c *gin.Context) string {
	// Check for X-Forwarded-For header in case of a proxy.
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}
	// Fall back to RemoteAddr.
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// GetDevice retrieves the User-Agent string from the gin context.
func GetDevice(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

// GetGeoLocation returns a formatted location string based on the client's IP address.
func GetGeoLocation(ip string) (string, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var info struct {
		City    string `json:"city"`
		Region  string `json:"region"`
		Country string `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s, %s, %s", info.City, info.Region, info.Country), nil
}
