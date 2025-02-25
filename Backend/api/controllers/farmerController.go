package controllers

import (
	"log"
	"math/big"
	"math/rand"

	"database/sql"
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ โครงสร้างข้อมูลที่รับจาก JSON Request
type FarmerRequest struct {
	UserID       string  `json:"userid"`
	CompanyName  string  `json:"company_name"`
	FirstName    string  `json:"firstname"`
	LastName     string  `json:"lastname"`
	Email        string  `json:"email"`
	Address      string  `json:"address"`
	AreaCode     string  `json:"areacode"`
	Phone        string  `json:"phone"`
	PostCode     string  `json:"postcode"`
	District     string  `json:"district"`
	SubDistrict  string  `json:"subdistrict"`
	Province     string  `json:"province"`
	Country      string  `json:"country"`
	LineID       *string `json:"lineid"`
	Facebook     *string `json:"facebook"`
	LocationLink *string `json:"location_link"`
}

// ✅ ใช้ Account จริงจาก Ganache แทนการสุ่ม
func getGanacheAccount() string {
	client, err := rpc.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Println("❌ Failed to connect to Ganache:", err)
		return "0x0000000000000000000000000000000000000000"
	}

	var accounts []common.Address
	err = client.Call(&accounts, "eth_accounts")
	if err != nil {
		log.Println("❌ Failed to get accounts from Ganache:", err)
		return "0x0000000000000000000000000000000000000000"
	}

	// ✅ เลือก Account ที่ยังไม่ถูกใช้
	selected := accounts[rand.Intn(len(accounts))] // สุ่ม 1 อันจาก Account ที่มีอยู่
	return selected.Hex()
}

func CreateFarmer(c *fiber.Ctx) error {
	var req FarmerRequest

	// ✅ แปลง JSON เป็น Struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// ✅ ตรวจสอบว่า User ID มีอยู่ในฐานข้อมูล `users` หรือไม่
	var user models.User
	if err := database.DB.Where("userid = ?", req.UserID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User ID not found"})
	}

	// ✅ ตรวจสอบว่าผู้ใช้เคยลงทะเบียนเป็น Farmer แล้วหรือไม่
	var existingFarmer models.Farmer
	err := database.DB.Where("userid = ?", req.UserID).First(&existingFarmer).Error
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already registered as a farmer"})
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ อัปเดต Role ของ User เป็น "farmer"
	if err := database.DB.Model(&models.User{}).Where("userid = ?", req.UserID).Update("role", "farmer").Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	// ✅ สร้าง FarmerID ใหม่ (FAYYNNNNN)
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('farmer_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate farmer ID"})
	}
	yearPrefix := time.Now().Format("06")
	farmerID := fmt.Sprintf("FA%s%05d", yearPrefix, sequence)

	// ✅ ดึง Wallet จริงจาก Ganache
	walletAddress := getGanacheAccount()
	fmt.Println("📌 DEBUG - Assigned Wallet Address:", walletAddress)

	// ✅ รวม `area code` กับ `phone`
	fullPhone := fmt.Sprintf("%s %s", strings.TrimSpace(req.AreaCode), strings.TrimSpace(req.Phone))

	// ✅ ตรวจสอบ `companyname` ถ้าว่างให้ใช้ "N/A"
	companyName := "N/A"
	if strings.TrimSpace(req.CompanyName) != "" {
		companyName = strings.TrimSpace(req.CompanyName)
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

	// ✅ สร้างข้อมูล Farmer
	farmer := models.Farmer{
		FarmerID:      farmerID,
		UserID:        req.UserID,
		FarmerName:    req.FirstName + " " + req.LastName,
		CompanyName:   companyName,
		Address:       strings.TrimSpace(req.Address),
		District:      strings.TrimSpace(req.District),    // ✅ เปลี่ยนจาก `city` เป็น `district`
		SubDistrict:   strings.TrimSpace(req.SubDistrict), // ✅ เพิ่ม `subdistrict`
		Province:      strings.TrimSpace(req.Province),
		Country:       strings.TrimSpace(req.Country),
		PostCode:      strings.TrimSpace(req.PostCode),
		Telephone:     fullPhone,
		LineID:        lineID,
		Facebook:      facebook,
		LocationLink:  locationLink,
		CreatedOn:     time.Now(),
		Email:         req.Email,
		WalletAddress: walletAddress, // ✅ ใช้ Wallet ที่ Generate อัตโนมัติ
	}

	// ✅ บันทึกลง Database
	if err := database.DB.Create(&farmer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save farmer data"})
	}

	// ✅ 🔗 ลงทะเบียนฟาร์มบน Blockchain
	txHash, err := services.BlockchainServiceInstance.RegisterFarmOnBlockchain(walletAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register farm on blockchain", "details": err.Error()})
	}

	fmt.Println("✅ Farmer Registered on Blockchain:", txHash)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Farmer registered successfully",
		"farmer_id":     farmerID,
		"walletAddress": walletAddress, // ✅ คืนค่า Wallet Address ให้ Frontend
		"txHash":        txHash,
	})
}

// ✅ ฟังก์ชันสำหรับดึงข้อมูลฟาร์มเมอร์ตาม UserID
func GetFarmerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var farmer models.Farmer

	// ค้นหาข้อมูล Farmer จาก userID
	if err := database.DB.Where("userid = ?", id).First(&farmer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farmer not found"})
	}

	// ✅ แยกชื่อ-นามสกุลออกจาก `FarmerName`
	nameParts := strings.SplitN(farmer.FarmerName, " ", 2)
	firstName := nameParts[0]
	lastName := ""
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}

	// ✅ แยก Area Code ออกจากเบอร์โทรศัพท์
	areaCode := "+66" // ค่าเริ่มต้น (ประเทศไทย)
	phoneNumber := farmer.Telephone

	if strings.HasPrefix(farmer.Telephone, "+") {
		parts := strings.SplitN(farmer.Telephone, " ", 2)
		if len(parts) == 2 {
			areaCode = parts[0]    // ดึงรหัสประเทศ
			phoneNumber = parts[1] // ดึงเบอร์โทรจริง
		}
	}

	// ✅ ตรวจสอบค่าว่างของ `sql.NullString`
	lineID := ""
	if farmer.LineID.Valid {
		lineID = farmer.LineID.String
	}

	facebook := ""
	if farmer.Facebook.Valid {
		facebook = farmer.Facebook.String
	}

	locationLink := ""
	if farmer.LocationLink.Valid {
		locationLink = farmer.LocationLink.String
	}

	// ✅ สร้าง JSON Response
	response := fiber.Map{
		"farmerID":    farmer.FarmerID,
		"userID":      farmer.UserID,
		"firstName":   firstName,
		"lastName":    lastName,
		"companyName": farmer.CompanyName,
		"address":     farmer.Address,
		"city":        farmer.District,
		"province":    farmer.Province,
		"country":     farmer.Country,
		"postCode":    farmer.PostCode,
		"areaCode":    areaCode,
		"telephone":   phoneNumber,
		"email":       farmer.Email,
		"wallet":      farmer.WalletAddress,
		"lineID":      lineID,
		"facebook":    facebook,
		"location":    locationLink,
	}

	return c.JSON(response)
}

// ✅ ดึงข้อมูลฟาร์มของ User ที่ล็อกอินอยู่
func GetFarmerByUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึงจาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	fmt.Println("🔍 [GetFarmerByUser] Fetching farmer for userID:", userID)

	var farmer models.Farmer
	if err := database.DB.Where("userid = ?", userID).First(&farmer).Error; err != nil {
		fmt.Println("❌ [GetFarmerByUser] Farmer not found for userID:", userID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farmer profile not found"})
	}

	fmt.Println("✅ [GetFarmerByUser] Farmer data found:", farmer.FarmerID)

	// ✅ แยกชื่อ-นามสกุลออกจาก `FarmerName`
	nameParts := strings.SplitN(farmer.FarmerName, " ", 2)
	firstName := nameParts[0]
	lastName := ""
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}

	// ✅ แยก Area Code ออกจากเบอร์โทรศัพท์
	areaCode := "+66"
	phoneNumber := farmer.Telephone
	if strings.HasPrefix(farmer.Telephone, "+") {
		parts := strings.SplitN(farmer.Telephone, " ", 2)
		if len(parts) == 2 {
			areaCode = parts[0]
			phoneNumber = parts[1]
		}
	}

	// ✅ ตรวจสอบค่าว่างของ `sql.NullString`
	lineID := ""
	if farmer.LineID.Valid {
		lineID = farmer.LineID.String
	}

	facebook := ""
	if farmer.Facebook.Valid {
		facebook = farmer.Facebook.String
	}

	locationLink := ""
	if farmer.LocationLink.Valid {
		locationLink = farmer.LocationLink.String
	}

	// ✅ สร้าง JSON Response
	response := fiber.Map{
		"farmerID":    farmer.FarmerID,
		"userID":      farmer.UserID,
		"firstName":   firstName,
		"lastName":    lastName,
		"companyName": farmer.CompanyName,
		"address":     farmer.Address,
		"district":    farmer.District,
		"subdistrict": farmer.SubDistrict,
		"province":    farmer.Province,
		"country":     farmer.Country,
		"postCode":    farmer.PostCode,
		"areaCode":    areaCode,
		"telephone":   phoneNumber,
		"email":       farmer.Email,
		"wallet":      farmer.WalletAddress,
		"lineID":      lineID,
		"facebook":    facebook,
		"location":    locationLink,
	}

	return c.JSON(response)
}

func UpdateFarmer(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [UpdateFarmer] Fetching farmer for userID:", userID)

	var farmer models.Farmer
	if err := database.DB.Where("userid = ?", userID).First(&farmer).Error; err != nil {
		fmt.Println("❌ [UpdateFarmer] Farmer not found for userID:", userID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farmer profile not found"})
	}

	// โครงสร้าง request สำหรับรับข้อมูลแบบ snake_case
	var req struct {
		FirstName    string  `json:"firstname"`
		LastName     string  `json:"lastname"`
		CompanyName  string  `json:"company_name"`
		Address      string  `json:"address"`
		District     string  `json:"district"`
		SubDistrict  string  `json:"subdistrict"`
		Province     string  `json:"province"`
		Country      string  `json:"country"`
		PostCode     string  `json:"postcode"`
		AreaCode     string  `json:"area_code"`
		Phone        string  `json:"phone"`
		LineID       *string `json:"lineid"`
		Facebook     *string `json:"facebook"`
		LocationLink *string `json:"location_link"`
		CertFile     string  `json:"cert_file"` // ไฟล์ใบเซอร์ (ในรูปแบบ Base64) จาก frontend
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// ประกาศตัวแปรสำหรับเก็บ CID ของใบรับรอง
	var certCID string
	var err error

	if req.CertFile != "" {
		fmt.Println("📌 Uploading new certification file to IPFS...")
		// อัปโหลดไฟล์ใบรับรอง (Base64) ไปยัง IPFS
		certCID, err = ipfsService.UploadBase64File(req.CertFile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload certification file to IPFS"})
		}
		fmt.Println("✅ Certification file uploaded to IPFS with CID:", certCID)
	} else {
		// ดึง CID เดิมจาก Blockchain (ใช้ event ID ที่อิงกับ FarmerID)
		eventID := fmt.Sprintf("EVENT-%s", farmer.FarmerID)
		certification, err := services.BlockchainServiceInstance.GetCertificationFromBlockchain(eventID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch certification CID from blockchain"})
		}
		// สมมติว่า models.Certification มีฟิลด์ CertificationCID เป็น string
		certCID = certification.CertificationCID
		fmt.Println("📌 Using existing Certification CID:", certCID)
	}

	// รวม Area Code กับ Phone ให้เป็นเบอร์โทรศัพท์ที่สมบูรณ์
	fullPhone := fmt.Sprintf("%s %s", req.AreaCode, req.Phone)

	// สร้างข้อมูล Farmer ใหม่เพื่ออัปเดต (แมปจาก snake_case payload)
	updatedFarmer := models.Farmer{
		// คงค่า FarmerID, UserID, Email, WalletAddress จากข้อมูลเดิม
		FarmerID: farmer.FarmerID,
		UserID:   farmer.UserID,
		// รวม first name กับ last name ให้เป็น FarmerName
		FarmerName:  fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		CompanyName: req.CompanyName,
		Address:     req.Address,
		// สมมติว่าฟิลด์ District และ SubDistrict มีอยู่ใน models.Farmer
		District:      req.District,
		SubDistrict:   req.SubDistrict,
		Province:      req.Province,
		Country:       req.Country,
		PostCode:      req.PostCode,
		Telephone:     fullPhone,
		Email:         farmer.Email,
		WalletAddress: farmer.WalletAddress,
		LineID: sql.NullString{String: func() string {
			if req.LineID != nil {
				return *req.LineID
			} else {
				return ""
			}
		}(), Valid: req.LineID != nil},
		Facebook: sql.NullString{String: func() string {
			if req.Facebook != nil {
				return *req.Facebook
			} else {
				return ""
			}
		}(), Valid: req.Facebook != nil},
		LocationLink: sql.NullString{String: func() string {
			if req.LocationLink != nil {
				return *req.LocationLink
			} else {
				return ""
			}
		}(), Valid: req.LocationLink != nil},
	}

	// อัปเดตข้อมูลลง PostgreSQL
	if err := database.DB.Model(&farmer).Updates(updatedFarmer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update farm information"})
	}

	// ถ้ามีการอัปโหลดใบรับรองใหม่ ให้บันทึกลง Blockchain
	if req.CertFile != "" {
		txHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(
			fmt.Sprintf("EVENT-%s", farmer.FarmerID),
			"Farmer",
			farmer.FarmerID,
			certCID,
			big.NewInt(time.Now().Unix()), // วันที่ออกใบเซอร์
			big.NewInt(time.Now().AddDate(1, 0, 0).Unix()), // วันหมดอายุ (1 ปี)
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update certification on blockchain"})
		}
		fmt.Println("✅ Certification updated on Blockchain:", txHash)
	}

	return c.JSON(fiber.Map{
		"message": "Farm information updated successfully!",
		"certCID": certCID,
	})
}
