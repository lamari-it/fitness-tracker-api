package utils

import (
	"crypto/tls"
	"fmt"
	"lamari-fit-api/config"
	"log"
	"net/smtp"
	"strings"
)

// EmailService handles sending emails
type EmailService struct {
	host      string
	port      string
	username  string
	password  string
	fromEmail string
	fromName  string
	appURL    string
}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	cfg := config.AppConfig
	return &EmailService{
		host:      cfg.SMTPHost,
		port:      cfg.SMTPPort,
		username:  cfg.SMTPUsername,
		password:  cfg.SMTPPassword,
		fromEmail: cfg.SMTPFromEmail,
		fromName:  cfg.SMTPFromName,
		appURL:    cfg.AppURL,
	}
}

// SendTrainerInvitation sends an invitation email to a potential client
func (e *EmailService) SendTrainerInvitation(toEmail, trainerName, invitationToken string) error {
	subject := fmt.Sprintf("%s has invited you to LamariFit", trainerName)

	invitationLink := fmt.Sprintf("%s/invitations/accept?token=%s", e.appURL, invitationToken)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Trainer Invitation</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2c3e50;">You've Been Invited!</h2>

        <p>Hi there,</p>

        <p><strong>%s</strong> has invited you to become their client on LamariFit.</p>

        <p>LamariFit is a fitness tracking platform that helps trainers and clients work together to achieve fitness goals.</p>

        <div style="margin: 30px 0;">
            <a href="%s" style="background-color: #3498db; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">
                Accept Invitation
            </a>
        </div>

        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #7f8c8d;">%s</p>

        <p>This invitation will expire in 7 days.</p>

        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">

        <p style="color: #7f8c8d; font-size: 12px;">
            If you didn't expect this invitation, you can safely ignore this email.
        </p>
    </div>
</body>
</html>
`, trainerName, invitationLink, invitationLink)

	return e.sendEmail(toEmail, subject, body)
}

// sendEmail sends an email using SMTP
func (e *EmailService) sendEmail(to, subject, body string) error {
	// Check if email is configured
	if e.username == "" || e.password == "" {
		log.Printf("Email not configured. Would send to: %s, Subject: %s", to, subject)
		return nil // Don't fail if email is not configured (for development)
	}

	from := fmt.Sprintf("%s <%s>", e.fromName, e.fromEmail)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%s", e.host, e.port)

	// For TLS connections (port 465)
	if e.port == "465" {
		return e.sendEmailTLS(to, message, auth, addr)
	}

	// For STARTTLS connections (port 587)
	return smtp.SendMail(addr, auth, e.fromEmail, []string{to}, []byte(message))
}

// sendEmailTLS sends email using TLS (for port 465)
func (e *EmailService) sendEmailTLS(to, message string, auth smtp.Auth, addr string) error {
	tlsConfig := &tls.Config{
		ServerName: strings.Split(addr, ":")[0],
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err = client.Mail(e.fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}

// IsConfigured returns true if email service is properly configured
func (e *EmailService) IsConfigured() bool {
	return e.username != "" && e.password != "" && e.host != ""
}
