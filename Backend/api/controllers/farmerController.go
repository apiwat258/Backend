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

	// ✅ ตรวจสอบ Role (ต้องไม่มี Role หรือไม่ใช่ `"farmer"`)
	if user.Role == "farmer" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already has a farmer role"})
	}

	// ✅ ตรวจสอบว่าอีเมลของฟาร์มซ้ำหรือไม่
	var existingFarmer models.Farmer
	if err := database.DB.Where("email = ?", req.Email).First(&existingFarmer).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Farm email is already in use"})
	}

	// ✅ สร้าง FarmerID ใหม่ (FAYYNNNNN)
	var sequence int64
	if err := database.DB.Raw("SELECT nextval('farmer_id_seq')").Scan(&sequence).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate farmer ID"})
	}
	yearPrefix := time.Now().Format("06")
	farmerID := fmt.Sprintf("FA%s%05d", yearPrefix, sequence) // ✅ ใช้ `farmerID` เป็น `entityID`

	// ✅ ดึง Wallet จริงจาก Ganache
	walletAddress := getGanacheAccount()

	// ✅ รวม `area code` กับ `phone`
	fullPhone := fmt.Sprintf("%s %s", strings.TrimSpace(req.AreaCode), strings.TrimSpace(req.Phone))

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
		EntityID:      farmerID, // ✅ ใช้ FarmerID เป็น EntityID
		CompanyName:   strings.TrimSpace(req.CompanyName),
		Address:       strings.TrimSpace(req.Address),
		District:      strings.TrimSpace(req.District),
		SubDistrict:   strings.TrimSpace(req.SubDistrict),
		Province:      strings.TrimSpace(req.Province),
		Country:       strings.TrimSpace(req.Country),
		PostCode:      strings.TrimSpace(req.PostCode),
		Telephone:     fullPhone,
		Email:         req.Email, // ✅ เก็บอีเมลของฟาร์ม
		LineID:        lineID,
		Facebook:      facebook,
		LocationLink:  locationLink,
		CreatedOn:     time.Now(),
		WalletAddress: walletAddress,
	}

	// ✅ บันทึกลง Database
	if err := database.DB.Create(&farmer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save farmer data"})
	}

	// ✅ 🔗 อัปเดต `entityID` และ Role ในตาราง `users`
	updateData := map[string]interface{}{
		"entityid": farmerID,
		"role":     "farmer",
	}
	if err := database.DB.Model(&models.User{}).Where("userid = ?", req.UserID).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role & entity ID"})
	}

	// ✅ ลงทะเบียน User ลง Smart Contract กลาง (UserRegistry)
	txHash, err := services.BlockchainServiceInstance.RegisterUserOnBlockchain(walletAddress, 1) // 1 = Farmer Role
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user on blockchain", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Farmer registered successfully",
		"farmer_id":     farmerID,  // ✅ `entityID` = `farmerID`
		"farm_email":    req.Email, // ✅ คืนค่า Email ของฟาร์ม
		"walletAddress": walletAddress,
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
	//nameParts := strings.SplitN(//farmer.//FarmerName, " ", 2)
	//firstName := nameParts[0]
	//lastName := ""
	//if len(nameParts) > 1 {
	//	lastName = nameParts[1]
	//}

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
		"farmerID": farmer.FarmerID,
		//"userID":      farmer.UserID,
		//"firstName":   firstName,
		//"lastName":    lastName,
		"companyName": farmer.CompanyName,
		"address":     farmer.Address,
		"city":        farmer.District,
		"province":    farmer.Province,
		"country":     farmer.Country,
		"postCode":    farmer.PostCode,
		"areaCode":    areaCode,
		"telephone":   phoneNumber,
		//"email":       farmer.Email,
		"wallet":   farmer.WalletAddress,
		"lineID":   lineID,
		"facebook": facebook,
		"location": locationLink,
	}

	return c.JSON(response)
}

// ✅ ดึงข้อมูลฟาร์มของ User ที่ล็อกอินอยู่
func GetFarmerByUser(c *fiber.Ctx) error {
    userID, ok := c.Locals("userID").(string) // ✅ ดึง userID จาก Middleware
    if !ok || userID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }
    fmt.Println("🔍 [GetFarmerByUser] Fetching entityID for userID:", userID)

    // ✅ ค้นหา `entityID` ในตาราง Users ก่อน
    var user models.User
    if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
        fmt.Println("❌ [GetFarmerByUser] User not found:", userID)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }

    // ✅ ตรวจสอบ Role ก่อนดึงข้อมูลฟาร์ม
    if user.Role != "farmer" {
        fmt.Println("⚠️ [GetFarmerByUser] User is not a farmer:", userID)
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not a farmer"})
    }

    // ✅ ค้นหา `farmer` โดยใช้ `entityID`
    var farmer models.Farmer
    if err := database.DB.Where("farmerid = ?", user.EntityID).First(&farmer).Error; err != nil {
        fmt.Println("❌ [GetFarmerByUser] Farmer not found for entityID:", user.EntityID)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farmer profile not found"})
    }

    fmt.Println("✅ [GetFarmerByUser] Farmer data found:", farmer.FarmerID)

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

    // ✅ ส่ง JSON Response
    response := fiber.Map{
        "farmerID":   farmer.FarmerID,
        "address":    farmer.Address,
        "district":   farmer.District,
        "subdistrict": farmer.SubDistrict,
        "province":   farmer.Province,
        "country":    farmer.Country,
        "postCode":   farmer.PostCode,
        "areaCode":   areaCode,
        "telephone":  phoneNumber,
        "wallet":     farmer.WalletAddress,
        "lineID":     lineID,
        "facebook":   facebook,
        "location":   locationLink,
    }

    return c.JSON(response)
}


func UpdateFarmer(c *fiber.Ctx) error {
    userID, ok := c.Locals("userID").(string)
    if !ok || userID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    fmt.Println("🔍 [UpdateFarmer] Fetching entityID for userID:", userID)

    // ✅ ค้นหา `entityID` ในตาราง Users ก่อน
    var user models.User
    if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
        fmt.Println("❌ [UpdateFarmer] User not found:", userID)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }

    // ✅ ตรวจสอบ Role ก่อนอัปเดตข้อมูลฟาร์ม
    if user.Role != "farmer" {
        fmt.Println("⚠️ [UpdateFarmer] User is not a farmer:", userID)
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not a farmer"})
    }

    // ✅ ค้นหา `farmer` โดยใช้ `entityID`
    var farmer models.Farmer
    if err := database.DB.Where("farmerid = ?", user.EntityID).First(&farmer).Error; err != nil {
        fmt.Println("❌ [UpdateFarmer] Farmer not found for entityID:", user.EntityID)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farmer profile not found"})
    }

    fmt.Println("✅ [UpdateFarmer] Farmer data found:", farmer.FarmerID)

    // ✅ อ่านข้อมูลใหม่จาก Request Body
    var req struct {
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
        CertFile     string  `json:"cert_file"` // ✅ ฟิลด์ใบเซอร์
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
    }

    // ✅ Debug: ดูค่า `CertFile`
    fmt.Println("📌 [UpdateFarmer] Received CertFile:", req.CertFile)

    // ✅ ตรวจสอบว่ามีใบเซอร์อยู่แล้วหรือไม่
    var latestCertCID string
    existingCert, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(farmer.FarmerID)
    if err == nil && len(existingCert) > 0 {
        for _, cert := range existingCert {
            if cert.IsActive {
                latestCertCID = cert.CertificationCID
                break
            }
        }
    }

    // ✅ ตรวจสอบว่า `cert_file` มีการเปลี่ยนแปลงหรือไม่
    var certCID string = latestCertCID
    if req.CertFile != "" && req.CertFile != latestCertCID {
        if strings.HasPrefix(req.CertFile, "Qm") {
            // ✅ เป็น CID อยู่แล้ว → ใช้ค่าที่ได้รับมา
            certCID = req.CertFile
        } else if strings.HasPrefix(req.CertFile, "data:") {
            // ✅ เป็น Base64 → อัปโหลดไป IPFS
            certCID, err = ipfsService.UploadBase64File(req.CertFile)
            if err != nil || certCID == "" {
                return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload certification file to IPFS"})
            }
        } else {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cert_file format"})
        }
    }

    // ✅ รวม Area Code กับ Phone
    fullPhone := fmt.Sprintf("%s %s", req.AreaCode, req.Phone)

    // ✅ Handle `nil` values เพื่อป้องกัน Panic
    lineID := sql.NullString{}
    if req.LineID != nil {
        lineID = sql.NullString{String: *req.LineID, Valid: true}
    }

    facebook := sql.NullString{}
    if req.Facebook != nil {
        facebook = sql.NullString{String: *req.Facebook, Valid: true}
    }

    locationLink := sql.NullString{}
    if req.LocationLink != nil {
        locationLink = sql.NullString{String: *req.LocationLink, Valid: true}
    }

    // ✅ อัปเดตข้อมูลฟาร์มใน PostgreSQL
    updatedFarmer := models.Farmer{
        FarmerID:     farmer.FarmerID,
        CompanyName:  req.CompanyName,
        Address:      req.Address,
        District:     req.District,
        SubDistrict:  req.SubDistrict,
        Province:     req.Province,
        Country:      req.Country,
        PostCode:     req.PostCode,
        Telephone:    fullPhone,
        WalletAddress: farmer.WalletAddress,
        LineID:       lineID,
        Facebook:     facebook,
        LocationLink: locationLink,
    }

    if err := database.DB.Model(&farmer).Updates(updatedFarmer).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update farm information"})
    }

    // ✅ อัปเดต Blockchain เฉพาะกรณี `certCID` เปลี่ยนแปลง
    if certCID != "" && certCID != latestCertCID {
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
    }

    return c.JSON(fiber.Map{
        "message": "Farm information updated successfully!",
        "certCID": certCID,
    })
}

