package controllers

import (
	"log"
	"math/big"
	"math/rand"

	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"
	"finalyearproject/Backend/utils"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
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
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [CreateFarmer] Creating farm for userID:", userID)

	// ✅ รับค่าจาก `FormData`
	farmName := strings.TrimSpace(c.FormValue("farmName"))
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

	// ✅ รวมเบอร์โทร
	fullPhone := fmt.Sprintf("%s %s", areaCode, phone)

	// ✅ ตรวจสอบ Role ของ User
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.Role == "farmer" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already has a farmer role"})
	}

	// ✅ ตรวจสอบ Email ซ้ำ
	var existingFarmer models.Farmer
	if err := database.DB.Where("email = ?", email).First(&existingFarmer).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Farm email is already in use"})
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

	// ✅ สร้าง `farmerID`
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('farmer_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate farmer ID"})
	}
	yearPrefix := time.Now().Format("06")
	farmerID := fmt.Sprintf("FA%s%05d", yearPrefix, sequence)

	// ✅ ดึง Wallet จาก Ganache
	walletAddress := getGanacheAccount()

	// ✅ ลงทะเบียน User บน Blockchain (ถ้ายังไม่ได้ลงทะเบียน)
	txHash, err := services.BlockchainServiceInstance.RegisterUserOnBlockchain(walletAddress, 1) // 1 = Farmer Role
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user on blockchain"})
	}
	fmt.Println("✅ User registered on Blockchain. Transaction Hash:", txHash)

	// ✅ ถ้ามีใบเซอร์ → บันทึกลง Blockchain
	if certCID != "" {
		// ✅ สร้าง `eventID`
		eventID := fmt.Sprintf("EVENT-%s-%s", farmerID, uuid.New().String())

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
			"farmer",
			farmerID,
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
	farmer := models.Farmer{
		FarmerID:      farmerID,
		CompanyName:   farmName,
		Address:       address,
		District:      district,
		SubDistrict:   subdistrict,
		Province:      province,
		PostCode:      postCode,
		Telephone:     fullPhone,
		Email:         email,
		WalletAddress: walletAddress,
		LocationLink:  location,
		CreatedOn:     time.Now(),
	}

	if err := database.DB.Create(&farmer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save farmer data"})
	}

	// ✅ อัปเดต `entityID` และ Role ใน `users`
	updateData := map[string]interface{}{
		"entityid": farmerID,
		"role":     "farmer",
	}
	if err := database.DB.Model(&models.User{}).Where("userid = ?", userID).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Farmer registered successfully",
		"farmer_id":     farmerID,
		"farm_email":    email,
		"walletAddress": walletAddress,
		"location_link": location,
		"cert_cid":      certCID, // ✅ ส่ง CID ของใบเซอร์กลับไป (ถ้ามี)
	})
}

// ✅ API: ดึงข้อมูลฟาร์มจาก Entity ID ของ User
func GetFarmByUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง User ID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [GetFarmByUser] Fetching farm data for userID:", userID)

	// ✅ ค้นหา EntityID ของ User
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// ✅ ใช้ EntityID ค้นหาในตาราง Farmer
	var farmer models.Farmer
	if err := database.DB.Where("farmerid = ?", user.EntityID).First(&farmer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farmer not found"})
	}

	// ✅ แยก areaCode และ phoneNumber ออกจาก Telephone
	areaCode, phoneNumber := utils.ExtractAreaCodeAndPhone(farmer.Telephone)

	// ✅ ส่งข้อมูลฟาร์มกลับไป
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"farmer_id":     farmer.FarmerID,
		"farm_name":     farmer.CompanyName,
		"address":       farmer.Address,
		"district":      farmer.District,
		"subdistrict":   farmer.SubDistrict,
		"province":      farmer.Province,
		"post_code":     farmer.PostCode,
		"areaCode":      areaCode,    // ✅ รหัสพื้นที่
		"telephone":     phoneNumber, // ✅ หมายเลขโทรศัพท์
		"email":         farmer.Email,
		"walletAddress": farmer.WalletAddress,
		"location_link": farmer.LocationLink,
		"created_on":    farmer.CreatedOn,
	})
}

func UpdateFarmer(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [UpdateFarmer] Updating farm for userID:", userID)

	// ✅ ตรวจสอบว่า User มีฟาร์มอยู่หรือไม่
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.EntityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User does not have a registered farm"})
	}

	// ✅ ดึงฟาร์มจาก EntityID ของ User
	var farmer models.Farmer
	if err := database.DB.Where("farmerid = ?", user.EntityID).First(&farmer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farm not found"})
	}

	// ✅ รับค่าจาก `FormData`
	farmName := strings.TrimSpace(c.FormValue("farmName"))
	address := strings.TrimSpace(c.FormValue("address"))
	district := strings.TrimSpace(c.FormValue("district"))
	subdistrict := strings.TrimSpace(c.FormValue("subdistrict"))
	province := strings.TrimSpace(c.FormValue("province"))
	postCode := strings.TrimSpace(c.FormValue("postCode"))
	phone := strings.TrimSpace(c.FormValue("phone"))
	areaCode := strings.TrimSpace(c.FormValue("areaCode"))
	location := strings.TrimSpace(c.FormValue("location_link")) // ✅ คงค่า location ไว้ใช้งาน

	// ✅ รวมเบอร์โทร
	fullPhone := fmt.Sprintf("%s %s", areaCode, phone)

	// ✅ ตรวจสอบว่ามีการเปลี่ยนแปลงหรือไม่
	updates := map[string]interface{}{}

	if farmName != "" && farmName != farmer.CompanyName {
		updates["companyname"] = farmName
	}
	if address != "" && address != farmer.Address {
		updates["address"] = address
	}
	if district != "" && district != farmer.District {
		updates["district"] = district
	}
	if subdistrict != "" && subdistrict != farmer.SubDistrict {
		updates["subdistrict"] = subdistrict
	}
	if province != "" && province != farmer.Province {
		updates["province"] = province
	}
	if postCode != "" && postCode != farmer.PostCode {
		updates["postcode"] = postCode
	}
	if fullPhone != "" && fullPhone != farmer.Telephone {
		updates["telephone"] = fullPhone
	}
	if location != "" && location != farmer.LocationLink {
		updates["location_link"] = location
	}

	// ✅ ถ้าไม่มีการเปลี่ยนแปลง ให้แจ้งเตือน
	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No changes detected"})
	}

	// ✅ อัปเดตข้อมูลฟาร์ม
	if err := database.DB.Model(&models.Farmer{}).Where("farmerid = ?", user.EntityID).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update farm data"})
	}

	fmt.Println("✅ Farm updated successfully:", user.EntityID)

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Farm updated successfully",
		"farmer_id": user.EntityID,
	})
}
