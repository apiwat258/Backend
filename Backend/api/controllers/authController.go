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

// LoginResponse struct
type LoginResponse struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

// Login handles user authentication
func Login(c *fiber.Ctx) error {
	var req LoginRequest

	// ✅ Bind JSON data
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// ✅ ค้นหาผู้ใช้จาก Database
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// ✅ ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// ✅ สร้าง JWT Token
	token, err := middleware.GenerateToken(user.UserID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// ✅ ตั้งค่าคุกกี้
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,   // ✅ ถ้าใช้ HTTPS ให้เป็น `true`
		SameSite: "None", // ✅ ต้องใช้ `None` ถ้าทำงาน Cross-Site
		Path:     "/",
		Domain:   "10.110.194.195", // ✅ ใส่ให้ตรงกับ Frontend
	})

	// ✅ Debug ตรวจสอบว่าคุกกี้ถูกเซ็ตหรือไม่
	fmt.Println("✅ [Login] Token set in cookie for user:", user.UserID)

	// ✅ ส่ง Role กลับไปให้ Frontend
	response := LoginResponse{
		Message: "Login successful",
		Role:    user.Role,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// ✅ ดึง Role จาก JWT Token ใน Cookie
func GetUserRole(c *fiber.Ctx) error {
	// ✅ ใช้ Fiber API ดึง Token โดยตรง
	tokenString := c.Cookies("auth_token")
	if tokenString == "" {
		fmt.Println("❌ [GetUserRole] No token found in Cookie")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No authentication token found. Please login again.",
		})
	}

	// ✅ Validate Token
	claims, err := middleware.ValidateToken(tokenString)
	if err != nil {
		fmt.Println("❌ [GetUserRole] Invalid token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token. Please login again.",
		})
	}

	fmt.Println("✅ [GetUserRole] Authenticated User Role:", claims.Role)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"role": claims.Role})
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

// Logout API: ลบคุกกี้
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // หมดอายุทันที
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful"})
}

// เพิ่มส่วนนี้ลงใน authController.go

// GetUserInfo handles fetching user information (email and password) from the user table.
func GetUserInfo(c *fiber.Ctx) error {
	// ดึง userID จาก context ที่ถูกตั้งโดย JWT middleware
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// ส่งกลับ email และ password (hashed) ของผู้ใช้
	return c.JSON(fiber.Map{
		"email":    user.Email,
		"password": user.Password,
	})
}

// UpdateUserInfo handles updating user's email and password.
func UpdateUserInfo(c *fiber.Ctx) error {
	// ดึง userID จาก context
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type UpdateRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// hash รหัสผ่านใหม่
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// อัปเดต email และ password ของผู้ใช้ในตาราง user
	if err := database.DB.Model(&models.User{}).Where("userid = ?", userID).Updates(models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user info"})
	}

	return c.JSON(fiber.Map{"message": "User info updated successfully"})
}
