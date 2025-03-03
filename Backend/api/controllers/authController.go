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

// RefreshTokenHandler - API สำหรับอัปเดต Token ใหม่
func RefreshTokenHandler(c *fiber.Ctx) error {
	// ✅ ดึง Token จาก Cookie หรือ Header
	tokenString := c.Cookies("auth_token")
	if tokenString == "" {
		tokenString = c.Get("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
	}

	if tokenString == "" {
		fmt.Println("❌ [RefreshToken] No token found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
	}

	// ✅ ตรวจสอบความถูกต้องของ Token
	claims, err := middleware.ValidateToken(tokenString)
	if err != nil {
		fmt.Println("❌ [RefreshToken] Invalid token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// ✅ ค้นหาผู้ใช้จาก Database
	var user models.User
	result := database.DB.Where("user_id = ?", claims.UserID).First(&user)
	if result.Error != nil {
		fmt.Println("❌ [RefreshToken] User not found:", claims.UserID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// ✅ สร้าง Token ใหม่ที่มี Entity ID ล่าสุด
	//newToken, err := middleware.GenerateToken(user.UserID, user.Email, user.Role, user.EntityID)
	if err != nil {
		fmt.Println("❌ [RefreshToken] Failed to generate token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// ✅ อัปเดต Cookie
	c.Cookie(&fiber.Cookie{
		Name: "auth_token",
		//Value:    newToken,
		Path:     "/",
		HTTPOnly: true,
	})

	// ✅ ส่ง Token ใหม่กลับไปให้ User
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token refreshed successfully",
		//"token":   newToken,
	})
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

	// ✅ ดึง `walletAddress` ตาม Role ของ User
	var walletAddress string
	switch user.Role {
	case "farmer":
		var farmer models.Farmer
		if err := database.DB.Where("farmerid = ?", user.EntityID).First(&farmer).Error; err == nil {
			walletAddress = farmer.WalletAddress
		}
	case "factory":
		var factory models.Factory
		if err := database.DB.Where("factoryid = ?", user.EntityID).First(&factory).Error; err == nil {
			walletAddress = factory.WalletAddress
		}
	case "logistics":
		var logistics models.Logistics
		if err := database.DB.Where("logisticsid = ?", user.EntityID).First(&logistics).Error; err == nil {
			walletAddress = logistics.WalletAddress
		}
	case "retailer":
		var retailer models.Retailer
		if err := database.DB.Where("retailerid = ?", user.EntityID).First(&retailer).Error; err == nil {
			walletAddress = retailer.WalletAddress
		}
	}

	// ✅ Debug: ตรวจสอบค่า `walletAddress`
	fmt.Println("🔍 [Login] Extracted WalletAddress:", walletAddress, "for Role:", user.Role, "EntityID:", user.EntityID)

	// ✅ Debug: ตรวจสอบค่า `walletAddress`
	fmt.Println("🔍 [Login] Extracted WalletAddress:", walletAddress, "for Role:", user.Role, "EntityID:", user.EntityID)

	// ✅ สร้าง JWT Token โดยเพิ่ม `walletAddress`
	token, err := middleware.GenerateToken(user.UserID, user.Email, user.Role, user.EntityID, walletAddress)
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
		Domain:   "", // ✅ ใส่ให้ตรงกับ Frontend
	})

	// ✅ Debug ตรวจสอบว่าคุกกี้ถูกเซ็ตหรือไม่
	fmt.Println("✅ [Login] Token set in cookie for user:", user.UserID, "EntityID:", user.EntityID, "Wallet:", walletAddress)

	// ✅ ส่ง Role, EntityID และ WalletAddress กลับไปให้ Frontend
	response := fiber.Map{
		"message":       "Login successful",
		"role":          user.Role,
		"entityID":      user.EntityID,
		"walletAddress": walletAddress,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// ✅ ดึง Role และ EntityID จาก JWT Token ใน Cookie
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

	fmt.Println("✅ [GetUserRole] Authenticated User Role:", claims.Role, "EntityID:", claims.EntityID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"role":     claims.Role,
		"entityID": claims.EntityID, // ✅ เพิ่ม EntityID กลับไปด้วย
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

	fmt.Println("📌 [UpdateUserRole] Updating role for:", req.Email, "New Role:", req.Role)

	// ✅ ตรวจสอบว่า Role ที่ส่งมาเป็นค่าที่ถูกต้องหรือไม่
	validRoles := map[string]bool{
		"farmer":    true,
		"factory":   true,
		"logistics": true,
		"retailer":  true,
		"user":      true,
	}

	if !validRoles[req.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role"})
	}

	// ✅ ค้นหาผู้ใช้จากฐานข้อมูล
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// ✅ อัปเดต Role ของผู้ใช้ (ไม่เช็ค entityID)
	user.Role = req.Role
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	fmt.Println("✅ [UpdateUserRole] User role updated successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role updated successfully"})
}

// Register a new user
func Register(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
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

	// ตรวจสอบว่า Username ซ้ำหรือไม่
	var existingUsername models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUsername).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// ✅ ดึงเลขลำดับถัดไปจาก Database
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('user_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate user ID"})
	}

	// ✅ สร้าง UserID ในรูปแบบ `YYNNNNN` (ปี + เลขลำดับ)
	yearPrefix := time.Now().Format("06")                 // ได้เลขปีสองหลัก เช่น "25" สำหรับปี 2025
	userID := fmt.Sprintf("%s%05d", yearPrefix, sequence) // ตัวเลข 5 หลัก, 00001-99999

	// กำหนดค่า EntityID เป็นค่า Default (ยังไม่ได้เลือก Role)
	entityID := "PENDING_ROLE" // ค่าเริ่มต้นก่อนเลือก Role

	// ✅ สร้าง User ในระบบ
	user := models.User{
		UserID:    userID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "pending", // ยังไม่มี Role ที่เลือก
		EntityID:  entityID,  // ยังไม่ได้เชื่อมกับ Entity
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// ✅ บันทึก User ลง Database
	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "User registered successfully",
		"user_id":  user.UserID,
		"username": user.Username,
		"email":    user.Email,
	})
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

// CheckEmailAvailability ตรวจสอบว่าอีเมลซ้ำหรือไม่
func CheckEmailAvailability(c *fiber.Ctx) error {
	// ✅ ดึงค่าพารามิเตอร์ `email` จาก Query String
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	// ✅ ตรวจสอบอีเมลในฐานข้อมูล
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		// ❌ อีเมลนี้ถูกใช้แล้ว
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"email":     email,
			"available": false,
		})
	}

	// ✅ อีเมลนี้สามารถใช้ได้
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"email":     email,
		"available": true,
	})
}
