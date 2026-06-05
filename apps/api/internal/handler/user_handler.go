package handler

import (
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler (svc service.UserService) *UserHandler {
	return &UserHandler{service : svc}
}

func (h *UserHandler) Insert(c *fiber.Ctx) error {
	var req *model.User

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "error": "invalid request body",
		})
	}

	err := h.service.Insert(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "error": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true, "data": req,
	})
}