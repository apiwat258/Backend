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

// ✅ สร้าง Tank ID โดยดึง FarmID จาก Token โดยตรง
func (rmc *RawMilkController) GenerateTankID(c *fiber.Ctx) error {
	farmID := c.Locals("entityID").(string)
	currentDate := time.Now().Format("20060102")
	key := farmID + "_" + currentDate

	rmc.Mutex.Lock()
	rmc.MilkTankCounter[key]++
	count := rmc.MilkTankCounter[key]
	rmc.Mutex.Unlock()

	tankID := fmt.Sprintf("%s-%s-%03d", farmID, currentDate, count)

	return c.Status(200).JSON(fiber.Map{
		"tankId": tankID,
	})
}

func (rmc *RawMilkController) CreateMilkTank(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Create Milk Tank")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	walletAddress := c.Locals("walletAddress").(string)

	// ✅ ตรวจสอบสิทธิ์
	if role != "farmer" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only farmers can create milk tanks"})
	}

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		MilkTankInfo    json.RawMessage `json:"milkTankInfo"`
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
		TankID          string `json:"TankId"`
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

	tankId := milkTankInfo.TankID

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
		oldPersonInCharge, hasOldPerson := tank["oldPersonInCharge"].(string) // ✅ ตรวจสอบว่ามี oldPersonInCharge ไหม

		// ✅ ถ้ามี Old Person In Charge ให้ใช้แทน
		if hasOldPerson && oldPersonInCharge != "" {
			personInCharge = oldPersonInCharge
		}

		// ✅ ถ้า searchQuery ว่าง → แสดงทั้งหมด, ถ้าไม่ว่าง → ค้นหาตาม Tank ID หรือ Person in Charge
		if searchQuery == "" || strings.Contains(strings.ToLower(tankId), searchQuery) || strings.Contains(strings.ToLower(personInCharge), searchQuery) {
			filteredMilkTanks = append(filteredMilkTanks, map[string]interface{}{
				"tankId":         strings.TrimRight(tankId, "\x00"),
				"personInCharge": personInCharge,         // ✅ ใช้ Old Person ถ้ามี
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

// ///For all/////////
func (rmc *RawMilkController) GetRawMilkTankDetails(c *fiber.Ctx) error {
	tankId := c.Params("tankId") // ✅ รับ tankId จาก URL Parameter
	fmt.Println("📌 Request received: Fetching milk tank details for:", tankId)

	// ✅ ดึงข้อมูลแท็งก์และประวัติจาก Blockchain
	rawMilk, history, err := rmc.BlockchainService.GetRawMilkTankDetails(tankId)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tank details:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch milk tank details"})
	}

	// ✅ สร้างโครงสร้างสำหรับ response
	responseData := fiber.Map{}

	// ✅ ตรวจสอบสถานะของแท็งก์
	var farmCID, factoryCID string

	if len(history) > 0 {
		// ✅ ดึง CID ของ `Status = 0` จากประวัติ (ฟาร์ม)
		for _, entry := range history {
			if status, ok := entry["status"].(uint8); ok && status == 0 {
				farmCID, _ = entry["qualityReportCID"].(string)
				break
			}
		}
	}

	fmt.Println("📌 Final farmRepo CID:", farmCID)
	fmt.Println("📌 Final factoryRepo CID:", rawMilk.QualityReportCID)

	// ✅ ดึงข้อมูลฟาร์มจาก IPFS
	if farmCID != "" {
		fmt.Println("📌 Retrieving farmRepo from IPFS... CID:", farmCID)
		ipfsData, err := rmc.IPFSService.GetFromIPFS(farmCID)
		if err != nil {
			fmt.Println("❌ Failed to fetch farmRepo from IPFS:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch farm quality report from IPFS"})
		}

		var ipfsFarmData map[string]interface{}
		err = json.Unmarshal(ipfsData, &ipfsFarmData)
		if err != nil {
			fmt.Println("❌ Failed to parse farmRepo JSON:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid JSON format from IPFS for farm"})
		}

		responseData["farmRepo"] = ipfsFarmData
	}

	// ✅ ดึงข้อมูลโรงงานจาก IPFS (ใช้ factoryCID ที่ถูกต้อง และดึงเพียงครั้งเดียว)
	if rawMilk.Status != 0 {
		factoryCID = rawMilk.QualityReportCID
		if factoryCID != "" {
			fmt.Println("📌 Retrieving factoryRepo from IPFS... CID:", factoryCID)
			ipfsData, err := rmc.IPFSService.GetFromIPFS(factoryCID)
			if err != nil {
				fmt.Println("❌ Failed to fetch factoryRepo from IPFS:", err)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch factory quality report from IPFS"})
			}

			var ipfsFactoryData map[string]interface{}
			err = json.Unmarshal(ipfsData, &ipfsFactoryData)
			if err != nil {
				fmt.Println("❌ Failed to parse factoryRepo JSON:", err)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid JSON format from IPFS for factory"})
			}

			responseData["factoryRepo"] = ipfsFactoryData
		}
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(responseData)
}

// ✅ ฟังก์ชันช่วยแยกข้อมูล `farmRepo`
func extractFarmRepo(history []map[string]interface{}) map[string]interface{} {
	if len(history) == 0 {
		return nil
	}
	latestEntry := history[0] // ✅ ดึงข้อมูลแรกสุด (ฟาร์มที่สร้าง)
	return map[string]interface{}{
		"farmName":        latestEntry["farmName"],
		"personInCharge":  latestEntry["personInCharge"],
		"quantity":        latestEntry["quantity"],
		"quantityUnit":    latestEntry["quantityUnit"],
		"temp":            latestEntry["temp"],
		"tempUnit":        latestEntry["tempUnit"],
		"pH":              latestEntry["pH"],
		"fat":             latestEntry["fat"],
		"protein":         latestEntry["protein"],
		"bacteria":        latestEntry["bacteria"],
		"bacteriaInfo":    latestEntry["bacteriaInfo"],
		"contaminants":    latestEntry["contaminants"],
		"contaminantInfo": latestEntry["contaminantInfo"],
		"abnormalChar":    latestEntry["abnormalChar"],
		"abnormalType":    latestEntry["abnormalType"],
	}
}

// ✅ ฟังก์ชันช่วยแยกข้อมูล `factoryRepo`
func extractFactoryRepo(ipfsRawMilkData map[string]interface{}) map[string]interface{} {
	rawMilkData, ok := ipfsRawMilkData["rawMilkData"].(map[string]interface{})
	if !ok {
		return nil
	}
	return map[string]interface{}{
		"personInCharge":  rawMilkData["recipientInfo"].(map[string]interface{})["personInCharge"],
		"location":        rawMilkData["recipientInfo"].(map[string]interface{})["location"],
		"pickUpTime":      rawMilkData["recipientInfo"].(map[string]interface{})["pickUpTime"],
		"quantity":        rawMilkData["quantity"],
		"quantityUnit":    rawMilkData["quantityUnit"],
		"temp":            rawMilkData["temperature"],
		"tempUnit":        rawMilkData["tempUnit"],
		"pH":              rawMilkData["pH"],
		"fat":             rawMilkData["fat"],
		"protein":         rawMilkData["protein"],
		"bacteria":        rawMilkData["bacteria"],
		"bacteriaInfo":    rawMilkData["bacteriaInfo"],
		"contaminants":    rawMilkData["contaminants"],
		"contaminantInfo": rawMilkData["contaminantInfo"],
		"abnormalChar":    rawMilkData["abnormalChar"],
		"abnormalType":    rawMilkData["abnormalType"],
	}
}

func (rmc *RawMilkController) GetQRCodeByTankID(c *fiber.Ctx) error {
	tankId := c.Params("tankId")
	fmt.Println("📌 Fetching QR Code for Tank ID:", tankId)

	// ✅ ดึงรายละเอียดแท็งก์จาก Blockchain (คืนค่า rawMilkData, history)
	rawMilkData, _, err := rmc.BlockchainService.GetRawMilkTankDetails(tankId)
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

// //////For Factory////
/*func (rmc *RawMilkController) GetFactoryRawMilkTanks(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Factory Raw Milk Tanks")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Factory เท่านั้น)
	if role != "factory" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only factories can view raw milk tanks"})
	}

	// ✅ ดึง `entityID` จาก JWT Token ที่ AuthMiddleware กำหนดไว้
	factoryID, ok := c.Locals("entityID").(string)
	if !ok || factoryID == "" {
		fmt.Println("❌ Factory ID is missing in Context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - Factory ID is missing"})
	}
	fmt.Println("✅ Factory ID from Context:", factoryID)

	// ✅ ดึงค่าที่พิมพ์ในช่องค้นหา (Search Query)
	searchQuery := strings.ToLower(c.Query("search", ""))

	// ✅ ดึงข้อมูลแท็งก์จาก Blockchain
	milkTanks, err := rmc.BlockchainService.GetMilkTanksByFactory(factoryID)
	if err != nil {
		fmt.Println("❌ Failed to fetch raw milk tanks for factory:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch raw milk tanks"})
	}

	// ✅ กรองข้อมูลตาม Search Query
	var filteredMilkTanks []map[string]interface{}
	for _, tank := range milkTanks {
		tankId := tank["tankId"].(string)
		personInCharge := tank["personInCharge"].(string)

		if searchQuery == "" || strings.Contains(strings.ToLower(tankId), searchQuery) || strings.Contains(strings.ToLower(personInCharge), searchQuery) {
			filteredMilkTanks = append(filteredMilkTanks, map[string]interface{}{
				"tankId":         strings.TrimRight(tankId, "\x00"),
				"personInCharge": personInCharge,
				"status":         tank["status"].(uint8),
				"moreInfoLink":   fmt.Sprintf("/Factory/FactoryDetails?id=%s", tankId),
			})
		}
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"displayedMilkTanks": filteredMilkTanks,
	})
}*/
func (rmc *RawMilkController) GetFactoryRawMilkTanks(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Factory Raw Milk Tanks")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Factory เท่านั้น)
	if role != "factory" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only factories can view raw milk tanks"})
	}

	// ✅ ดึง `entityID` จาก JWT Token
	factoryID, ok := c.Locals("entityID").(string)
	if !ok || factoryID == "" {
		fmt.Println("❌ Factory ID is missing in Context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - Factory ID is missing"})
	}
	fmt.Println("✅ Factory ID from Context:", factoryID)

	// ✅ ดึงค่าที่พิมพ์ในช่องค้นหา (Search Query)
	searchQuery := strings.ToLower(c.Query("search", ""))

	// ✅ ดึงข้อมูลแท็งก์จาก Blockchain
	milkTanks, err := rmc.BlockchainService.GetMilkTanksByFactory(factoryID)
	if err != nil {
		fmt.Println("❌ Failed to fetch raw milk tanks for factory:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch raw milk tanks"})
	}

	// ✅ เตรียม response
	var displayedMilkTanks []map[string]interface{}

	for _, tank := range milkTanks {
		tankId := strings.TrimRight(tank["tankId"].(string), "\x00")
		personInCharge := tank["personInCharge"].(string)
		status := tank["status"].(uint8)

		// ✅ ดึง `farmID` จาก `tankId`
		parts := strings.Split(tankId, "-")
		if len(parts) < 1 {
			fmt.Println("❌ Invalid Tank ID format:", tankId)
			continue
		}
		farmID := parts[0]

		// ✅ ดึง `farmName` และ `location` จาก PostgreSQL
		var farmer models.Farmer
		if err := rmc.DB.Where("farmerid = ?", farmID).First(&farmer).Error; err != nil {
			fmt.Println("❌ Failed to fetch farm details:", err)
			continue
		}

		// ✅ ใช้ `GetRawMilkTankDetails` ดึง factoryCID
		rawMilkDetails, _, err := rmc.BlockchainService.GetRawMilkTankDetails(tankId)
		if err != nil {
			fmt.Println("❌ Failed to fetch tank details for:", tankId)
			continue
		}

		factoryCID := rawMilkDetails.QualityReportCID
		fmt.Println("📌 Found factoryRepo CID:", factoryCID)

		// ✅ ดึงข้อมูลจาก IPFS (ถ้ามี CID)
		var quantity, quantityUnit, temperature, tempUnit string
		if factoryCID != "" {
			ipfsData, err := rmc.IPFSService.GetFromIPFS(factoryCID)
			if err == nil {
				var ipfsFactoryData map[string]interface{}
				if err := json.Unmarshal(ipfsData, &ipfsFactoryData); err == nil {
					// ✅ ดึงค่าที่ต้องการจาก `factoryRepo.rawMilkData`
					if rawMilkData, exists := ipfsFactoryData["rawMilkData"].(map[string]interface{}); exists {
						if q, ok := rawMilkData["quantity"].(float64); ok {
							quantity = fmt.Sprintf("%.2f", q)
						}
						if qUnit, ok := rawMilkData["quantityUnit"].(string); ok {
							quantityUnit = qUnit
						}
						if temp, ok := rawMilkData["temperature"].(float64); ok {
							temperature = fmt.Sprintf("%.2f", temp)
						}
						if tUnit, ok := rawMilkData["tempUnit"].(string); ok {
							tempUnit = tUnit
						}
					}
				}
			}
		}

		// ✅ รวม `quantity` + `quantityUnit` และ `temperature` + `tempUnit`
		quantityInfo := quantity + " " + quantityUnit
		temperatureInfo := temperature + " " + tempUnit

		// ✅ กรองข้อมูลตาม Search Query
		if searchQuery == "" || strings.Contains(strings.ToLower(tankId), searchQuery) || strings.Contains(strings.ToLower(personInCharge), searchQuery) {
			displayedMilkTanks = append(displayedMilkTanks, map[string]interface{}{
				"tankId":         tankId,
				"personInCharge": personInCharge,
				"status":         status,
				"moreInfoLink":   fmt.Sprintf("/Factory/FactoryDetails?id=%s", tankId),
				"farmName":       farmer.CompanyName, // ✅ เพิ่มชื่อฟาร์ม
				"location":       farmer.Province,    // ✅ เพิ่มโลเคชันของฟาร์ม
				"quantity":       quantityInfo,       // ✅ เพิ่มจำนวน
				"temperature":    temperatureInfo,    // ✅ เพิ่มอุณหภูมิ
			})
		}
	}

	// ✅ ส่ง Response กลับไปที่ Frontend (เหมือนโครงสร้างเดิม)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"displayedMilkTanks": displayedMilkTanks,
	})
}

func (rmc *RawMilkController) UpdateMilkTankStatus(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Update Milk Tank Status")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	walletAddress := c.Locals("walletAddress").(string)

	// ✅ ตรวจสอบสิทธิ์
	if role != "factory" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only factories can update milk tanks"})
	}

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		TankID   string `json:"tankId"`
		Approved bool   `json:"approved"`
		Input    struct {
			RecipientInfo struct {
				PersonInCharge string `json:"personInCharge"`
				Location       string `json:"location"`
				PickUpTime     string `json:"pickUpTime"`
			} `json:"RecipientInfo"`
			Quantity struct {
				Quantity        float64 `json:"quantity"`
				QuantityUnit    string  `json:"quantityUnit"`
				Temp            float64 `json:"temp"`
				TempUnit        string  `json:"tempUnit"`
				PH              float64 `json:"pH"`
				Fat             float64 `json:"fat"`
				Protein         float64 `json:"protein"`
				Bacteria        bool    `json:"bacteria"`
				BacteriaInfo    string  `json:"bacteriaInfo"`
				Contaminants    bool    `json:"contaminants"`
				ContaminantInfo string  `json:"contaminantInfo"`
				AbnormalChar    bool    `json:"abnormalChar"`
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
			} `json:"Quantity"`
		} `json:"input"`
	}

	// ✅ ตรวจสอบ JSON Request
	if err := c.BodyParser(&request); err != nil {
		fmt.Println("❌ Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ ตรวจสอบค่าที่จำเป็น
	if request.TankID == "" || request.Input.RecipientInfo.PersonInCharge == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// ✅ รวมข้อมูลอัปโหลดไป IPFS
	milkMetadata := map[string]interface{}{
		"recipientInfo": map[string]interface{}{
			"personInCharge": request.Input.RecipientInfo.PersonInCharge,
			"location":       request.Input.RecipientInfo.Location,
			"pickUpTime":     request.Input.RecipientInfo.PickUpTime,
		},
		"quantity":        request.Input.Quantity.Quantity,
		"quantityUnit":    request.Input.Quantity.QuantityUnit,
		"temperature":     request.Input.Quantity.Temp,
		"tempUnit":        request.Input.Quantity.TempUnit,
		"pH":              request.Input.Quantity.PH,
		"fat":             request.Input.Quantity.Fat,
		"protein":         request.Input.Quantity.Protein,
		"bacteria":        request.Input.Quantity.Bacteria,
		"bacteriaInfo":    request.Input.Quantity.BacteriaInfo,
		"contaminants":    request.Input.Quantity.Contaminants,
		"contaminantInfo": request.Input.Quantity.ContaminantInfo,
		"abnormalChar":    request.Input.Quantity.AbnormalChar,
		"abnormalType":    request.Input.Quantity.AbnormalType,
	}

	// ✅ อัปโหลดข้อมูลไปยัง IPFS
	qualityReportCID, err := rmc.IPFSService.UploadMilkDataToIPFS(milkMetadata, nil)
	if err != nil {
		fmt.Println("❌ Failed to upload to IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload quality report"})
	}

	// ✅ อัปเดตสถานะไปยัง Blockchain
	txHash, err := rmc.BlockchainService.UpdateMilkTankStatus(
		walletAddress,
		request.TankID,
		request.Approved,
		request.Input.RecipientInfo.PersonInCharge,
		qualityReportCID,
	)
	if err != nil {
		fmt.Println("❌ Blockchain transaction failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":          "Milk tank status updated successfully",
		"tankId":           request.TankID,
		"txHash":           txHash,
		"qualityReportCID": qualityReportCID,
	})
}
