package controllers

import (
	"fmt"
	"mime"
	"net/smtp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stevenyambos/hmmm-backend/config"
	"github.com/stevenyambos/hmmm-backend/models"
	"github.com/stevenyambos/hmmm-backend/pkg/database"
)

func ContactForm(c *fiber.Ctx) error {
	var contact models.Contact
	if err := c.BodyParser(&contact); err != nil {
		return err
	}

	// Insérer la demande de contact dans la base de données
	_, err := database.DB.Exec(`
	INSERT INTO contact_form (username, email, subject, message, submitted_at)
	VALUES ($1, $2, $3, $4, $5)
	`, contact.Username, contact.Email, contact.Subject, contact.Message, time.Now())
	if err != nil {
		fmt.Println("Erreur lors de l'insertion en base de données :", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Erreur lors de l'insertion en base de données"})
	}

	fmt.Printf("Données reçues du formulaire : \nEmail: %s\nSubject: %s\nMessage: %s\n", contact.Email, contact.Subject, contact.Message)
	contact.Submitted_at = time.Now()

	// Récupération de l'adresse e-mail de l'expéditeur (elle est dynamique)
	senderEmail := contact.Email

	// Envoi d'e-mail
	err = sendEmail(senderEmail, contact)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de l'Email:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Erreur lors de l'envoi de l'Email"})
	}

	return c.JSON(contact)
}

func sendEmail(senderEmail string, contact models.Contact) error {
	// Initialisation des paramètres SMTP
	emailConfig := config.LoadEmailConfig()

	// Configuration du message e-mail
	message := []byte(
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"Subject: " + mime.QEncoding.Encode("UTF-8", contact.Subject) + "\r\n\r\n" +
		"<p><i>Formulaire de contact DouceursArtisanales.com:</i></p>" +
		"<p><strong>Expéditeur:</strong> " + senderEmail + "</p>" +
		"<p><strong>Nom complet:</strong> " + contact.Username + "</p>" +
		"<p><strong>Objet:<strong> " + contact.Subject + "</i></p>" +
		"<p><strong>Message:<strong> " + contact.Message + "</p>",
	)

	// Configuration de l'adresse e-mail expéditeur
	from := senderEmail

	// Configuration de l'authentification SMTP
	auth := smtp.PlainAuth(
		"",
		emailConfig.SMTP_Username,
		emailConfig.SMTP_Password,
		emailConfig.SMTP_Server,
	)

	// Envoi de l'e-mail
	err := smtp.SendMail(
		emailConfig.SMTP_Server+":"+emailConfig.SMTP_Port,
		auth,
		from,
		[]string{emailConfig.Receiver_Email},
		message,
	)
	fmt.Println("Email envoyé avec succès !")

	return err

}
