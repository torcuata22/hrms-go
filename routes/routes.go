package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/torcuata22/hrms-mongo/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/employees", controller.GetEmployees)
	app.Get("/employee/:id", controller.GetEmployee)
	app.Post("/employees", controller.CreateEmployee)
	app.Put("/employee/:id", controller.UpdateEmployee)
	app.Delete("/employee/:id", controller.DeleteEmployee)
}
