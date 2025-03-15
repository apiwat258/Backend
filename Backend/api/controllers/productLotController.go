package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"finalyearproject/Backend/database"
	"finalyearproject/Backend/services"
)

// ✅ ProductLotController โครงสร้าง
type ProductLotController struct {
	DB                *gorm.DB
	IPFSService       *services.IPFSService
	BlockchainService *services.BlockchainService
	QRService         *services.QRCodeService // ✅ เพิ่ม QR Code Service
}

// ✅ แก้ไข Constructor ให้รับ QRService ด้วย
func NewProductLotController(db *gorm.DB, blockchainService *services.BlockchainService, ipfsService *services.IPFSService, qrService *services.QRCodeService) *ProductLotController {
	return &ProductLotController{
		DB:                db,
		BlockchainService: blockchainService,
		IPFSService:       ipfsService,
		QRService:         qrService, // ✅ เพิ่ม QRService เข้าไป
	}
}

// ✅ ฟังก์ชันสร้าง Product Lot พร้อม Tracking Event
func (plc *ProductLotController) CreateProductLot(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Create Product Lot")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	factoryID := c.Locals("entityID").(string)
	walletAddress := c.Locals("walletAddress").(string)
	userID := c.Locals("userID").(string)

	// ✅ ตรวจสอบสิทธิ์
	if role != "factory" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only factories can create product lots"})
	}

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		ProductID         string          `json:"productId"`
		Grade             bool            `json:"grade"`
		MilkTankIDs       []string        `json:"milkTankIds"`
		QualityData       json.RawMessage `json:"qualityData"`
		ShippingAddresses []struct {
			RetailerID  string `json:"retailerId"`
			CompanyName string `json:"companyName"`
			Address     string `json:"address"`
			Province    string `json:"province"`
			District    string `json:"district"`
			SubDistrict string `json:"subDistrict"`
			PostalCode  string `json:"postalCode"`
			Location    string `json:"location"`
		} `json:"shippingAddresses"`
	}
	// ✅ Debug: ตรวจสอบ JSON ที่ได้รับทั้งหมด
	body := c.Body()
	fmt.Println("📌 Debug - Raw Request Body:", string(body))

	// ✅ แปลง JSON
	if err := json.Unmarshal(body, &request); err != nil {
		fmt.Println("❌ Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ Debug: แสดงค่า Shipping Addresses ที่รับเข้ามา
	fmt.Println("📌 Debug - Parsed Shipping Addresses:", request.ShippingAddresses)

	if len(request.ShippingAddresses) == 0 {
		fmt.Println("❌ No Shipping Addresses Found!")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No Shipping Addresses provided"})
	}

	if err := json.Unmarshal(c.Body(), &request); err != nil {
		fmt.Println("❌ Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if strings.TrimSpace(request.ProductID) == "" || len(request.MilkTankIDs) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Product ID and at least one Milk Tank are required"})
	}

	// ✅ ค้นหา Inspector Name
	var inspectorName string
	err := database.DB.Table("users").Where("userid = ?", userID).Select("username").Scan(&inspectorName).Error
	if err != nil || inspectorName == "" {
		fmt.Println("❌ Failed to find inspector name:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve inspector name"})
	}

	// ✅ สร้าง `lotId`
	lotId := plc.generateLotID(factoryID)

	// ✅ อัปโหลด `Quality & Nutrition` ไปที่ IPFS
	qualityCID, err := plc.IPFSService.UploadDataToIPFS(map[string]interface{}{
		"qualityData": json.RawMessage(request.QualityData),
	})
	if err != nil {
		fmt.Println("❌ Failed to upload quality data to IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload quality data"})
	}

	// ✅ ส่งธุรกรรมไปที่ Blockchain (สร้าง Product Lot)
	txHash, err := plc.BlockchainService.CreateProductLot(
		walletAddress,
		lotId,
		request.ProductID,
		inspectorName,
		request.Grade,
		qualityCID,
		request.MilkTankIDs,
	)
	if err != nil {
		fmt.Println("❌ Blockchain transaction failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
	}
	fmt.Println("📌 Debug - Checking Shipping Addresses")
	fmt.Println("   - Total Addresses:", len(request.ShippingAddresses))

	if len(request.ShippingAddresses) == 0 {
		fmt.Println("❌ No Shipping Addresses Found! Skipping Tracking Event creation.")
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"message": "Product Lot created, but no shipping addresses provided.",
			"lotId":   lotId,
			"txHash":  txHash,
			"ipfsCID": qualityCID,
		})
	}

	// ✅ สร้าง Tracking Event สำหรับทุก Retailer
	var trackingTxHashes []string
	for _, shipping := range request.ShippingAddresses {
		fmt.Println("📌 Debug - Processing Shipping Address:", shipping.RetailerID)

		if shipping.RetailerID == "" {
			fmt.Println("❌ Skipping empty Retailer ID")
			continue
		}

		trackingID := plc.GenerateTrackingID(lotId, shipping.RetailerID)
		fmt.Println("📌 Debug - Generated Tracking ID:", trackingID)

		// ✅ ค้นหา Factory Name จากตาราง `dairyfactory`
		// ✅ ค้นหา Factory Name จากตาราง `dairyfactory`
		var factoryName string
		err := database.DB.Table("dairyfactory").Where("factoryid = ?", factoryID).Select("companyname").Scan(&factoryName).Error
		if err != nil || factoryName == "" {
			fmt.Println("❌ Failed to find factory name:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve factory name"})
		}

		// ✅ เตรียมข้อมูล QR Code (เพิ่มข้อมูลให้ครบ)
		qrData := map[string]interface{}{
			"trackingId":   trackingID,
			"productLotId": lotId,
			"retailer": map[string]string{
				"retailerId":  shipping.RetailerID,
				"companyName": shipping.CompanyName,
				"address":     shipping.Address,
				"province":    shipping.Province,
				"district":    shipping.District,
				"subDistrict": shipping.SubDistrict,
				"postalCode":  shipping.PostalCode,
				"location":    shipping.Location,
			},
			"factory": map[string]string{
				"factoryId":   factoryID,
				"factoryName": factoryName,
			},
		}

		// ✅ แปลงเป็น JSON
		qrDataJSON, err := json.Marshal(qrData)
		if err != nil {
			fmt.Println("❌ Failed to encode QR data:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode QR data"})
		}

		// ✅ สร้าง QR Code ใหม่ที่เก็บข้อมูลทั้งหมด
		qrImageCID, err := plc.QRService.GenerateQRCodeforFactory(string(qrDataJSON))
		if err != nil {
			fmt.Println("❌ Failed to generate and upload QR Code:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate and upload QR Code"})
		}

		fmt.Println("✅ Debug - QR Code CID:", qrImageCID)
		if err != nil {
			fmt.Println("❌ Failed to generate and upload QR Code:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate and upload QR Code"})
		}
		fmt.Println("✅ Debug - QR Code CID:", qrImageCID)

		txHashTracking, err := plc.BlockchainService.CreateTrackingEvent(
			walletAddress,
			trackingID,
			lotId,
			shipping.RetailerID,
			qrImageCID,
		)
		if err != nil {
			fmt.Println("❌ Blockchain tracking event failed for Retailer:", shipping.RetailerID)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain tracking event failed"})
		}

		fmt.Println("✅ Tracking Event Created on Blockchain:", txHashTracking)
		trackingTxHashes = append(trackingTxHashes, txHashTracking)
	}

	// ✅ ส่ง Response กลับไปให้ Frontend
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":          "Product Lot and Tracking Events created successfully",
		"lotId":            lotId,
		"txHash":           txHash,
		"ipfsCID":          qualityCID,
		"inspector":        inspectorName,
		"trackingTxHashes": trackingTxHashes, // ✅ คืนค่าทั้งหมดที่สร้างได้
	})
}

/*
func (plc *ProductLotController) CreateProductLot(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Create Product Lot")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	factoryID := c.Locals("entityID").(string)
	walletAddress := c.Locals("walletAddress").(string)
	userID := c.Locals("userID").(string) // ✅ ใช้ userID เพื่อนำไปดึง Inspector Name

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Factory เท่านั้น)
	if role != "factory" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only factories can create product lots"})
	}

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		ProductID   string          `json:"productId"`
		Grade       bool            `json:"grade"`
		MilkTankIDs []string        `json:"milkTankIds"`
		QualityData json.RawMessage `json:"qualityData"` // ✅ เก็บข้อมูลโภชนาการแบบ JSON
	}

	// ✅ แปลง JSON
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		fmt.Println("❌ Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ ตรวจสอบว่า Product ID ต้องไม่ว่าง
	if strings.TrimSpace(request.ProductID) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Product ID is required"})
	}

	// ✅ ตรวจสอบว่าเลือก Milk Tank อย่างน้อย 1 ตัว
	if len(request.MilkTankIDs) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "At least one Milk Tank must be selected"})
	}

	// ✅ ค้นหา Inspector Name จาก Database โดยใช้ User ID
	var inspectorName string
	err := database.DB.Table("users").Where("userid = ?", userID).Select("username").Scan(&inspectorName).Error
	if err != nil || inspectorName == "" {
		fmt.Println("❌ Failed to find inspector name:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve inspector name"})
	}

	// ✅ สร้าง `lotId`
	lotId := plc.generateLotID(factoryID)

	// ✅ อัปโหลด `Quality & Nutrition` ไปที่ IPFS
	qualityCID, err := plc.IPFSService.UploadDataToIPFS(map[string]interface{}{
		"qualityData": json.RawMessage(request.QualityData),
	})
	if err != nil {
		fmt.Println("❌ Failed to upload quality data to IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload quality data"})
	}

	// ✅ ส่งธุรกรรมไปที่ Blockchain
	txHash, err := plc.BlockchainService.CreateProductLot(
		walletAddress,
		lotId,
		request.ProductID,
		inspectorName,
		request.Grade,
		qualityCID,
		request.MilkTankIDs,
	)
	if err != nil {
		fmt.Println("❌ Blockchain transaction failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
	}

	// ✅ ส่ง Response กลับไปให้ Frontend
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":   "Product Lot created successfully",
		"lotId":     lotId,
		"txHash":    txHash,
		"ipfsCID":   qualityCID,
		"inspector": inspectorName,
	})
}*/

// ✅ ฟังก์ชันสร้าง Lot ID (ใช้ Factory ID)
func (plc *ProductLotController) generateLotID(factoryID string) string {
	return fmt.Sprintf("LOT-%s-%d", factoryID, time.Now().Unix())
}

func (plc *ProductLotController) GenerateTrackingID(lotID string, retailerID string) string {
	// ✅ ดึงตัวเลข 6 หลักสุดท้ายของ Lot ID (ถ้า Lot ID ยาวกว่า 6 ตัว)
	lotSuffix := lotID[len(lotID)-6:]

	// ✅ ดึงเลข 3 หลักสุดท้ายของ Retailer ID เช่น RE000025 → 025
	retailerSuffix := retailerID[len(retailerID)-3:]

	// ✅ เพิ่มตัวเลขสุ่ม 3 หลัก เพื่อให้ Tracking ID ไม่ซ้ำกัน
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(900) + 100 // ได้ค่า 100-999

	// ✅ สร้าง Tracking ID ตามโครงสร้างที่อ่านง่ายและนำไปใช้งานได้
	return fmt.Sprintf("TRK-%s-%s-%d", lotSuffix, retailerSuffix, randomNumber)
}

// ✅ ฟังก์ชันดึงข้อมูล Product Lot Details
func (pc *ProductLotController) GetProductLotDetails(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Product Lot Details")

	// ✅ ดึง `lotId` จาก URL Parameter
	lotID := c.Params("lotId")
	if lotID == "" {
		fmt.Println("❌ Product Lot ID is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product Lot ID is required"})
	}
	fmt.Println("✅ Product Lot ID:", lotID)

	// ✅ ดึงข้อมูล Product Lot จาก Blockchain
	productLotData, err := pc.BlockchainService.GetProductLotByLotID(lotID)
	if err != nil {
		fmt.Println("❌ Failed to fetch product lot from blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product lot details"})
	}

	// ✅ ดึงข้อมูล Product จาก Smart Contract
	productID := productLotData.ProductID
	productData, err := pc.BlockchainService.GetProductDetails(productID)
	if err != nil {
		fmt.Println("❌ Failed to fetch product from blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product details"})
	}

	// ✅ ดึงข้อมูล JSON จาก IPFS ของ Product (เพื่อหา quantityUnit)
	productIPFSCID := productData["productCID"].(string)
	productIPFSData, err := pc.IPFSService.GetJSONFromIPFS(productIPFSCID)
	if err != nil {
		fmt.Println("❌ Failed to fetch product data from IPFS:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product data"})
	}

	// ✅ Debug Log ตรวจสอบว่า `productIPFSData` มีอะไรอยู่บ้าง
	fmt.Println("📌 Debug: productIPFSData =", productIPFSData)

	// ✅ ตรวจสอบว่า Nutrition มีค่าหรือไม่
	NutritionData, ok := productIPFSData["nutrition"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Error: Nutrition data is missing or incorrect")
		fmt.Println("📌 Debug: Available keys in productIPFSData:", reflect.ValueOf(productIPFSData).MapKeys()) // ✅ เช็คว่า key มีอะไรบ้าง
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Nutrition data structure is incorrect"})
	}

	// ✅ ดึงข้อมูล JSON จาก IPFS ของ Product Lot (Quality & Nutrition)
	ipfsCID := productLotData.QualityAndNutritionCID
	ipfsData, err := pc.IPFSService.GetJSONFromIPFS(ipfsCID)
	if err != nil {
		fmt.Println("❌ Failed to fetch quality & nutrition data from IPFS:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch quality & nutrition data"})
	}

	// ✅ Debug Log ตรวจสอบว่า `ipfsData` มีอะไรอยู่บ้าง
	fmt.Println("📌 Debug: ipfsData =", ipfsData)

	// ✅ ตรวจสอบว่า qualityData มีอยู่จริงหรือไม่
	qualityDataMap, ok := ipfsData["qualityData"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Error: qualityData is missing or incorrect")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "qualityData structure is incorrect"})
	}

	// ✅ ดึงข้อมูล quality
	qualityData, ok := qualityDataMap["quality"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Error: Quality data is missing or incorrect")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Quality data structure is incorrect"})
	}

	// ✅ ดึงข้อมูล nutrition (แก้จาก := เป็น =)
	nutritionData, ok := qualityDataMap["nutrition"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Error: Nutrition data is missing or incorrect")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Nutrition data structure is incorrect"})
	}
	// ✅ แปลง `grade` เป็นข้อความ
	var gradeText string
	if productLotData.Grade {
		gradeText = "Passed"
	} else {
		gradeText = "Failed"
	}

	// ✅ แปลง `inspectionDate` เป็น `YYYY-MM-DD HH:mm:ss`
	inspectionTime := time.Unix(productLotData.InspectionDate.Unix(), 0).Format("2006-01-02 15:04:05")
	// ✅ ตรวจสอบว่า `MilkTankIDs` มีค่า หรือไม่ ถ้าไม่มี ให้ตั้งค่าเป็นอาเรย์ว่าง
	var milkTankIDs []string
	if len(productLotData.MilkTankIDs) > 0 {
		milkTankIDs = productLotData.MilkTankIDs
	} else {
		milkTankIDs = []string{} // ตั้งค่าเป็นอาเรย์ว่างเพื่อป้องกัน error
	}
	// ✅ ดึงข้อมูล Tracking จาก Blockchain
	trackingIds, _, qrCodeCIDs, err := pc.BlockchainService.GetTrackingByLotId(lotID)
	if err != nil {
		fmt.Println("❌ Failed to fetch tracking data:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tracking data"})
	}

	// ✅ ตรวจสอบว่ามี Tracking Events หรือไม่
	var trackingDataArray []fiber.Map
	for i := range trackingIds {
		// ✅ ดึงข้อมูล QR Code Data
		qrCodeData, err := pc.QRService.ReadQRCodeFromCID(qrCodeCIDs[i])
		if err != nil {
			fmt.Println("❌ Failed to decode QR Code from CID:", qrCodeCIDs[i])
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode QR Code"})
		}

		// ✅ ดึงภาพ QR Code จาก IPFS
		qrCodeBase64, err := pc.IPFSService.GetImageBase64FromIPFS(qrCodeCIDs[i])
		if err != nil {
			fmt.Println("❌ Failed to fetch QR Code image:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch QR Code image"})
		}

		// ✅ เพิ่มข้อมูล Tracking เข้าไปในอาเรย์
		trackingDataArray = append(trackingDataArray, fiber.Map{
			"trackingId": trackingIds[i],
			"qrCodeData": qrCodeData,
			"qrCodeImg":  fmt.Sprintf("data:image/png;base64,%s", qrCodeBase64),
		})
	}
	// ✅ ส่งค่าที่แก้ไขกลับไปที่ Frontend
	response := fiber.Map{
		"GeneralInfo": fiber.Map{
			"productId":    productID,
			"productName":  productData["productName"],
			"category":     productData["category"],
			"description":  productIPFSData["description"],
			"quantity":     productIPFSData["quantity"],
			"quantityUnit": NutritionData["quantityUnit"],
		},
		"selectMilkTank": fiber.Map{
			"tankIds":         milkTankIDs, // ✅ ใช้ตัวแปรที่ตรวจสอบแล้ว
			"temp":            qualityData["temp"],
			"tempUnit":        qualityData["tempUnit"],
			"pH":              qualityData["pH"],
			"fat":             qualityData["fat"],
			"protein":         qualityData["protein"],
			"bacteria":        qualityData["bacteria"],
			"bacteriaInfo":    qualityData["bacteriaInfo"],
			"contaminants":    qualityData["contaminants"],
			"contaminantInfo": qualityData["contaminantInfo"],
			"abnormalChar":    qualityData["abnormalChar"],
			"abnormalType":    qualityData["abnormalType"],
		},
		"Quality": fiber.Map{
			"grade":          gradeText,
			"inspectionDate": inspectionTime,
			"inspector":      productLotData.Inspector,
		},
		"nutrition":         nutritionData,
		"shippingAddresses": trackingDataArray, // ✅ เพิ่ม Tracking Data เข้าไป

	}

	// ✅ ส่งข้อมูลให้ Frontend
	return c.Status(http.StatusOK).JSON(response)

}

// ✅ ฟังก์ชันดึง Product Lots ทั้งหมดของโรงงาน
func (plc *ProductLotController) GetFactoryProductLots(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Factory Product Lots")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	factoryWallet := c.Locals("walletAddress").(string)

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Factory เท่านั้น)
	if role != "factory" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only factories can view product lots"})
	}

	// ✅ ดึงค่าที่พิมพ์ในช่องค้นหา (Search Query)
	searchQuery := strings.ToLower(c.Query("search", ""))

	// ✅ ดึงข้อมูล Product Lots จาก Blockchain (ได้ค่าครบเลย)
	productLots, err := plc.BlockchainService.GetProductLotsByFactory(factoryWallet)
	if err != nil {
		fmt.Println("❌ Failed to fetch product lots:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product lots"})
	}

	// ✅ กรองข้อมูลตาม Search Query
	var filteredProductLots []map[string]interface{}
	for _, lot := range productLots {
		lotID := lot["Product Lot No"]
		productName := lot["Product Name"]
		personInCharge := lot["Person In Charge"]

		// ✅ ถ้า searchQuery ว่าง → แสดงทั้งหมด, ถ้าไม่ว่าง → ค้นหาตาม Lot ID หรือ Product Name
		if searchQuery == "" || strings.Contains(strings.ToLower(lotID), searchQuery) || strings.Contains(strings.ToLower(productName), searchQuery) {
			filteredProductLots = append(filteredProductLots, map[string]interface{}{
				"productLotNo":   lotID,
				"productName":    productName,
				"personInCharge": personInCharge,
				"moreInfoLink":   fmt.Sprintf("/Factory/ProductLot/Details?id=%s", lotID),
			})
		}
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"displayedProductLots": filteredProductLots,
		"addNewLotLink":        "/Factory/CreateProductLot",
	})
}

func (plc *ProductLotController) GetAllTrackingIds(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get All Tracking IDs")

	// ✅ ดึง Tracking IDs จาก Smart Contract
	trackingIds, err := plc.BlockchainService.GetAllTrackingIds()
	if err != nil {
		fmt.Println("❌ Failed to fetch tracking IDs:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tracking IDs"})
	}

	// ✅ สร้างลิงก์สำหรับดูรายละเอียด Tracking แต่ละอัน
	var trackingList []map[string]interface{}
	for _, trackingId := range trackingIds {
		trackingList = append(trackingList, map[string]interface{}{
			"trackingId":   strings.TrimRight(trackingId, "\x00"),
			"moreInfoLink": fmt.Sprintf("/Tracking/Details?id=%s", trackingId),
		})
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"trackingList":       trackingList,
		"addNewTrackingLink": "/Tracking/Create",
	})
}

func (plc *ProductLotController) UpdateLogisticsCheckpoint(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Update Logistics Checkpoints")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	walletAddress := c.Locals("walletAddress").(string)

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Logistics เท่านั้น)
	if role != "logistics" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only logistics can update checkpoints"})
	}

	// ✅ Debug - แสดงข้อมูล JSON ที่รับเข้ามาก่อนแปลง
	bodyBytes := c.Body()
	fmt.Println("📡 Received Raw JSON Body:", string(bodyBytes))

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		TrackingID  string `json:"trackingId"`
		Checkpoints struct {
			Before []Checkpoint `json:"before"`
			During []Checkpoint `json:"during"`
			After  []Checkpoint `json:"after"`
		} `json:"checkpoints"`
	}

	// ✅ ตรวจสอบ JSON Request
	if err := c.BodyParser(&request); err != nil {
		fmt.Println("❌ Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ Debug - แสดงข้อมูลหลังแปลง JSON สำเร็จ
	fmt.Printf("✅ Parsed Request Data:\nTrackingID: %s\nCheckpoints: %+v\n", request.TrackingID, request.Checkpoints)

	// ✅ ตรวจสอบค่าที่จำเป็น
	if request.TrackingID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing Tracking ID"})
	}

	// ✅ อัปโหลดข้อมูล `ReceiverInfo` ไปยัง IPFS
	uploadToIPFS := func(cp Checkpoint) (string, error) {
		// 🛠 แปลง Struct → map[string]interface{}
		jsonData, err := json.Marshal(cp)
		if err != nil {
			return "", fmt.Errorf("❌ Failed to marshal checkpoint data: %v", err)
		}

		var checkpointMap map[string]interface{}
		if err := json.Unmarshal(jsonData, &checkpointMap); err != nil {
			return "", fmt.Errorf("❌ Failed to unmarshal checkpoint data: %v", err)
		}

		// ✅ Debug ข้อมูลที่กำลังอัปโหลดไป IPFS
		fmt.Println("📡 Uploading Checkpoint Data to IPFS:", checkpointMap)

		// ✅ อัปโหลดไปยัง IPFS
		cid, err := plc.IPFSService.UploadDataToIPFS(checkpointMap)
		if err != nil {
			return "", fmt.Errorf("❌ Failed to upload checkpoint data to IPFS: %v", err)
		}

		// ✅ Debug CID ที่ได้จาก IPFS
		fmt.Println("✅ Uploaded to IPFS, CID:", cid)
		return cid, nil
	}

	// ✅ ประมวลผล Checkpoints
	allCheckpoints := []BlockchainCheckpoint{}
	// ✅ Debug - เช็คค่าที่ส่งมาจาก JSON
	fmt.Println("📌 Received JSON Data:", request)

	processCheckpoints := func(checkpoints []Checkpoint, checkType uint8) error {
		for _, cp := range checkpoints {
			// แปลงเวลา string เป็น Unix Timestamp
			pickupUnix := parseTimeStringToUnix(cp.PickupTime)
			deliveryUnix := parseTimeStringToUnix(cp.DeliveryTime)

			// อัปโหลด Checkpoint data (ทั้งหมด) ไปยัง IPFS
			cid, err := uploadToIPFS(cp)
			if err != nil {
				return err
			}

			// รวม firstName + lastName จาก cp
			personInCharge := cp.FirstName + " " + cp.LastName

			// เพิ่มข้อมูลไปยังรายการที่ส่งไป Blockchain
			allCheckpoints = append(allCheckpoints, BlockchainCheckpoint{
				CheckType:      checkType,
				PickupTime:     uint64(pickupUnix),
				DeliveryTime:   uint64(deliveryUnix),
				Quantity:       uint64(cp.Quantity),
				Temperature:    int64(cp.Temperature),
				PersonInCharge: personInCharge,
				ReceiverCID:    cid,
			})
		}
		return nil
	}

	// ✅ ประมวลผลแต่ละช่วงเวลา
	if err := processCheckpoints(request.Checkpoints.Before, 0); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := processCheckpoints(request.Checkpoints.During, 1); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := processCheckpoints(request.Checkpoints.After, 2); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// ✅ ตรวจสอบว่ามี Checkpoint อย่างน้อยหนึ่งรายการ
	if len(allCheckpoints) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No logistics checkpoints provided"})
	}

	// ✅ วนลูปส่ง Checkpoints ไปยัง Blockchain
	var txHashes []string
	for _, checkpoint := range allCheckpoints {
		txHash, err := plc.BlockchainService.UpdateLogisticsCheckpoint(
			walletAddress,
			request.TrackingID,
			checkpoint.PickupTime,
			checkpoint.DeliveryTime,
			checkpoint.Quantity,
			checkpoint.Temperature,
			checkpoint.PersonInCharge, // ✅ Backend รวมชื่อก่อนบันทึก
			checkpoint.CheckType,
			checkpoint.ReceiverCID,
		)
		if err != nil {
			fmt.Println("❌ Blockchain transaction failed:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
		}

		// ✅ Debug Hash ของ Transaction บน Blockchain
		fmt.Println("✅ Transaction Sent, Hash:", txHash)
		txHashes = append(txHashes, txHash)
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":    "Logistics checkpoints updated successfully",
		"trackingId": request.TrackingID,
		"txHashes":   txHashes,
	})
}

func parseTimeStringToUnix(timeStr string) int64 {
	// 🛠 แปลง "YYYY-MM-DDTHH:MM" เป็น Unix Timestamp
	layout := "2006-01-02T15:04"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("⚠️ Warning: Failed to parse time:", timeStr, "Error:", err)
		return 0 // ❌ ถ้าแปลงไม่ได้ ให้คืนค่าเป็น 0
	}
	return t.Unix()
}

// ✅ Structs ที่ใช้ในโค้ด
type Checkpoint struct {
	PickupTime   string `json:"deliverTime"` // ✅ ใช้ deliverTime จาก JSON
	DeliveryTime string `json:"recieveTime"` // ✅ ใช้ recieveTime จาก JSON
	Quantity     int    `json:"quantity"`    // ✅ ตรงกับ JSON
	Temperature  int    `json:"temp"`        // ✅ ตรงกับ JSON
	CompanyName  string `json:"companyName"` // ✅ ใช้ตรง ๆ จาก JSON
	FirstName    string `json:"firstName"`   // ✅ ใช้ตรง ๆ จาก JSON
	LastName     string `json:"lastName"`    // ✅ ใช้ตรง ๆ จาก JSON
	Email        string `json:"email"`       // ✅ ใช้ตรง ๆ จาก JSON
	Phone        string `json:"phoneNumber"` // ❌ JSON ใช้ phoneNumber แต่ Struct ใช้ Phone → ต้องแก้ให้ตรงกัน
	Address      string `json:"address"`     // ✅ ใช้ตรง ๆ จาก JSON
	Province     string `json:"province"`    // ✅ ใช้ตรง ๆ จาก JSON
	District     string `json:"district"`    // ✅ ใช้ตรง ๆ จาก JSON
	SubDistrict  string `json:"subDistrict"` // ✅ ใช้ตรง ๆ จาก JSON
	PostalCode   string `json:"postalCode"`  // ✅ ใช้ตรง ๆ จาก JSON
	Location     string `json:"location"`    // ✅ ใช้ตรง ๆ จาก JSON
}

type ReceiverInfo struct {
	CompanyName string `json:"companyName"`
	FirstName   string `json:"firstName"` // ✅ ใช้ FirstName + LastName
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phoneNumber"` // ✅ JSON ใช้ "phoneNumber" ต้องตรงกัน	Address     string `json:"address"`
	Province    string `json:"province"`
	District    string `json:"district"`
	SubDistrict string `json:"subDistrict"`
	PostalCode  string `json:"postalCode"`
	Location    string `json:"location"`
}

type BlockchainCheckpoint struct {
	CheckType      uint8
	PickupTime     uint64
	DeliveryTime   uint64
	Quantity       uint64
	Temperature    int64
	PersonInCharge string // ✅ Backend รวม FirstName + LastName
	ReceiverCID    string
}

type LogisticsCheckpoint struct {
	TrackingId        string       `json:"trackingId"`
	LogisticsProvider string       `json:"logisticsProvider"`
	PickupTime        uint64       `json:"pickupTime"`
	DeliveryTime      uint64       `json:"deliveryTime"`
	Quantity          uint64       `json:"quantity"`
	Temperature       int64        `json:"temperature"`
	PersonInCharge    string       `json:"personInCharge"`
	CheckType         uint8        `json:"checkType"`
	ReceiverCID       string       `json:"receiverCID"`
	ReceiverInfo      ReceiverInfo `json:"receiverInfo,omitempty"` // ✅ อัปเดตค่าจาก IPFS
}

// ✅ ฟังก์ชันดึงข้อมูล Logistics Checkpoints ตาม Tracking ID
func (plc *ProductLotController) GetLogisticsCheckpointsByTrackingId(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Logistics Checkpoints by Tracking ID")

	// ✅ รับ Tracking ID จาก Query Parameter
	trackingId := c.Query("trackingId")
	if trackingId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Tracking ID is required"})
	}

	// ✅ ดึงข้อมูลจาก Smart Contract ผ่าน BlockchainService
	beforeCheckpoints, duringCheckpoints, afterCheckpoints, err := plc.BlockchainService.GetLogisticsCheckpointsByTrackingId(trackingId)
	if err != nil {
		fmt.Println("❌ Failed to fetch logistics checkpoints:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch logistics checkpoints"})
	}

	// ✅ ฟังก์ชันดึงข้อมูล ReceiverInfo จาก IPFS
	enhanceCheckpointsWithIPFS := func(checkpoints []services.LogisticsCheckpoint) []map[string]interface{} {
		var enhancedCheckpoints []map[string]interface{}
		for _, cp := range checkpoints {
			fmt.Println("📡 Fetching Receiver Info from IPFS CID:", cp.ReceiverCID)
			ipfsData, err := plc.IPFSService.GetJSONFromIPFS(cp.ReceiverCID)
			if err != nil {
				fmt.Println("⚠️ Warning: Failed to fetch receiver info from IPFS:", err)
				continue
			}
			fmt.Println("✅ IPFS Data:", ipfsData) // ✅ Debug ข้อมูลที่ได้จาก IPFS

			// ✅ แปลง map[string]interface{} เป็น ReceiverInfo
			receiverInfo := map[string]interface{}{
				"companyName": getStringFromMap(ipfsData, "companyName"),
				"firstName":   getStringFromMap(ipfsData, "firstName"),
				"lastName":    getStringFromMap(ipfsData, "lastName"),
				"email":       getStringFromMap(ipfsData, "email"),
				"phone":       getStringFromMap(ipfsData, "phone"),
				"address":     getStringFromMap(ipfsData, "address"),
				"province":    getStringFromMap(ipfsData, "province"),
				"district":    getStringFromMap(ipfsData, "district"),
				"subDistrict": getStringFromMap(ipfsData, "subDistrict"),
				"postalCode":  getStringFromMap(ipfsData, "postalCode"),
				"location":    getStringFromMap(ipfsData, "location"),
			}

			// ✅ สร้าง JSON Response ใหม่
			enhancedCheckpoints = append(enhancedCheckpoints, map[string]interface{}{
				"trackingId":        cp.TrackingId,
				"logisticsProvider": cp.LogisticsProvider,
				"pickupTime":        cp.PickupTime,
				"deliveryTime":      cp.DeliveryTime,
				"quantity":          cp.Quantity,
				"temperature":       cp.Temperature,
				"personInCharge":    cp.PersonInCharge,
				"checkType":         cp.CheckType,
				"receiverCID":       cp.ReceiverCID,
				"receiverInfo":      receiverInfo, // ✅ เพิ่มข้อมูลจาก IPFS
			})
		}
		return enhancedCheckpoints
	}

	// ✅ เพิ่มข้อมูลจาก IPFS ใน Response JSON
	response := fiber.Map{
		"trackingId":        trackingId,
		"beforeCheckpoints": enhanceCheckpointsWithIPFS(beforeCheckpoints),
		"duringCheckpoints": enhanceCheckpointsWithIPFS(duringCheckpoints),
		"afterCheckpoints":  enhanceCheckpointsWithIPFS(afterCheckpoints),
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(response)
}

// ✅ ฟังก์ชันช่วยแปลงค่า map[string]interface{} -> string
func getStringFromMap(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

// /////// ร้านค้า//////////////
func (plc *ProductLotController) GetRetailerTracking(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Tracking Data by Retailer ID")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	retailerID, ok := c.Locals("entityID").(string) // ✅ ดึง `Retailer ID`
	if !ok || retailerID == "" {
		fmt.Println("❌ Retailer ID is missing in Context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - Retailer ID is missing"})
	}

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Retailer เท่านั้น)
	if role != "retailer" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only retailers can view their tracking data"})
	}

	fmt.Println("✅ Retailer ID from Context:", retailerID)

	// ✅ ดึง Tracking IDs ตาม Retailer ID จาก Blockchain
	trackingData, err := plc.BlockchainService.GetTrackingByRetailer(retailerID)
	if err != nil {
		fmt.Println("❌ Failed to fetch tracking data for retailer:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tracking data"})
	}

	// ✅ ตรวจสอบข้อมูลที่ได้รับ
	if len(trackingData) == 0 {
		fmt.Println("⚠️ No tracking data found for retailer:", retailerID)
		return c.Status(http.StatusOK).JSON(fiber.Map{"trackingList": []map[string]interface{}{}})
	}

	// ✅ กรองข้อมูล - แปลง Tracking ID ให้ถูกต้อง
	var trackingList []map[string]interface{}
	for _, item := range trackingData {
		trackingIDRaw, exists := item["trackingId"]
		if !exists {
			fmt.Println("⚠️ Skipping entry with missing tracking ID:", item)
			continue
		}

		// ✅ ตัดอักขระ \x00 และแปลงเป็น string
		trackingID, ok := trackingIDRaw.(string)
		if !ok {
			fmt.Println("❌ Invalid tracking ID format:", trackingIDRaw)
			continue
		}
		trackingID = strings.TrimSpace(trackingID)

		// ✅ ตรวจสอบว่า Tracking ID เป็นค่าว่างหรือไม่
		if trackingID == "" {
			fmt.Println("⚠️ Empty tracking ID found, skipping...")
			continue
		}

		// ✅ สร้างโครงสร้าง Response
		trackingList = append(trackingList, map[string]interface{}{
			"trackingId":   trackingID,
			"moreInfoLink": fmt.Sprintf("/Retailer/TrackingDetails?id=%s", trackingID),
		})
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"trackingList": trackingList,
	})
}

func (plc *ProductLotController) RetailerReceiveProduct(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Retailer Receiving Product")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	walletAddress := c.Locals("walletAddress").(string)
	fmt.Println("📌 Debug - Wallet Address:", walletAddress)
	entityId := c.Locals("entityID").(string) // ✅ ใช้ Entity ID ของ Retailer
	// ✅ ตรวจสอบสิทธิ์
	if role != "retailer" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only retailers can receive products"})
	}

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		TrackingId string `json:"trackingId"`
		Input      struct {
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
	if request.TrackingId == "" || request.Input.RecipientInfo.PersonInCharge == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// ✅ รวมข้อมูลอัปโหลดไป IPFS
	productMetadata := map[string]interface{}{
		"trackingId": request.TrackingId,
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

	// ✅ อัปโหลดข้อมูลไปยัง IPFS (ใช้ `UploadDataToIPFS` แทน)
	qualityReportCID, err := plc.IPFSService.UploadDataToIPFS(productMetadata)
	if err != nil {
		fmt.Println("❌ Failed to upload to IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload quality report"})
	}

	// ✅ อัปเดตสถานะไปยัง Blockchain
	txHash, err := plc.BlockchainService.RetailerReceiveProduct(
		walletAddress,
		request.TrackingId,
		entityId,
		qualityReportCID,
		request.Input.RecipientInfo.PersonInCharge,
	)
	if err != nil {
		fmt.Println("❌ Blockchain transaction failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":          "Product received successfully",
		"trackingId":       request.TrackingId,
		"txHash":           txHash,
		"qualityReportCID": qualityReportCID,
	})
}

func (plc *ProductLotController) GetRetailerReceivedProduct(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Retailer Received Product")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	entityId := c.Locals("entityID").(string)

	// ✅ ตรวจสอบสิทธิ์
	if role != "retailer" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only retailers can access this data"})
	}

	// ✅ รับ Tracking ID จาก Query Parameter
	trackingId := c.Query("trackingId")
	if trackingId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Tracking ID is required"})
	}

	fmt.Println("📌 Debug - Fetching Data for Tracking ID:", trackingId)

	// ✅ ดึงข้อมูลจาก Blockchain
	retailerData, err := plc.BlockchainService.GetRetailerConfirmation(trackingId)
	if err != nil {
		fmt.Println("❌ Failed to fetch retailer confirmation:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch retailer confirmation"})
	}

	// ✅ ตรวจสอบว่า Entity ID ตรงกับ Retailer ID ใน Blockchain หรือไม่
	if retailerData["retailerId"] != entityId {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: This retailer is not assigned to this tracking ID"})
	}

	// ✅ ดึงข้อมูลจาก IPFS โดยใช้ Quality CID
	qualityCID, ok := retailerData["qualityCID"].(string)
	if !ok || qualityCID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid or missing qualityCID"})
	}

	qualityData, err := plc.IPFSService.GetJSONFromIPFS(qualityCID)
	if err != nil {
		fmt.Println("❌ Failed to fetch quality report from IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve quality data from IPFS"})
	}

	fmt.Println("📌 Debug - Quality Data from IPFS:", qualityData)

	// ✅ ตรวจสอบโครงสร้าง `recipientInfo`
	recipientInfo, ok := qualityData["recipientInfo"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Missing or invalid recipientInfo structure")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid recipient info structure"})
	}

	// ✅ ดึงข้อมูล `Quantity`
	quantity, _ := qualityData["quantity"].(float64)
	quantityUnit, _ := qualityData["quantityUnit"].(string)
	temp, _ := qualityData["temperature"].(float64)
	tempUnit, _ := qualityData["tempUnit"].(string)
	pH, _ := qualityData["pH"].(float64)
	fat, _ := qualityData["fat"].(float64)
	protein, _ := qualityData["protein"].(float64)
	bacteria, _ := qualityData["bacteria"].(bool)
	bacteriaInfo, _ := qualityData["bacteriaInfo"].(string)
	contaminants, _ := qualityData["contaminants"].(bool)
	contaminantInfo, _ := qualityData["contaminantInfo"].(string)
	abnormalChar, _ := qualityData["abnormalChar"].(bool)
	abnormalType, _ := qualityData["abnormalType"].(map[string]interface{})

	// ✅ ป้องกันค่าที่เป็น `nil`
	if abnormalType == nil {
		abnormalType = map[string]interface{}{}
	}

	// ✅ จัดรูปแบบข้อมูลให้ตรงกับโครงสร้าง JSON ที่รับเข้ามา
	response := fiber.Map{
		"trackingId": trackingId,
		"Input": fiber.Map{
			"RecipientInfo": fiber.Map{
				"personInCharge": recipientInfo["personInCharge"],
				"location":       recipientInfo["location"],
				"pickUpTime":     recipientInfo["pickUpTime"],
			},
			"Quantity": fiber.Map{
				"quantity":        quantity,
				"quantityUnit":    quantityUnit,
				"temp":            temp,
				"tempUnit":        tempUnit,
				"pH":              pH,
				"fat":             fat,
				"protein":         protein,
				"bacteria":        bacteria,
				"bacteriaInfo":    bacteriaInfo,
				"contaminants":    contaminants,
				"contaminantInfo": contaminantInfo,
				"abnormalChar":    abnormalChar,
				"abnormalType":    abnormalType,
			},
		},
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(http.StatusOK).JSON(response)
}
