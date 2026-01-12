package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// EmailConfig holds the email service configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
	FromName     string
	BaseURL      string
}

// EmailService handles email sending operations
type EmailService struct {
	config EmailConfig
}

// NewEmailService creates a new email service with configuration from environment
func NewEmailService() *EmailService {
	return &EmailService{
		config: EmailConfig{
			SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvOrDefault("SMTP_PORT", "587"),
			SMTPUser:     getEnvOrDefault("SMTP_USER", ""),
			SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),
			FromAddress:  getEnvOrDefault("SMTP_FROM_ADDRESS", "noreply@entorno35.com"),
			FromName:     getEnvOrDefault("SMTP_FROM_NAME", "Entorno35 - NOM-035"),
			BaseURL:      getEnvOrDefault("APP_BASE_URL", "http://localhost:3000"),
		},
	}
}

// NewEmailServiceWithConfig creates an email service with explicit configuration
func NewEmailServiceWithConfig(config EmailConfig) *EmailService {
	return &EmailService{config: config}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// AssessmentEmailData contains the data for assessment invitation emails
type AssessmentEmailData struct {
	StaffName     string
	CompanyName   string
	AssessmentURL string
	ExpiresAt     time.Time
	Period        int
}

// SendAssessmentInvitation sends an assessment invitation email
func (s *EmailService) SendAssessmentInvitation(to string, data AssessmentEmailData) error {
	subject := fmt.Sprintf("Evaluacion NOM-035 - %s", data.CompanyName)

	body, err := s.renderAssessmentTemplate(data)
	if err != nil {
		return fmt.Errorf("error al renderizar plantilla: %w", err)
	}

	return s.sendEmail(to, subject, body)
}

// renderAssessmentTemplate renders the assessment invitation email template
func (s *EmailService) renderAssessmentTemplate(data AssessmentEmailData) (string, error) {
	tmpl := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Evaluacion NOM-035</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #1e293b;
            margin: 0;
            padding: 0;
            background-color: #f8fafc;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .email-card {
            background-color: #ffffff;
            border-radius: 12px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #0f172a 0%, #1e3a5f 100%);
            color: #ffffff;
            padding: 32px 24px;
            text-align: center;
        }
        .header h1 {
            margin: 0 0 8px 0;
            font-size: 24px;
            font-weight: 700;
        }
        .header p {
            margin: 0;
            font-size: 14px;
            opacity: 0.9;
        }
        .content {
            padding: 32px 24px;
        }
        .greeting {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 16px;
        }
        .message {
            color: #475569;
            margin-bottom: 24px;
        }
        .cta-button {
            display: inline-block;
            background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
            color: #ffffff !important;
            text-decoration: none;
            padding: 14px 32px;
            border-radius: 8px;
            font-weight: 600;
            font-size: 16px;
            text-align: center;
            margin: 24px 0;
        }
        .info-box {
            background-color: #f0f9ff;
            border: 1px solid #bae6fd;
            border-radius: 8px;
            padding: 16px;
            margin: 24px 0;
        }
        .info-box h3 {
            margin: 0 0 8px 0;
            font-size: 14px;
            color: #0369a1;
        }
        .info-box p {
            margin: 0;
            font-size: 13px;
            color: #0c4a6e;
        }
        .warning-box {
            background-color: #fef3c7;
            border: 1px solid #fbbf24;
            border-radius: 8px;
            padding: 16px;
            margin: 24px 0;
        }
        .warning-box p {
            margin: 0;
            font-size: 13px;
            color: #92400e;
        }
        .footer {
            background-color: #f8fafc;
            padding: 24px;
            text-align: center;
            border-top: 1px solid #e2e8f0;
        }
        .footer p {
            margin: 0;
            font-size: 12px;
            color: #64748b;
        }
        .footer a {
            color: #3b82f6;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="email-card">
            <div class="header">
                <h1>Entorno35</h1>
                <p>Plataforma de Cumplimiento NOM-035-STPS-2018</p>
            </div>
            
            <div class="content">
                <p class="greeting">Hola {{.StaffName}},</p>
                
                <p class="message">
                    {{.CompanyName}} te invita a completar tu evaluacion de riesgo psicosocial 
                    correspondiente al periodo {{.Period}}, conforme a la NOM-035-STPS-2018.
                </p>
                
                <div style="text-align: center;">
                    <a href="{{.AssessmentURL}}" class="cta-button">
                        Comenzar Evaluacion
                    </a>
                </div>
                
                <div class="info-box">
                    <h3>Informacion Importante</h3>
                    <p>
                        • El cuestionario toma aproximadamente 15-20 minutos<br>
                        • Tus respuestas son completamente confidenciales<br>
                        • Solo tu puedes ver tus respuestas individuales
                    </p>
                </div>
                
                <div class="warning-box">
                    <p>
                        <strong>Este enlace expira el {{.ExpiresAt.Format "02/01/2006"}}.</strong><br>
                        Por favor completa la evaluacion antes de esta fecha.
                    </p>
                </div>
                
                <p class="message" style="font-size: 13px;">
                    Si tienes problemas para acceder al enlace, copia y pega la siguiente URL en tu navegador:<br>
                    <span style="color: #3b82f6; word-break: break-all;">{{.AssessmentURL}}</span>
                </p>
            </div>
            
            <div class="footer">
                <p>
                    Este correo fue enviado automaticamente por la plataforma Entorno35.<br>
                    Por favor no responda a este mensaje.
                </p>
                <p style="margin-top: 12px;">
                    <a href="{{.AssessmentURL}}">Acceder a la evaluacion</a>
                </p>
            </div>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("assessment").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// sendEmail sends an email using SMTP
func (s *EmailService) sendEmail(to, subject, body string) error {
	// If SMTP is not configured, log the email instead of sending
	if s.config.SMTPUser == "" || s.config.SMTPPassword == "" {
		log.Printf("[EMAIL] Would send email to: %s\nSubject: %s\nBody length: %d bytes",
			to, subject, len(body))
		return nil // Return success - email is "sent" (logged)
	}

	// Construct email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddress)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build message
	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	// Setup authentication
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(addr, auth, s.config.FromAddress, []string{to}, []byte(message.String()))
	if err != nil {
		return fmt.Errorf("error al enviar correo: %w", err)
	}

	log.Printf("[EMAIL] Email sent successfully to: %s", to)
	return nil
}

// GetAssessmentURL generates the full assessment URL for a token
func (s *EmailService) GetAssessmentURL(token string) string {
	return fmt.Sprintf("%s/assessment/%s", s.config.BaseURL, token)
}

// IsConfigured returns true if email service is properly configured
func (s *EmailService) IsConfigured() bool {
	return s.config.SMTPUser != "" && s.config.SMTPPassword != ""
}
