package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
