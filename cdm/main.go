package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stevenyambos/hmmm-backend/pkg/database"
	"github.com/stevenyambos/hmmm-backend/routes"
)

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

func main() {

	// Lancement de la BDD
	database.InitDB()

	// Instance de Fiber
	var app = fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	if err := app.Listen(":3000"); err != nil {
		fmt.Print("Connexion à la Base de donnée établit")
		panic(err)
	}

}
