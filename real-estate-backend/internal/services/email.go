package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/smtp"

	"real-estate-backend/internal/config"
)

// EmailService handles email sending functionality
type EmailService struct {
	config *config.EmailConfig
}

// NewEmailService creates a new email service
func NewEmailService(config *config.EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// VerificationEmailData holds data for verification email template
type VerificationEmailData struct {
	UserName         string
	VerificationURL  string
	CompanyName      string
	SupportEmail     string
	ExpirationHours  int
}

// SendVerificationEmail sends an email verification email
func (s *EmailService) SendVerificationEmail(to, userName, verificationToken string) error {
	// Create verification URL
	verificationURL := fmt.Sprintf("%s/api/v1/verify-email?token=%s", s.config.BaseURL, verificationToken)
	
	// Prepare email data
	data := VerificationEmailData{
		UserName:         userName,
		VerificationURL:  verificationURL,
		CompanyName:      "Real Estate Platform",
		SupportEmail:     s.config.SupportEmail,
		ExpirationHours:  24,
	}

	// Generate email content
	subject := "Verify Your Email Address"
	htmlBody, err := s.generateVerificationEmailHTML(data)
	if err != nil {
		return fmt.Errorf("failed to generate email content: %w", err)
	}

	textBody := s.generateVerificationEmailText(data)

	return s.sendEmail(to, subject, textBody, htmlBody)
}

// generateVerificationEmailHTML generates HTML email content
func (s *EmailService) generateVerificationEmailHTML(data VerificationEmailData) (string, error) {
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2c5aa0; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9f9f9; padding: 30px; border-radius: 0 0 5px 5px; }
        .button { display: inline-block; background-color: #28a745; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .button:hover { background-color: #218838; }
        .footer { margin-top: 30px; font-size: 12px; color: #666; text-align: center; }
        .warning { background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; border-radius: 5px; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.CompanyName}}</h1>
        <h2>Email Verification Required</h2>
    </div>
    <div class="content">
        <p>Hello {{.UserName}},</p>
        
        <p>Thank you for registering with {{.CompanyName}}! To complete your registration and start using our platform, please verify your email address.</p>
        
        <p style="text-align: center;">
            <a href="{{.VerificationURL}}" class="button">Verify Email Address</a>
        </p>
        
        <div class="warning">
            <strong>⚠️ Important:</strong> This verification link will expire in {{.ExpirationHours}} hours. If you didn't create an account with us, please ignore this email.
        </div>
        
        <p>If the button above doesn't work, you can copy and paste this link into your browser:</p>
        <p style="word-break: break-all; background-color: #e9ecef; padding: 10px; border-radius: 3px; font-family: monospace;">{{.VerificationURL}}</p>
        
        <p>If you have any questions or need assistance, please contact our support team at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>.</p>
        
        <p>Best regards,<br>The {{.CompanyName}} Team</p>
    </div>
    <div class="footer">
        <p>This is an automated email. Please do not reply to this message.</p>
        <p>&copy; 2025 {{.CompanyName}}. All rights reserved.</p>
    </div>
</body>
</html>`

	tmpl, err := template.New("verification").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateVerificationEmailText generates plain text email content
func (s *EmailService) generateVerificationEmailText(data VerificationEmailData) string {
	return fmt.Sprintf(`
Hello %s,

Thank you for registering with %s! To complete your registration and start using our platform, please verify your email address.

Click the link below to verify your email:
%s

IMPORTANT: This verification link will expire in %d hours. If you didn't create an account with us, please ignore this email.

If you have any questions or need assistance, please contact our support team at %s.

Best regards,
The %s Team

---
This is an automated email. Please do not reply to this message.
© 2025 %s. All rights reserved.
`, data.UserName, data.CompanyName, data.VerificationURL, data.ExpirationHours, data.SupportEmail, data.CompanyName, data.CompanyName)
}

// sendEmail sends an email using SMTP
func (s *EmailService) sendEmail(to, subject, textBody, htmlBody string) error {
	// Create message
	message := fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("From: %s\r\n", s.config.FromEmail)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: multipart/alternative; boundary=\"boundary123\"\r\n"
	message += "\r\n"
	message += "--boundary123\r\n"
	message += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	message += "\r\n"
	message += textBody + "\r\n"
	message += "--boundary123\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	message += "\r\n"
	message += htmlBody + "\r\n"
	message += "--boundary123--\r\n"

	// Setup authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Send email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		auth,
		s.config.FromEmail,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken() string {
	bytes := make([]byte, 32) // 256 bits
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// PasswordResetEmailData holds data for password reset email template
type PasswordResetEmailData struct {
	UserName       string
	ResetURL       string
	CompanyName    string
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, userName, resetToken string) error {
	// Create password reset URL
	resetURL := fmt.Sprintf("%s/web/reset-password?token=%s", s.config.BaseURL, resetToken)

	// Prepare email data
	data := PasswordResetEmailData{
		UserName:    userName,
		ResetURL:    resetURL,
		CompanyName: "Real Estate Platform",
	}

	// Generate email content
	subject := "Reset Your Password"
	htmlBody, err := s.generatePasswordResetEmailHTML(data)
	if err != nil {
		return fmt.Errorf("failed to generate email content: %w", err)
	}

	textBody := s.generatePasswordResetEmailText(data)

	return s.sendEmail(to, subject, textBody, htmlBody)
}

// generatePasswordResetEmailHTML generates HTML email content for password reset
func (s *EmailService) generatePasswordResetEmailHTML(data PasswordResetEmailData) (string, error) {
	templateString := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #c0392b; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9f9f9; padding: 30px; border-radius: 0 0 5px 5px; }
        .button { display: inline-block; background-color: #e74c3c; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .button:hover { background-color: #c0392b; }
        .footer { margin-top: 30px; font-size: 12px; color: #666; text-align: center; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.CompanyName}}</h1>
        <h2>Password Reset Request</h2>
    </div>
    <div class="content">
        <p>Hello {{.UserName}},</p>
        <p>We received a request to reset your password for your account at {{.CompanyName}}. You can reset your password by clicking the link below:</p>
        <p style="text-align: center;">
            <a href="{{.ResetURL}}" class="button">Reset Password</a>
        </p>
        <p>If you did not request a password reset, please ignore this email or contact support if you have questions.</p>
        <p>Thank you,<br>The {{.CompanyName}} Team</p>
    </div>
    <div class="footer">
        <p>This is an automated email. Please do not reply to this message.</p>
    </div>
</body>
</html>`
	tmpl, err := template.New("password_reset").Parse(templateString)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// generatePasswordResetEmailText generates plain text email content for password reset
func (s *EmailService) generatePasswordResetEmailText(data PasswordResetEmailData) string {
	return fmt.Sprintf(`
Hello %s,

We received a request to reset your password for your account at %s. You can reset your password by clicking the link below:
%s

If you did not request a password reset, please ignore this email or contact support if you have questions.

Thank you,
The %s Team
`, data.UserName, data.CompanyName, data.ResetURL, data.CompanyName)
}

// SendWelcomeEmail sends a welcome email after successful verification
func (s *EmailService) SendWelcomeEmail(to, userName string) error {
	subject := "Welcome to Real Estate Platform!"
	
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #28a745; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9f9f9; padding: 30px; border-radius: 0 0 5px 5px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎉 Welcome to Real Estate Platform!</h1>
    </div>
    <div class="content">
        <p>Hello %s,</p>
        <p>Congratulations! Your email has been successfully verified and your account is now active.</p>
        <p>You can now:</p>
        <ul>
            <li>Browse property listings</li>
            <li>Save your favorite properties</li>
            <li>Contact landlords directly</li>
            <li>Create and manage your profile</li>
        </ul>
        <p>If you're a landlord, you can also start listing your properties!</p>
        <p>Thank you for joining our community!</p>
        <p>Best regards,<br>The Real Estate Platform Team</p>
    </div>
</body>
</html>`, userName)

	textBody := fmt.Sprintf(`
Hello %s,

Congratulations! Your email has been successfully verified and your account is now active.

You can now:
- Browse property listings
- Save your favorite properties  
- Contact landlords directly
- Create and manage your profile

If you're a landlord, you can also start listing your properties!

Thank you for joining our community!

Best regards,
The Real Estate Platform Team
`, userName)

	return s.sendEmail(to, subject, textBody, htmlBody)
}
