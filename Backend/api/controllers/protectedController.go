package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ProtectedRoute - API ที่ต้องใช้ JWT Token
func ProtectedRoute(c *fiber.Ctx) error {
	// ดึงค่าจาก JWT Middleware
	userID := c.Locals("userID")
	email := c.Locals("email")

	if userID == nil || email == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Welcome to the protected route!",
		"user_id": userID,
		"email":   email,
	})
}
