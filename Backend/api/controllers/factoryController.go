package controllers

import (
	"database/sql"
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ API สำหรับลงทะเบียนโรงงาน (Factory)
func CreateFactory(c *fiber.Ctx) error {
	type FactoryRequest struct {
		UserID       string  `json:"userid"`
		CompanyName  string  `json:"company_name"`
		FirstName    string  `json:"firstname"`
		LastName     string  `json:"lastname"`
		Email        string  `json:"email"`
		Address      string  `json:"address"`
		Address2     *string `json:"address2"`
		AreaCode     *string `json:"areacode"`
		Phone        string  `json:"phone"`
		PostCode     string  `json:"post"`
		City         string  `json:"city"`
		Province     string  `json:"province"`
		Country      string  `json:"country"`
		LineID       *string `json:"lineid"`
		Facebook     *string `json:"facebook"`
		LocationLink *string `json:"location_link"`
	}

	var req FactoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// ✅ ตรวจสอบว่า User ID มีอยู่ในฐานข้อมูล `users` หรือไม่
	var user models.User
	if err := database.DB.Where("userid = ?", req.UserID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User ID not found in users table"})
	}

	// ✅ ตรวจสอบว่าผู้ใช้เคยลงทะเบียนเป็น Factory แล้วหรือไม่
	var existingFactory models.Factory
	err := database.DB.Where("userid = ?", req.UserID).First(&existingFactory).Error

	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already registered as a factory owner"})
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ Log ที่ชัดเจนขึ้น
	fmt.Println("UserID", req.UserID, "is not registered as a factory yet. Proceeding with registration.")

	// ✅ อัปเดต Role ของ User เป็น "factory"
	if err := database.DB.Model(&models.User{}).Where("userid = ?", req.UserID).Update("role", "factory").Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	// ✅ สร้าง FactoryID ใหม่ (DFYYNNNNN)
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('dairyfactory_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate factory ID"})
	}
	yearPrefix := time.Now().Format("06")                      // ใช้เลขสองหลักของปี เช่น 24 สำหรับ 2024
	factoryID := fmt.Sprintf("DF%s%05d", yearPrefix, sequence) // ได้เป็น DFYYNNNNN

	// ✅ รวม `address2` กับ `address`
	fullAddress := strings.TrimSpace(req.Address)
	if req.Address2 != nil && strings.TrimSpace(*req.Address2) != "" {
		fullAddress = fullAddress + ", " + strings.TrimSpace(*req.Address2)
	}

	// ✅ รวม `area code` กับ `phone`
	fullPhone := strings.TrimSpace(req.Phone)
	if req.AreaCode != nil && strings.TrimSpace(*req.AreaCode) != "" {
		areaCode := strings.TrimSpace(*req.AreaCode)
		if !strings.HasPrefix(areaCode, "+") {
			areaCode = "+" + areaCode
		}
		fullPhone = areaCode + " " + fullPhone
	}

	// ✅ ตรวจสอบ `companyname` ถ้าว่างให้ใช้ "N/A"
	companyName := strings.TrimSpace(req.CompanyName)
	if companyName == "" {
		companyName = "N/A"
	}

	// ✅ ถ้า `province` ว่างให้ใช้ `city`
	province := strings.TrimSpace(req.Province)
	if province == "" {
		province = req.City
	}

	// ✅ ตรวจสอบ `email` ถ้าเป็น `""` ให้ใช้ NULL
	//email := sql.NullString{}
	if strings.TrimSpace(req.Email) != "" {
		//email = sql.NullString{String: strings.TrimSpace(req.Email), Valid: true}
	}

	// ✅ แปลง `*string` เป็น `sql.NullString`
	lineID := sql.NullString{}
	if req.LineID != nil && strings.TrimSpace(*req.LineID) != "" {
		lineID = sql.NullString{String: strings.TrimSpace(*req.LineID), Valid: true}
	}

	facebook := sql.NullString{}
	if req.Facebook != nil && strings.TrimSpace(*req.Facebook) != "" {
		facebook = sql.NullString{String: strings.TrimSpace(*req.Facebook), Valid: true}
	}

	locationLink := sql.NullString{}
	if req.LocationLink != nil && strings.TrimSpace(*req.LocationLink) != "" {
		locationLink = sql.NullString{String: strings.TrimSpace(*req.LocationLink), Valid: true}
	}

	if req.UserID == "" {
		fmt.Println("❌ Error: UserID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "UserID is required"})
	}

	// ✅ สร้างข้อมูล Factory
	factory := models.Factory{
		FactoryID: factoryID,
		//UserID:       req.UserID,
		//Username:     strings.TrimSpace(req.FirstName) + " " + strings.TrimSpace(req.LastName),
		CompanyName: companyName,
		Address:     fullAddress,
		//City:         req.City,
		Province:     province,
		Country:      req.Country,
		PostCode:     req.PostCode,
		Telephone:    fullPhone,
		LineID:       lineID,
		Facebook:     facebook,
		LocationLink: locationLink,
		CreatedOn:    time.Now(),
		//Email:        email.String,
	}

	// ✅ บันทึกลง Database
	if err := database.DB.Create(&factory).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save factory data"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Factory registered successfully", "factory_id": factory.FactoryID})
}
