package services

import (
	"fmt"
	"log"
	"os"
)

type EmailService struct {
	fromEmail string
	appURL    string
}

func NewEmailService() *EmailService {
	return &EmailService{
		fromEmail: getEnv("EMAIL_FROM", "noreply@autoelys.com"),
		appURL:    getEnv("APP_URL", "http://localhost:8080"),
	}
}

func (s *EmailService) SendPasswordResetEmail(toEmail, resetToken string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.appURL, resetToken)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`
Hello,

You have requested to reset your password. Please click the link below to reset your password:

%s

This link will expire in 1 hour.

If you did not request this, please ignore this email.

Best regards,
AutoElys Team
`, resetURL)

	// In production, use a proper email service (SendGrid, AWS SES, etc.)
	// For now, we'll log the email
	log.Printf("===== PASSWORD RESET EMAIL =====")
	log.Printf("To: %s", toEmail)
	log.Printf("Subject: %s", subject)
	log.Printf("Body:\n%s", body)
	log.Printf("Reset Token: %s", resetToken)
	log.Printf("================================")

	// TODO: Implement actual email sending
	// Example with SMTP:
	// return s.sendSMTP(toEmail, subject, body)

	return nil
}

func (s *EmailService) SendWelcomeEmail(toEmail, firstName string) error {
	subject := "Welcome to AutoElys!"
	body := fmt.Sprintf(`
Hello %s,

Welcome to AutoElys! Your account has been successfully created.

We're excited to have you on board.

Best regards,
AutoElys Team
`, firstName)

	log.Printf("===== WELCOME EMAIL =====")
	log.Printf("To: %s", toEmail)
	log.Printf("Subject: %s", subject)
	log.Printf("Body:\n%s", body)
	log.Printf("=========================")

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
