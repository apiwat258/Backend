package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"finalyearproject/Backend/services"

	"github.com/gofiber/fiber/v2"
)

// RawMilkRequest ใช้รับค่าจากฟอร์มของฟาร์ม
type RawMilkRequest struct {
	FarmWallet      string   `json:"farmWallet"`
	FarmName        string   `json:"farmName"`
	MilkTankNum     int      `json:"milkTankNum"`
	PersonInCharge  string   `json:"personInCharge"`
	Quantity        float64  `json:"quantity"`
	QuantityUnit    string   `json:"quantityUnit"`
	Temperature     float64  `json:"temperature"`
	TemperatureUnit string   `json:"temperatureUnit"`
	PH              float64  `json:"pH"`
	Fat             float64  `json:"fat"`
	Protein         float64  `json:"protein"`
	BacteriaTest    string   `json:"bacteriaTest,omitempty"`
	Contaminants    string   `json:"contaminants,omitempty"`
	AbnormalChecks  []string `json:"abnormalChecks,omitempty"`
	Location        string   `json:"location"`
	IPFSCid         string   `json:"ipfsCid"`
}

// generateRawMilkID - ฟังก์ชันสร้าง RawMilkID อัตโนมัติ
func generateRawMilkID(farmWallet string) string {
	// ดึงวันที่ปัจจุบันในรูปแบบ YYYYMMDD
	currentDate := time.Now().Format("20060102")

	// ใช้ SHA-256 Hash เพื่อให้ ID มีความเป็นเอกลักษณ์
	hashInput := fmt.Sprintf("%s-%s-%d", farmWallet, currentDate, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(hashInput))

	// แปลงค่า Hash เป็น Hex และตัดให้เหลือ 16 ตัวอักษร
	rawMilkID := hex.EncodeToString(hash[:])[:16] // 🛑 ลดความยาวให้ไม่ยาวเกินไป

	return rawMilkID
}

// AddRawMilkHandler รับข้อมูลจากฟาร์มและบันทึกลง Blockchain
func AddRawMilkHandler(c *fiber.Ctx) error {
	var request RawMilkRequest

	// ✅ ใช้ c.BodyParser เพื่อแปลง JSON request เป็น struct
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// ✅ สร้าง RawMilkID อัตโนมัติ
	rawMilkID := generateRawMilkID(request.FarmWallet)
	fmt.Printf("✅ Generated RawMilkID: %s\n", rawMilkID)

	// ✅ Debug: Log ค่าที่ได้รับ
	fmt.Printf("Received Raw Milk Data: %+v\n", request)

	// ✅ เช็คว่ามี BlockchainServiceInstance หรือไม่
	if services.BlockchainServiceInstance == nil {
		log.Println("❌ Blockchain service is not initialized")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain service is not initialized"})
	}

	// ✅ ส่งข้อมูลไปยัง Blockchain
	txHash, err := services.BlockchainServiceInstance.StoreRawMilkOnBlockchain(
		rawMilkID, // ✅ ใช้ RawMilkID ที่สร้างอัตโนมัติ
		request.FarmWallet,
		request.Temperature,
		request.PH,
		request.Fat,
		request.Protein,
		request.IPFSCid,
	)
	if err != nil {
		log.Println("❌ Failed to store raw milk on blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store raw milk on blockchain"})
	}

	// ✅ ส่ง response กลับไป
	return c.JSON(fiber.Map{
		"message":   "Raw milk data stored on blockchain",
		"txHash":    txHash,
		"rawMilkID": rawMilkID, // ✅ ส่งค่า RawMilkID กลับไป
	})
}

// GetRawMilkHandler - ดึงข้อมูล Raw Milk จาก Blockchain
func GetRawMilkHandler(c *fiber.Ctx) error {
	rawMilkID := c.Params("id") // ✅ รับ rawMilkID จาก URL

	// ✅ ตรวจสอบว่า BlockchainService ถูก initialize หรือไม่
	if services.BlockchainServiceInstance == nil {
		log.Println("❌ Blockchain service is not initialized")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain service is not initialized"})
	}

	// ✅ ดึงข้อมูลจาก Blockchain
	rawMilk, err := services.BlockchainServiceInstance.GetRawMilkFromBlockchain(rawMilkID)
	if err != nil {
		log.Println("❌ Failed to fetch raw milk data:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch raw milk data"})
	}

	// ✅ ส่ง response กลับไป
	return c.JSON(fiber.Map{
		"farmWallet":  rawMilk.FarmWallet,
		"temperature": rawMilk.Temperature,
		"pH":          rawMilk.PH,
		"fat":         rawMilk.Fat,
		"protein":     rawMilk.Protein,
		"ipfsCid":     rawMilk.IPFSCid,
		"status":      rawMilk.Status,
		"timestamp":   rawMilk.Timestamp,
	})
}
