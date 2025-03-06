package controllers

import (
	"encoding/json"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type RawMilkController struct {
	DB                *gorm.DB
	BlockchainService *services.BlockchainService
	QRCodeService     *services.QRCodeService
	IPFSService       *services.IPFSService
	MilkTankCounter   map[string]int
	Mutex             sync.Mutex
}

// ✅ อัปเดต Constructor ให้รองรับ `MilkTankCounter` และ `Mutex`
func NewRawMilkController(
	db *gorm.DB,
	blockchainService *services.BlockchainService,
	ipfsService *services.IPFSService,
	qrCodeService *services.QRCodeService,
) *RawMilkController {
	return &RawMilkController{
		DB:                db,
		BlockchainService: blockchainService,
		IPFSService:       ipfsService,
		QRCodeService:     qrCodeService,
		MilkTankCounter:   make(map[string]int), // ✅ กำหนดค่าให้ MilkTankCounter
		Mutex:             sync.Mutex{},         // ✅ กำหนดค่าให้ Mutex
	}
}

// ✅ ฟังก์ชันสร้าง Tank ID (FarmID + วันที่ + Running Number)
func (rmc *RawMilkController) generateTankID(farmID string) string {
	rmc.Mutex.Lock()
	defer rmc.Mutex.Unlock()

	// ✅ ดึงวันที่ปัจจุบันในรูปแบบ YYYYMMDD
	currentDate := time.Now().Format("20060102")

	// ✅ คีย์สำหรับเก็บ Running Number (FarmID + วันที่)
	key := farmID + "_" + currentDate

	// ✅ ถ้าไม่มีข้อมูลเก่า หรือเป็นวันใหม่ ให้รีเซ็ตเลขลำดับ
	if _, exists := rmc.MilkTankCounter[key]; !exists {
		rmc.MilkTankCounter[key] = 1
	} else {
		rmc.MilkTankCounter[key]++
	}

	// ✅ สร้าง Tank ID => FarmID + วันที่ + Running Number (3 หลัก)
	tankID := fmt.Sprintf("%s-%s-%03d", farmID, currentDate, rmc.MilkTankCounter[key])

	fmt.Println("✅ Generated Tank ID:", tankID)
	return tankID
}

func (rmc *RawMilkController) CreateMilkTank(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Create Milk Tank")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	farmID := c.Locals("entityID").(string)
	walletAddress := c.Locals("walletAddress").(string)

	// ✅ ตรวจสอบสิทธิ์
	if role != "farmer" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only farmers can create milk tanks"})
	}

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		MilkTankInfo    json.RawMessage `json:"milkTankInfo"` // ✅ เก็บข้อมูล MilkTankInfo แบบดิบ (Raw JSON)
		ShippingAddress struct {
			CompanyName string `json:"companyName"`
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			Email       string `json:"email"`
			AreaCode    string `json:"areaCode"`
			PhoneNumber string `json:"phoneNumber"`
			Address     string `json:"address"`
			Province    string `json:"province"`
			District    string `json:"district"`
			SubDistrict string `json:"subDistrict"`
			PostalCode  string `json:"postalCode"`
			Location    string `json:"location"`
		} `json:"shippingAddress"`
	}

	// ✅ ใช้ json.Unmarshal() เพื่อแปลงข้อมูล
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		fmt.Println("❌ Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ แปลง MilkTankInfo ที่เป็น Raw JSON ให้อยู่ใน Struct
	var milkTankInfo struct {
		FarmName        string `json:"farmName"`
		PersonInCharge  string `json:"personInCharge"`
		Quantity        string `json:"quantity"`
		QuantityUnit    string `json:"quantityUnit"`
		Temp            string `json:"temp"`
		TempUnit        string `json:"tempUnit"`
		PH              string `json:"pH"`
		Fat             string `json:"fat"`
		Protein         string `json:"protein"`
		Bacteria        bool   `json:"bacteria"`
		BacteriaInfo    string `json:"bacteriaInfo"`
		Contaminants    bool   `json:"contaminants"`
		ContaminantInfo string `json:"contaminantInfo"`
		AbnormalChar    bool   `json:"abnormalChar"`
		AbnormalType    struct {
			SmellBad      bool `json:"smellBad"`
			SmellNotFresh bool `json:"smellNotFresh"`
			AbnormalColor bool `json:"abnormalColor"`
			Sour          bool `json:"sour"`
			Bitter        bool `json:"bitter"`
			Cloudy        bool `json:"cloudy"`
			Lumpy         bool `json:"lumpy"`
			Separation    bool `json:"separation"`
		} `json:"abnormalType"`
	}

	if err := json.Unmarshal(request.MilkTankInfo, &milkTankInfo); err != nil {
		fmt.Println("❌ Error parsing MilkTankInfo:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid MilkTankInfo data"})
	}

	// ✅ Debug Log เพื่อตรวจสอบค่าที่ได้รับ
	fmt.Printf("📌 Debug - Full MilkTankInfo Struct: %+v\n", milkTankInfo)
	fmt.Println("📌 Debug - Received Person In Charge:", milkTankInfo.PersonInCharge)

	// ✅ ตรวจสอบค่า PersonInCharge ก่อนใช้งาน
	if milkTankInfo.PersonInCharge == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "personInCharge is required"})
	}

	// ✅ ค้นหา FactoryID จาก CompanyName
	var factory models.Factory
	if err := rmc.DB.Where("companyname = ?", request.ShippingAddress.CompanyName).First(&factory).Error; err != nil {
		fmt.Println("❌ Factory not found:", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Factory not found"})
	}

	// ✅ สร้าง `tankId`
	tankId := rmc.generateTankID(farmID)

	// ✅ แปลงค่าที่เป็น string → uint64
	quantity, _ := strconv.ParseUint(milkTankInfo.Quantity, 10, 64)
	temp, _ := strconv.ParseUint(milkTankInfo.Temp, 10, 64)
	ph, _ := strconv.ParseUint(milkTankInfo.PH, 10, 64)
	fat, _ := strconv.ParseUint(milkTankInfo.Fat, 10, 64)
	protein, _ := strconv.ParseUint(milkTankInfo.Protein, 10, 64)

	// ✅ รวมข้อมูลอัปโหลดไป IPFS
	milkMetadata := map[string]interface{}{
		"farmName":        milkTankInfo.FarmName,
		"personInCharge":  milkTankInfo.PersonInCharge,
		"quantity":        quantity,
		"quantityUnit":    milkTankInfo.QuantityUnit,
		"temperature":     temp,
		"tempUnit":        milkTankInfo.TempUnit,
		"pH":              ph,
		"fat":             fat,
		"protein":         protein,
		"bacteria":        milkTankInfo.Bacteria,
		"bacteriaInfo":    milkTankInfo.BacteriaInfo,
		"contaminants":    milkTankInfo.Contaminants,
		"contaminantInfo": milkTankInfo.ContaminantInfo,
		"abnormalChar":    milkTankInfo.AbnormalChar,
		"abnormalType":    milkTankInfo.AbnormalType,
		"shippingAddress": request.ShippingAddress,
	}

	// ✅ อัปโหลด IPFS
	qualityReportCID, err := rmc.IPFSService.UploadMilkDataToIPFS(milkMetadata, nil)
	if err != nil {
		fmt.Println("❌ Failed to upload to IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload quality report"})
	}

	// ✅ สร้าง QR Code
	qrCodeCID, err := rmc.QRCodeService.GenerateQRCode(tankId)
	if err != nil {
		fmt.Println("❌ Failed to generate QR Code:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate QR Code"})
	}

	// ✅ ส่งไป Blockchain
	txHash, err := rmc.BlockchainService.CreateMilkTank(
		walletAddress,
		tankId,
		factory.FactoryID,
		milkTankInfo.PersonInCharge,
		qualityReportCID,
		qrCodeCID,
	)
	if err != nil {
		fmt.Println("❌ Blockchain transaction failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
	}

	// ✅ ส่ง Response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":          "Milk tank created successfully",
		"tankId":           tankId,
		"txHash":           txHash,
		"qrCodeCID":        qrCodeCID,
		"qualityReportCID": qualityReportCID,
	})
}

func (rmc *RawMilkController) GetFarmRawMilkTanks(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Farm Raw Milk Tanks")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	farmerWallet := c.Locals("walletAddress").(string)

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Farmer เท่านั้น)
	if role != "farmer" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only farmers can view raw milk tanks"})
	}

	// ✅ ดึงค่าที่พิมพ์ในช่องค้นหา (Search Query)
	searchQuery := strings.ToLower(c.Query("search", ""))

	// ✅ ดึงข้อมูลแท็งก์จาก Blockchain
	milkTanks, err := rmc.BlockchainService.GetMilkTanksByFarmer(farmerWallet)
	if err != nil {
		fmt.Println("❌ Failed to fetch raw milk tanks:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch raw milk tanks"})
	}

	// ✅ กรองข้อมูลตาม Search Query
	var filteredMilkTanks []map[string]interface{}
	for _, tank := range milkTanks {
		tankId := tank["tankId"].(string)
		personInCharge := tank["personInCharge"].(string)

		// ✅ ถ้า searchQuery ว่าง → แสดงทั้งหมด, ถ้าไม่ว่าง → ค้นหาตาม Tank ID หรือ Person in Charge
		if searchQuery == "" || strings.Contains(strings.ToLower(tankId), searchQuery) || strings.Contains(strings.ToLower(personInCharge), searchQuery) {
			filteredMilkTanks = append(filteredMilkTanks, map[string]interface{}{
				"tankId":         tankId,
				"personInCharge": personInCharge,
				"status":         tank["status"].(uint8), // แปลงค่า Enum เป็นเลข
				"moreInfoLink":   fmt.Sprintf("/Farmer/FarmDetails?id=%s", tankId),
			})
		}
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"displayedMilkTanks": filteredMilkTanks,
		"addNewTankLink":     "/Farmer/FarmCreateRM",
	})
}

// ✅ ฟังก์ชันดึงข้อมูลแท็งก์นมดิบตาม Tank ID
func (rmc *RawMilkController) GetRawMilkTankDetails(c *fiber.Ctx) error {
	tankId := c.Params("tankId") // ✅ รับ tankId จาก URL Parameter
	fmt.Println("📌 Request received: Fetching milk tank details for:", tankId)

	// ✅ ดึงข้อมูลแท็งก์จาก Blockchain
	rawMilk, err := rmc.BlockchainService.GetRawMilkTankDetails(tankId)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tank details:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch milk tank details"})
	}

	// ✅ ดึงข้อมูลรายละเอียดเพิ่มเติมจาก IPFS โดยใช้ `QualityReportCID`
	ipfsCID := rawMilk.QualityReportCID
	ipfsData, err := rmc.IPFSService.GetFromIPFS(ipfsCID)
	if err != nil {
		fmt.Println("❌ Failed to fetch data from IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch quality report from IPFS"})
	}

	// ✅ แปลงข้อมูลจาก JSON (IPFS)
	var ipfsRawMilkData map[string]interface{}
	err = json.Unmarshal(ipfsData, &ipfsRawMilkData)
	if err != nil {
		fmt.Println("❌ Failed to parse IPFS JSON:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid JSON format from IPFS"})
	}

	// ✅ ตรวจสอบว่ามี `rawMilkData` อยู่ใน JSON หรือไม่
	rawMilkData, ok := ipfsRawMilkData["rawMilkData"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Missing rawMilkData in IPFS response")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Missing raw milk data in IPFS"})
	}

	// ✅ ตรวจสอบว่ามี `shippingAddress` หรือไม่
	var shippingAddress map[string]interface{}
	if rawMilkData["shippingAddress"] != nil {
		shippingAddress, ok = rawMilkData["shippingAddress"].(map[string]interface{})
		if !ok {
			fmt.Println("❌ Invalid shippingAddress format in IPFS")
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid shipping address format in IPFS"})
		}
	} else {
		shippingAddress = map[string]interface{}{} // ✅ ใช้ค่าเริ่มต้นว่าง
	}

	// ✅ ป้องกัน Panic Error โดยตรวจสอบค่าก่อนแปลง Type
	getString := func(key string, data map[string]interface{}) string {
		if value, ok := data[key].(string); ok {
			return value
		}
		return ""
	}

	getFloat64 := func(key string, data map[string]interface{}) float64 {
		if value, ok := data[key].(float64); ok {
			return value
		}
		return 0.0
	}

	getBool := func(key string, data map[string]interface{}) bool {
		if value, ok := data[key].(bool); ok {
			return value
		}
		return false
	}

	// ✅ สร้างโครงสร้างข้อมูลตามที่ Frontend ต้องการ
	responseData := fiber.Map{
		"milkTankInfo": fiber.Map{
			"farmName":        getString("farmName", rawMilkData),
			"milkTankNo":      rawMilk.TankId, // ✅ ใช้ Tank ID ที่ดึงจาก Blockchain
			"personInCharge":  rawMilk.PersonInCharge,
			"quantity":        getFloat64("quantity", rawMilkData),
			"quantityUnit":    getString("quantityUnit", rawMilkData),
			"temp":            getFloat64("temperature", rawMilkData),
			"tempUnit":        getString("tempUnit", rawMilkData),
			"pH":              getFloat64("pH", rawMilkData),
			"fat":             getFloat64("fat", rawMilkData),
			"protein":         getFloat64("protein", rawMilkData),
			"bacteria":        getBool("bacteria", rawMilkData),
			"bacteriaInfo":    getString("bacteriaInfo", rawMilkData),
			"contaminants":    getBool("contaminants", rawMilkData),
			"contaminantInfo": getString("contaminantInfo", rawMilkData),
			"abnormalChar":    getBool("abnormalChar", rawMilkData),
			"abnormalType":    rawMilkData["abnormalType"], // ✅ ส่งทั้ง Object กลับ
		},
		"shippingAddress": fiber.Map{
			"companyName": getString("companyName", shippingAddress),
			"firstName":   getString("firstName", shippingAddress),
			"lastName":    getString("lastName", shippingAddress),
			"email":       getString("email", shippingAddress),
			"areaCode":    getString("areaCode", shippingAddress),
			"phoneNumber": getString("phoneNumber", shippingAddress),
			"address":     getString("address", shippingAddress),
			"province":    getString("province", shippingAddress),
			"district":    getString("district", shippingAddress),
			"subDistrict": getString("subDistrict", shippingAddress),
			"postalCode":  getString("postalCode", shippingAddress),
			"location":    getString("location", shippingAddress),
		},
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(responseData)
}

func (rmc *RawMilkController) GetQRCodeByTankID(c *fiber.Ctx) error {
	tankId := c.Params("tankId")
	fmt.Println("📌 Fetching QR Code for Tank ID:", tankId)

	// ✅ ดึงรายละเอียดแท็งก์จาก Blockchain
	rawMilkData, err := rmc.BlockchainService.GetRawMilkTankDetails(tankId)
	if err != nil {
		fmt.Println("❌ Failed to fetch tank details:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tank details"})
	}

	// ✅ ตรวจสอบว่ามี QR Code CID หรือไม่
	if rawMilkData.QrCodeCID == "" {
		fmt.Println("❌ QR Code not found for Tank ID:", tankId)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "QR Code not found"})
	}

	// ✅ ดึง QR Code จาก IPFS
	qrCodeBase64, err := rmc.IPFSService.GetImageBase64FromIPFS(rawMilkData.QrCodeCID)
	if err != nil {
		fmt.Println("❌ Failed to retrieve QR Code from IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve QR Code"})
	}

	// ✅ ส่ง Base64 QR Code กลับไปที่ Frontend
	return c.JSON(fiber.Map{
		"tankId":    tankId,
		"qrCodeCID": rawMilkData.QrCodeCID,
		"qrCodeImg": fmt.Sprintf("data:image/png;base64,%s", qrCodeBase64),
	})
}
