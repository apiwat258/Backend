package controllers

import (
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/middleware"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest struct
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles user authentication
func Login(c *fiber.Ctx) error {
	var req LoginRequest

	// Bind JSON data
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// Find user in database
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.UserID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Return token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func UpdateUserRole(c *fiber.Ctx) error {
	type RoleRequest struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	var req RoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	user.Role = req.Role
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role updated successfully"})
}

// Register a new user
func Register(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// ตรวจสอบว่าอีเมลนี้มีอยู่ในระบบแล้วหรือไม่
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already exists"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// ดึงค่าลำดับถัดไปจากลำดับในฐานข้อมูล
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('user_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate user ID"})
	}

	// สร้าง userid ในรูปแบบ YYNNNNN
	yearPrefix := time.Now().Format("06") // ได้เลขปีสองหลัก เช่น "25" สำหรับปี 2025
	userID := fmt.Sprintf("%s%05d", yearPrefix, sequence)

	user := models.User{
		UserID:   userID,
		Email:    req.Email,
		Password: hashedPassword,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully", "user_id": user.UserID})
}
