package controllers

import (
	"encoding/json"
	"fmt"
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
}

func NewProductLotController(db *gorm.DB, blockchainService *services.BlockchainService, ipfsService *services.IPFSService) *ProductLotController {
	return &ProductLotController{
		DB:                db,
		BlockchainService: blockchainService,
		IPFSService:       ipfsService,
	}
}

// ✅ ฟังก์ชันสร้าง Product Lot
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
}

// ✅ ฟังก์ชันสร้าง Lot ID (ใช้ Factory ID)
func (plc *ProductLotController) generateLotID(factoryID string) string {
	return fmt.Sprintf("LOT-%s-%d", factoryID, time.Now().Unix())
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

	// ✅ จัดรูปแบบข้อมูลที่ส่งไป Frontend
	response := fiber.Map{
		"GeneralInfo": fiber.Map{
			"productId":    productID,
			"productName":  productData["productName"],
			"category":     productData["category"],
			"description":  productIPFSData["description"],
			"quantity":     productIPFSData["quantity"],
			"quantityUnit": NutritionData["quantityUnit"], // ✅ ใช้จาก IPFS ของ Product
		},
		"selectMilkTank": fiber.Map{
			"tankIds":         productLotData.MilkTankIDs,
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
			"inspectionDate": inspectionTime, // ✅ แปลงเป็น Timestamp
			"inspector":      productLotData.Inspector,
		},
		"nutrition": nutritionData, // ✅ ใช้ nutritionData ที่ถูกต้อง
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
