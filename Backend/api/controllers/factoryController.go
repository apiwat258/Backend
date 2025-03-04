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

func CreateFactory(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [CreateFactory] Creating factory for userID:", userID)

	// ✅ รับค่าจาก `FormData`
	companyName := strings.TrimSpace(c.FormValue("factoryName"))
	email := strings.TrimSpace(c.FormValue("email"))
	address := strings.TrimSpace(c.FormValue("address"))
	district := strings.TrimSpace(c.FormValue("district"))
	subdistrict := strings.TrimSpace(c.FormValue("subdistrict"))
	province := strings.TrimSpace(c.FormValue("province"))
	postCode := strings.TrimSpace(c.FormValue("postCode"))
	phone := strings.TrimSpace(c.FormValue("phone"))
	areaCode := strings.TrimSpace(c.FormValue("areaCode"))
	location := strings.TrimSpace(c.FormValue("location_link")) // ✅ คงค่า location ไว้ใช้งาน
	certCID := strings.TrimSpace(c.FormValue("cert_cid"))       // ✅ รับค่า certCID จาก Frontend
	lineID := strings.TrimSpace(c.FormValue("lineID"))
	facebook := strings.TrimSpace(c.FormValue("facebook"))

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
	factoryID := fmt.Sprintf("FAC%s%05d", yearPrefix, sequence)

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
			walletAddress,
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
		"cert_cid":      certCID,
	})
}

// ✅ API: ดึงข้อมูลโรงงานจาก Entity ID ของ User
func GetFactoryByUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง User ID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [GetFactoryByUser] Fetching factory data for userID:", userID)

	// ✅ ค้นหา EntityID ของ User
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// ✅ ใช้ EntityID ค้นหาในตาราง Factory
	var factory models.Factory
	if err := database.DB.Where("factoryid = ?", user.EntityID).First(&factory).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Factory not found"})
	}

	// ✅ แยก areaCode และ phoneNumber ออกจาก Telephone
	areaCode, phoneNumber := utils.ExtractAreaCodeAndPhone(factory.Telephone)

	// ✅ ส่งข้อมูลโรงงานกลับไป
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"factory_id":    factory.FactoryID,
		"factoryName":   factory.CompanyName,
		"address":       factory.Address,
		"district":      factory.District,
		"subdistrict":   factory.SubDistrict,
		"province":      factory.Province,
		"country":       factory.Country,
		"post_code":     factory.PostCode,
		"areaCode":      areaCode,    // ✅ รหัสพื้นที่
		"telephone":     phoneNumber, // ✅ หมายเลขโทรศัพท์
		"email":         factory.Email,
		"walletAddress": factory.WalletAddress,
		"location_link": factory.LocationLink.String,
		"line_id":       factory.LineID.String,   // ✅ เพิ่ม LineID
		"facebook":      factory.Facebook.String, // ✅ เพิ่ม Facebook
		"created_on":    factory.CreatedOn,
	})
}

func UpdateFactory(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [UpdateFactory] Updating factory for userID:", userID)

	// ✅ ตรวจสอบว่า User มีโรงงานอยู่หรือไม่
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.EntityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User does not have a registered factory"})
	}

	// ✅ ดึงข้อมูลโรงงานจาก EntityID ของ User
	var factory models.Factory
	if err := database.DB.Where("factoryid = ?", user.EntityID).First(&factory).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Factory not found"})
	}

	// ✅ รับค่าจาก `FormData`
	companyName := strings.TrimSpace(c.FormValue("factoryName"))
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

	if companyName != "" && companyName != factory.CompanyName {
		updates["companyname"] = companyName
	}
	if address != "" && address != factory.Address {
		updates["address"] = address
	}
	if district != "" && district != factory.District {
		updates["district"] = district
	}
	if subdistrict != "" && subdistrict != factory.SubDistrict {
		updates["subdistrict"] = subdistrict
	}
	if province != "" && province != factory.Province {
		updates["province"] = province
	}
	if postCode != "" && postCode != factory.PostCode {
		updates["postcode"] = postCode
	}
	if fullPhone != "" && fullPhone != factory.Telephone {
		updates["telephone"] = fullPhone
	}
	if location != "" && location != factory.LocationLink.String {
		updates["location_link"] = location
	}
	if lineID != "" && lineID != factory.LineID.String {
		updates["lineid"] = lineID
	}
	if facebook != "" && facebook != factory.Facebook.String {
		updates["facebook"] = facebook
	}

	// ✅ ถ้าไม่มีการเปลี่ยนแปลง ให้แจ้งเตือน
	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No changes detected"})
	}

	// ✅ อัปเดตข้อมูลโรงงาน
	if err := database.DB.Model(&models.Factory{}).Where("factoryid = ?", user.EntityID).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update factory data"})
	}

	fmt.Println("✅ Factory updated successfully:", user.EntityID)

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Factory updated successfully",
		"factory_id": user.EntityID,
	})
}
