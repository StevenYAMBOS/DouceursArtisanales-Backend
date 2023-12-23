package controllers

import (
	"fmt"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/stevenyambos/hmmm-backend/models"
	"github.com/stevenyambos/hmmm-backend/pkg/database"
)

func ContactForm(c *fiber.Ctx) error {
	var contact models.Contact
	if err := c.BodyParser(&contact); err != nil {
		return err
	}

	// Insérez la demande de contact dans la base de données
	_, err := database.DB.Exec(`
	INSERT INTO contact_form (id, email, subject, message, submitted_at)
	VALUES ($1, $2, $3, $4, $5)
	`, contact.ID, contact.Email, contact.Subject, contact.Message, time.Now())
	if err != nil {
		fmt.Println("Erreur lors de l'insertion en base de données :", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Erreur lors de l'insertion en base de données"})
	}

	fmt.Printf("Données reçues du formulaire : \nEmail: %s\nSubject: %s\nMessage: %s\n", contact.Email, contact.Subject, contact.Message)

	contact.Submitted_at = time.Now()

	return c.JSON(contact)
}
