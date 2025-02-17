package controllers

import (
	"log"
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

	// ✅ แปลง `*string` เป็น `sql.NullString`
	email := sql.NullString{}
	if strings.TrimSpace(req.Email) != "" {
		email = sql.NullString{String: strings.TrimSpace(req.Email), Valid: true}
	}

	// ✅ สร้างข้อมูล Farmer
	farmer := models.Farmer{
		FarmerID:      farmerID,
		UserID:        req.UserID,
		FarmerName:    req.FirstName + " " + req.LastName,
		CompanyName:   companyName,
		Address:       fullAddress,
		City:          req.City,
		Province:      province,
		Country:       req.Country,
		PostCode:      req.PostCode,
		Telephone:     fullPhone,
		CreatedOn:     time.Now(),
		Email:         email.String,
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
		"city":        farmer.City,
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
		"city":        farmer.City,
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

// ✅ อัปเดตข้อมูลฟาร์ม
func UpdateFarmer(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var updatedFarmer models.Farmer
	if err := c.BodyParser(&updatedFarmer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
	}

	// ค้นหาข้อมูลฟาร์มเดิม
	var farmer models.Farmer
	if err := database.DB.Where("user_id = ?", userID).First(&farmer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farmer profile not found"})
	}

	// อัปเดตข้อมูล
	database.DB.Model(&farmer).Updates(updatedFarmer)

	return c.JSON(fiber.Map{"message": "Farm information updated successfully!"})
}
