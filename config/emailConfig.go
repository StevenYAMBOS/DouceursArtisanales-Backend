package config

import (
	"os"
)

type EmailConfig struct {
	SMTP_Server    string
	SMTP_Port      string
	SMTP_Username  string
	SMTP_Password  string
	Sender_Email   string
	Receiver_Email string
}

// InitializeEmailConfig initialise la configuration pour l'envoi d'e-mails
func LoadEmailConfig() *EmailConfig{
	return &EmailConfig{
	// Configuration des paramètres SMTP et adresses Email
	SMTP_Server : os.Getenv("SMTP_SERVER"),
	SMTP_Port : os.Getenv("SMTP_PORT"),
	SMTP_Username : os.Getenv("SMTP_USERNAME"),
	SMTP_Password : os.Getenv("SMTP_PASSWORD"),
	Sender_Email : os.Getenv("SENDER_EMAIL"),
	Receiver_Email : os.Getenv("RECEIVER_EMAIL"),
	}
}

