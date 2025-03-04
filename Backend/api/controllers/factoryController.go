package controllers

import (
	"database/sql"
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateFactory(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [CreateFactory] Creating factory for userID:", userID)

	// ✅ รับค่าจาก `FormData`
	companyName := strings.TrimSpace(c.FormValue("company_name"))
	email := strings.TrimSpace(c.FormValue("email"))
	address := strings.TrimSpace(c.FormValue("address"))
	district := strings.TrimSpace(c.FormValue("district"))
	subdistrict := strings.TrimSpace(c.FormValue("subdistrict"))
	province := strings.TrimSpace(c.FormValue("province"))
	country := strings.TrimSpace(c.FormValue("country"))
	postCode := strings.TrimSpace(c.FormValue("postcode"))
	phone := strings.TrimSpace(c.FormValue("phone"))
	areaCode := strings.TrimSpace(c.FormValue("areaCode"))
	location := strings.TrimSpace(c.FormValue("location_link"))
	certCID := strings.TrimSpace(c.FormValue("cert_cid"))

	// ✅ รวมเบอร์โทร
	fullPhone := fmt.Sprintf("%s %s", areaCode, phone)

	// ✅ ตรวจสอบ Role ของ User
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.Role == "factory" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already has a factory role"})
	}

	// ✅ ตรวจสอบ Email ซ้ำ
	var existingFactory models.Factory
	if err := database.DB.Where("email = ?", email).First(&existingFactory).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Factory email is already in use"})
	}

	// ✅ ถ้ามีใบเซอร์ → ตรวจสอบว่าซ้ำใน Blockchain หรือไม่
	if certCID != "" {
		cidUnique, err := services.BlockchainServiceInstance.CheckUserCertification(certCID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check certification CID"})
		}
		if !cidUnique {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certification CID already exists in Blockchain"})
		}
	}

	// ✅ สร้าง `factoryID`
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('factory_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate factory ID"})
	}
	yearPrefix := time.Now().Format("06")
	factoryID := fmt.Sprintf("FA%s%05d", yearPrefix, sequence)

	// ✅ ดึง Wallet จาก Ganache
	walletAddress := getGanacheAccount()

	// ✅ ลงทะเบียน User บน Blockchain (ถ้ายังไม่ได้ลงทะเบียน)
	txHash, err := services.BlockchainServiceInstance.RegisterUserOnBlockchain(walletAddress, 2) // 2 = Factory Role
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user on blockchain"})
	}
	fmt.Println("✅ User registered on Blockchain. Transaction Hash:", txHash)

	// ✅ ถ้ามีใบเซอร์ → บันทึกลง Blockchain
	if certCID != "" {
		// ✅ สร้าง `eventID`
		eventID := fmt.Sprintf("EVENT-%s-%s", factoryID, uuid.New().String())

		// ✅ วันที่ออกใบรับรอง และวันหมดอายุ (1 ปี)
		issuedDate := time.Now()
		expiryDate := issuedDate.AddDate(1, 0, 0)

		// ✅ แปลงวันที่เป็น *big.Int
		issuedDateBigInt := big.NewInt(issuedDate.Unix())
		expiryDateBigInt := big.NewInt(expiryDate.Unix())

		// ✅ บันทึกใบเซอร์ลง Blockchain
		certTxHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(
			walletAddress, // ✅ เพิ่ม Wallet Address
			eventID,
			"factory",
			factoryID,
			certCID,
			issuedDateBigInt,
			expiryDateBigInt,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store certification on blockchain"})
		}
		fmt.Println("✅ Certification stored on Blockchain. Transaction Hash:", certTxHash)
	}

	// ✅ บันทึกลง Database
	factory := models.Factory{
		FactoryID:     factoryID,
		CompanyName:   companyName,
		Address:       address,
		District:      district,
		SubDistrict:   subdistrict,
		Province:      province,
		Country:       country,
		PostCode:      postCode,
		Telephone:     fullPhone,
		WalletAddress: walletAddress,
		LocationLink:  sql.NullString{String: location, Valid: location != ""},
		CreatedOn:     time.Now(),
	}

	if err := database.DB.Create(&factory).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save factory data"})
	}

	// ✅ อัปเดต `entityID` และ Role ใน `users`
	updateData := map[string]interface{}{
		"entityid": factoryID,
		"role":     "factory",
	}
	if err := database.DB.Model(&models.User{}).Where("userid = ?", userID).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Factory registered successfully",
		"factory_id":    factoryID,
		"factory_email": email,
		"walletAddress": walletAddress,
		"location_link": location,
		"cert_cid":      certCID, // ✅ ส่ง CID ของใบเซอร์กลับไป (ถ้ามี)
	})
}
