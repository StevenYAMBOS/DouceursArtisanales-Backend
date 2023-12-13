package routes

import (

	"github.com/gofiber/fiber/v2"
	"github.com/stevenyambos/hmmm-backend/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/user-profile", controllers.Protected(), controllers.User)
	app.Get("/", controllers.UserProfile)
	app.Post("/logout", controllers.Logout)
}