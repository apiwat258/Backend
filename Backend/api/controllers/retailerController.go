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

func CreateRetailer(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [CreateRetailer] Creating retailer for userID:", userID)

	// ✅ รับค่าจาก `FormData`
	companyName := strings.TrimSpace(c.FormValue("retailerName"))
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
	if user.Role == "retailer" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already has a retailer role"})
	}

	// ✅ ตรวจสอบ Email ซ้ำ
	var existingRetailer models.Retailer
	if err := database.DB.Where("email = ?", email).First(&existingRetailer).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Retailer email is already in use"})
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

	// ✅ สร้าง `retailerID`
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('retailer_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate retailer ID"})
	}
	yearPrefix := time.Now().Format("06")
	retailerID := fmt.Sprintf("RET%s%05d", yearPrefix, sequence)

	// ✅ ดึง Wallet จาก Ganache
	walletAddress := getGanacheAccount()

	// ✅ ลงทะเบียน User บน Blockchain (ถ้ายังไม่ได้ลงทะเบียน)
	txHash, err := services.BlockchainServiceInstance.RegisterUserOnBlockchain(walletAddress, 4) // 4 = Retailer Role
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user on blockchain"})
	}
	fmt.Println("✅ User registered on Blockchain. Transaction Hash:", txHash)

	// ✅ ถ้ามีใบเซอร์ → บันทึกลง Blockchain
	if certCID != "" {
		// ✅ สร้าง `eventID`
		eventID := fmt.Sprintf("EVENT-%s-%s", retailerID, uuid.New().String())

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
			"retailer",
			retailerID,
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
	retailer := models.Retailer{
		RetailerID:    retailerID,
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
		LocationLink:  location,
		LineID:        sql.NullString{String: lineID, Valid: lineID != ""},
		Facebook:      sql.NullString{String: facebook, Valid: facebook != ""},
		CreatedOn:     time.Now(),
	}

	if err := database.DB.Create(&retailer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save retailer data"})
	}

	// ✅ อัปเดต `entityID` และ Role ใน `users`
	updateData := map[string]interface{}{
		"entityid": retailerID,
		"role":     "retailer",
	}
	if err := database.DB.Model(&models.User{}).Where("userid = ?", userID).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Retailer registered successfully",
		"retailer_id":    retailerID,
		"retailer_email": email,
		"walletAddress":  walletAddress,
		"location_link":  location,
		"cert_cid":       certCID,
	})
}

func GetRetailerByUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง User ID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [GetRetailerByUser] Fetching retailer data for userID:", userID)

	// ✅ ค้นหา EntityID ของ User
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// ✅ ใช้ EntityID ค้นหาในตาราง Retailer
	var retailer models.Retailer
	if err := database.DB.Where("retailerid = ?", user.EntityID).First(&retailer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Retailer not found"})
	}

	// ✅ แยก areaCode และ phoneNumber ออกจาก Telephone
	areaCode, phoneNumber := utils.ExtractAreaCodeAndPhone(retailer.Telephone)

	// ✅ ส่งข้อมูลร้านค้ากลับไป
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"retailer_id":   retailer.RetailerID,
		"retailerName":  retailer.CompanyName,
		"address":       retailer.Address,
		"district":      retailer.District,
		"subdistrict":   retailer.SubDistrict,
		"province":      retailer.Province,
		"country":       retailer.Country,
		"post_code":     retailer.PostCode,
		"areaCode":      areaCode,    // ✅ รหัสพื้นที่
		"telephone":     phoneNumber, // ✅ หมายเลขโทรศัพท์
		"email":         retailer.Email,
		"walletAddress": retailer.WalletAddress,
		"location_link": retailer.LocationLink,
		"line_id":       retailer.LineID.String,   // ✅ เพิ่ม LineID
		"facebook":      retailer.Facebook.String, // ✅ เพิ่ม Facebook
		"created_on":    retailer.CreatedOn,
	})
}

func UpdateRetailer(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [UpdateRetailer] Updating retailer for userID:", userID)

	// ✅ ตรวจสอบว่า User มีร้านค้าอยู่หรือไม่
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.EntityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User does not have a registered retailer"})
	}

	// ✅ ดึงข้อมูลร้านค้าจาก EntityID ของ User
	var retailer models.Retailer
	if err := database.DB.Where("retailerid = ?", user.EntityID).First(&retailer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Retailer not found"})
	}

	// ✅ รับค่าจาก `FormData`
	companyName := strings.TrimSpace(c.FormValue("retailerName"))
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

	if companyName != "" && companyName != retailer.CompanyName {
		updates["companyname"] = companyName
	}
	if address != "" && address != retailer.Address {
		updates["address"] = address
	}
	if district != "" && district != retailer.District {
		updates["district"] = district
	}
	if subdistrict != "" && subdistrict != retailer.SubDistrict {
		updates["subdistrict"] = subdistrict
	}
	if province != "" && province != retailer.Province {
		updates["province"] = province
	}
	if postCode != "" && postCode != retailer.PostCode {
		updates["postcode"] = postCode
	}
	if fullPhone != "" && fullPhone != retailer.Telephone {
		updates["telephone"] = fullPhone
	}
	if location != "" && location != retailer.LocationLink {
		updates["location_link"] = location
	}
	if lineID != "" && lineID != retailer.LineID.String {
		updates["lineid"] = lineID
	}
	if facebook != "" && facebook != retailer.Facebook.String {
		updates["facebook"] = facebook
	}

	// ✅ ถ้าไม่มีการเปลี่ยนแปลง ให้แจ้งเตือน
	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No changes detected"})
	}

	// ✅ อัปเดตข้อมูลร้านค้า
	if err := database.DB.Model(&models.Retailer{}).Where("retailerid = ?", user.EntityID).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update retailer data"})
	}

	fmt.Println("✅ Retailer updated successfully:", user.EntityID)

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Retailer updated successfully",
		"retailer_id": user.EntityID,
	})
}

// GetAllRetailers ดึง retailerID และชื่อร้านค้าทั้งหมด
func GetAllRetailers(c *fiber.Ctx) error {
	var retailers []models.Retailer // ✅ ใช้ Model เต็ม

	// ✅ ลอง Query โดยใช้ Model เต็ม
	result := database.DB.Find(&retailers)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch retailers"})
	}

	// ✅ Debug Log เช็คค่าที่ GORM ดึงออกมา
	fmt.Println("📡 Query Result:", retailers)

	// ✅ ถ้าไม่มีข้อมูล ให้แจ้งเตือน
	if len(retailers) == 0 {
		fmt.Println("⚠️ No retailers found in database")
		return c.JSON([]models.Retailer{})
	}

	// ✅ สร้าง Array ใหม่ที่มีเฉพาะ `retailer_id` และ `company_name`
	var simplifiedRetailers []struct {
		RetailerID  string `json:"retailer_id"`
		CompanyName string `json:"company_name"`
	}

	for _, retailer := range retailers {
		simplifiedRetailers = append(simplifiedRetailers, struct {
			RetailerID  string `json:"retailer_id"`
			CompanyName string `json:"company_name"`
		}{
			RetailerID:  retailer.RetailerID,
			CompanyName: retailer.CompanyName,
		})
	}

	// ✅ Debug Log เช็คค่าก่อนส่งออกไป
	fmt.Println("📡 Simplified Query Result:", simplifiedRetailers)

	// ✅ ส่งข้อมูลกลับไป
	return c.JSON(simplifiedRetailers)
}

// GetRetailerByID ดึงข้อมูลร้านค้าตาม retailerID
func GetRetailerByID(c *fiber.Ctx) error {
	retailerID := c.Params("id")

	var retailer models.Retailer
	if err := database.DB.Where("retailerid = ?", retailerID).First(&retailer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Retailer not found"})
	}

	return c.JSON(fiber.Map{
		"retailer_id":   retailer.RetailerID,
		"company_name":  retailer.CompanyName,
		"email":         retailer.Email,
		"telephone":     retailer.Telephone,
		"address":       retailer.Address,
		"province":      retailer.Province,
		"district":      retailer.District,
		"subdistrict":   retailer.SubDistrict,
		"post_code":     retailer.PostCode,
		"location_link": retailer.LocationLink,
	})
}

// GetRetailerUsernames ดึง username ทั้งหมดของร้านค้า
func GetRetailerUsernames(c *fiber.Ctx) error {
	// ✅ ดึงค่า retailer_id จาก query parameter
	retailerID := c.Query("retailer_id")               // <-- ตรวจสอบว่าค่านี้ถูกต้องจริงๆ
	fmt.Println("📌 Received retailer_id:", retailerID) // ✅ Debug Log

	// ✅ ถ้าไม่มี retailer_id ให้ return error
	if retailerID == "" {
		fmt.Println("❌ Missing retailer_id") // ✅ Debug Log
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing retailer_id"})
	}

	// ✅ ตรวจสอบว่ามี retailer_id จริงใน database
	var count int64
	if err := database.DB.Model(&models.Retailer{}).Where("retailerid = ?", retailerID).Count(&count).Error; err != nil {
		fmt.Println("❌ Database Error:", err) // ✅ Debug Log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ ถ้าไม่พบ retailer_id
	if count == 0 {
		fmt.Println("❌ Retailer not found in database for ID:", retailerID) // ✅ Debug Log
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Retailer not found"})
	}

	// ✅ Query หาข้อมูล users ที่ entityid ตรงกับ retailer_id
	var users []struct {
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		EntityID  string `json:"entity_id"`
	}

	if err := database.DB.Raw(`
		SELECT username, 
		       SPLIT_PART(username, ' ', 1) AS first_name, 
		       SPLIT_PART(username, ' ', 2) AS last_name, 
		       entityid AS entity_id 
		FROM users WHERE role = 'retailer' AND entityid = ?
	`, retailerID).Scan(&users).Error; err != nil {
		fmt.Println("❌ Query Error:", err) // ✅ Debug Log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch usernames"})
	}

	// ✅ Debug Log เช็คค่าที่ดึงออกมา
	fmt.Println("📡 Query Result:", users)

	// ✅ ส่งข้อมูลกลับไป
	return c.JSON(users)
}
