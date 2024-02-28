package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"net/smtp"
)

var SmtpAuth *smtp.Auth

func AuthSMTP() {
	smtpServer := os.Getenv("GMAIL_SMTP_SERVER")
	smtpUsername := os.Getenv("GMAIL_SMTP_USERNAME")
	smtpPassword := os.Getenv("GMAIL_SMTP_PASSWORD")

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	SmtpAuth = &auth
}

func SendRegistrationEmail(email, name, username string) error {
	smtpServer := os.Getenv("GMAIL_SMTP_SERVER")
	smtpPort := 587
	senderEmail := "sender@example.com"

	// Load the HTML template
	htmlTemplate, err := template.ParseFiles("templates/welcome_email.html")
	if err != nil {
		return fmt.Errorf("failed to parse HTML template : %v", err)
	}

	// Create buffer to store the rendered HTML
	var body bytes.Buffer

	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: Welcome to EcoLens App!\n%s\n\n", headers)))

	if err := htmlTemplate.Execute(&body, map[string]string{"Name": name, "Username": username}); err != nil {
		return fmt.Errorf("failed to execute HTML template : %v", err)
	}

	smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), *SmtpAuth, senderEmail, []string{email}, body.Bytes())
	return nil
}

func SendLoginAlert(email, name, username string) error {
	return nil
}
