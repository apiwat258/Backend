package controllers

import (
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"time"

	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"
	"finalyearproject/Backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateLogistics(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [CreateLogistics] Creating logistics provider for userID:", userID)

	// ✅ รับค่าจาก `FormData`
	companyName := strings.TrimSpace(c.FormValue("logisticsName"))
	email := strings.TrimSpace(c.FormValue("email"))
	address := strings.TrimSpace(c.FormValue("address"))
	district := strings.TrimSpace(c.FormValue("district"))
	subdistrict := strings.TrimSpace(c.FormValue("subdistrict"))
	province := strings.TrimSpace(c.FormValue("province"))
	postCode := strings.TrimSpace(c.FormValue("postCode"))
	phone := strings.TrimSpace(c.FormValue("phone"))
	areaCode := strings.TrimSpace(c.FormValue("areaCode"))
	location := strings.TrimSpace(c.FormValue("location_link"))
	certCID := strings.TrimSpace(c.FormValue("cert_cid"))
	lineID := strings.TrimSpace(c.FormValue("lineID"))
	facebook := strings.TrimSpace(c.FormValue("facebook"))

	// ✅ รวมเบอร์โทร
	fullPhone := fmt.Sprintf("%s %s", areaCode, phone)

	// ✅ ตรวจสอบ Role ของ User
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.Role == "logistics" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already has a logistics role"})
	}

	// ✅ ตรวจสอบ Email ซ้ำ
	var existingLogistics models.Logistics
	if err := database.DB.Where("email = ?", email).First(&existingLogistics).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Logistics email is already in use"})
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

	// ✅ สร้าง `logisticsID`
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('logistics_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate logistics ID"})
	}
	yearPrefix := time.Now().Format("06")
	logisticsID := fmt.Sprintf("LOG%s%05d", yearPrefix, sequence)

	// ✅ ดึง Wallet จาก Ganache
	walletAddress := getGanacheAccount()

	// ✅ ลงทะเบียน User บน Blockchain (ถ้ายังไม่ได้ลงทะเบียน)
	txHash, err := services.BlockchainServiceInstance.RegisterUserOnBlockchain(walletAddress, 3) // 3 = Logistics Role
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user on blockchain"})
	}
	fmt.Println("✅ User registered on Blockchain. Transaction Hash:", txHash)

	// ✅ ถ้ามีใบเซอร์ → บันทึกลง Blockchain
	if certCID != "" {
		// ✅ สร้าง `eventID`
		eventID := fmt.Sprintf("EVENT-%s-%s", logisticsID, uuid.New().String())

		// ✅ วันที่ออกใบรับรอง และวันหมดอายุ (1 ปี)
		issuedDate := time.Now()
		expiryDate := issuedDate.AddDate(1, 0, 0)

		// ✅ แปลงวันที่เป็น *big.Int
		issuedDateBigInt := big.NewInt(issuedDate.Unix())
		expiryDateBigInt := big.NewInt(expiryDate.Unix())

		// ✅ บันทึกใบเซอร์ลง Blockchain
		certTxHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(
			walletAddress,
			eventID,
			"logistics",
			logisticsID,
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
	logistics := models.Logistics{
		LogisticsID:   logisticsID,
		CompanyName:   companyName,
		Address:       address,
		District:      district,
		SubDistrict:   subdistrict,
		Province:      province,
		Country:       "Thailand",
		PostCode:      postCode,
		Telephone:     fullPhone,
		Email:         email,
		WalletAddress: walletAddress,
		LocationLink:  sql.NullString{String: location, Valid: location != ""},
		LineID:        sql.NullString{String: lineID, Valid: lineID != ""},
		Facebook:      sql.NullString{String: facebook, Valid: facebook != ""},
		CreatedOn:     time.Now(),
	}

	if err := database.DB.Create(&logistics).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save logistics data"})
	}

	// ✅ อัปเดต `entityID` และ Role ใน `users`
	updateData := map[string]interface{}{
		"entityid": logisticsID,
		"role":     "logistics",
	}
	if err := database.DB.Model(&models.User{}).Where("userid = ?", userID).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":         "Logistics provider registered successfully",
		"logistics_id":    logisticsID,
		"logistics_email": email,
		"walletAddress":   walletAddress,
		"location_link":   location,
		"cert_cid":        certCID,
	})
}

func GetLogisticsByUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง User ID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [GetLogisticsByUser] Fetching logistics provider data for userID:", userID)

	// ✅ ค้นหา EntityID ของ User
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// ✅ ใช้ EntityID ค้นหาในตาราง Logistics
	var logistics models.Logistics
	if err := database.DB.Where("logisticsid = ?", user.EntityID).First(&logistics).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Logistics provider not found"})
	}

	// ✅ แยก areaCode และ phoneNumber ออกจาก Telephone
	areaCode, phoneNumber := utils.ExtractAreaCodeAndPhone(logistics.Telephone)

	// ✅ ส่งข้อมูลโลจิสติกส์กลับไป
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"logistics_id":  logistics.LogisticsID,
		"logisticsName": logistics.CompanyName,
		"address":       logistics.Address,
		"district":      logistics.District,
		"subdistrict":   logistics.SubDistrict,
		"province":      logistics.Province,
		"country":       logistics.Country,
		"post_code":     logistics.PostCode,
		"areaCode":      areaCode,    // ✅ รหัสพื้นที่
		"telephone":     phoneNumber, // ✅ หมายเลขโทรศัพท์
		"email":         logistics.Email,
		"walletAddress": logistics.WalletAddress,
		"location_link": logistics.LocationLink.String,
		"line_id":       logistics.LineID.String,   // ✅ เพิ่ม LineID
		"facebook":      logistics.Facebook.String, // ✅ เพิ่ม Facebook
		"created_on":    logistics.CreatedOn,
	})
}

func UpdateLogistics(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [UpdateLogistics] Updating logistics provider for userID:", userID)

	// ✅ ตรวจสอบว่า User มีโลจิสติกส์อยู่หรือไม่
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.EntityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User does not have a registered logistics provider"})
	}

	// ✅ ดึงข้อมูลโลจิสติกส์จาก EntityID ของ User
	var logistics models.Logistics
	if err := database.DB.Where("logisticsid = ?", user.EntityID).First(&logistics).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Logistics provider not found"})
	}

	// ✅ รับค่าจาก `FormData`
	logisticsName := strings.TrimSpace(c.FormValue("logisticsName"))
	address := strings.TrimSpace(c.FormValue("address"))
	district := strings.TrimSpace(c.FormValue("district"))
	subdistrict := strings.TrimSpace(c.FormValue("subdistrict"))
	province := strings.TrimSpace(c.FormValue("province"))
	postCode := strings.TrimSpace(c.FormValue("postCode"))
	phone := strings.TrimSpace(c.FormValue("phone"))
	areaCode := strings.TrimSpace(c.FormValue("areaCode"))
	location := strings.TrimSpace(c.FormValue("location_link"))
	lineID := strings.TrimSpace(c.FormValue("lineID"))
	facebook := strings.TrimSpace(c.FormValue("facebook"))

	// ✅ รวมเบอร์โทร
	fullPhone := fmt.Sprintf("%s %s", areaCode, phone)

	// ✅ ตรวจสอบว่ามีการเปลี่ยนแปลงหรือไม่
	updates := map[string]interface{}{}

	if logisticsName != "" && logisticsName != logistics.CompanyName {
		updates["companyname"] = logisticsName
	}
	if address != "" && address != logistics.Address {
		updates["address"] = address
	}
	if district != "" && district != logistics.District {
		updates["district"] = district
	}
	if subdistrict != "" && subdistrict != logistics.SubDistrict {
		updates["subdistrict"] = subdistrict
	}
	if province != "" && province != logistics.Province {
		updates["province"] = province
	}
	if postCode != "" && postCode != logistics.PostCode {
		updates["postcode"] = postCode
	}
	if fullPhone != "" && fullPhone != logistics.Telephone {
		updates["telephone"] = fullPhone
	}
	if location != "" && location != logistics.LocationLink.String {
		updates["location_link"] = location
	}
	if lineID != "" && lineID != logistics.LineID.String {
		updates["lineid"] = lineID
	}
	if facebook != "" && facebook != logistics.Facebook.String {
		updates["facebook"] = facebook
	}

	// ✅ ถ้าไม่มีการเปลี่ยนแปลง ให้แจ้งเตือน
	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No changes detected"})
	}

	// ✅ อัปเดตข้อมูลโลจิสติกส์
	if err := database.DB.Model(&models.Logistics{}).Where("logisticsid = ?", user.EntityID).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update logistics provider data"})
	}

	fmt.Println("✅ Logistics provider updated successfully:", user.EntityID)

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Logistics provider updated successfully",
		"logistics_id": user.EntityID,
	})
}
