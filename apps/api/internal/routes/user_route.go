package routes

import (
	"corporate-translator-api/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, getHandler *handler.UserHandler){
	api := app.Group("/api/v1")

	user := api.Group("/task")
	user.Post("/", getHandler.Insert)
}