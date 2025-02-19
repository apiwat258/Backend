package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"finalyearproject/Backend/services"

	"github.com/ethereum/go-ethereum/common"
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

// generateRawMilkID - สร้าง RawMilkID โดยให้ Blockchain & UI ใช้ ID เดียวกัน
func generateRawMilkID(farmWallet string) (string, [32]byte) {
	// ดึงวันที่ปัจจุบันในรูปแบบ YYYYMMDD
	currentDate := time.Now().Format("20060102")

	// ใช้ SHA-256 Hash เพื่อให้ ID มีความเป็นเอกลักษณ์
	hashInput := fmt.Sprintf("%s-%s-%d", farmWallet, currentDate, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(hashInput))

	// ✅ ใช้ 16 ตัวอักษรแรกสำหรับ UI
	shortID := hex.EncodeToString(hash[:])[:16]

	// ✅ คืนค่า 16-char ID + bytes32 Hash
	return shortID, hash
}

// AddRawMilkHandler รับข้อมูลจากฟาร์มและบันทึกลง Blockchain
func AddRawMilkHandler(c *fiber.Ctx) error {
	var request RawMilkRequest

	// ✅ แปลง JSON request เป็น struct
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// ✅ สร้าง RawMilkID ทั้งแบบ 16-char (UI) และ bytes32 (Blockchain)
	rawMilkShortID, rawMilkHash := generateRawMilkID(request.FarmWallet)
	fmt.Printf("✅ Generated RawMilkID: %s\n", rawMilkShortID)

	// ✅ Debug: Log ค่าที่ได้รับ
	fmt.Printf("Received Raw Milk Data: %+v\n", request)

	// ✅ เช็คว่า Blockchain Service ทำงานอยู่หรือไม่
	if services.BlockchainServiceInstance == nil {
		log.Println("❌ Blockchain service is not initialized")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain service is not initialized"})
	}

	// ✅ ส่งข้อมูลไปยัง Blockchain (ใช้ rawMilkHash เป็น bytes32)
	txHash, err := services.BlockchainServiceInstance.StoreRawMilkOnBlockchain(
		rawMilkHash, // ✅ ใช้ bytes32
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
		"rawMilkID": rawMilkShortID, // ✅ UI ใช้ 16-char ID
	})
}

// GetRawMilkHandler - ดึงข้อมูล Raw Milk จาก Blockchain
func GetRawMilkHandler(c *fiber.Ctx) error {
	rawMilkID := c.Params("id") // ✅ รับ 16-char ID จาก URL

	if services.BlockchainServiceInstance == nil {
		log.Println("❌ Blockchain service is not initialized")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain service is not initialized"})
	}

	// ✅ แปลง 16-char ID เป็น SHA-256 (`bytes32`)
	fullHash := sha256.Sum256([]byte(rawMilkID))

	// ✅ แปลง `[32]byte` → `common.Hash`
	fullHashCommon := common.BytesToHash(fullHash[:])

	// ✅ ดึงข้อมูลจาก Blockchain
	rawMilk, err := services.BlockchainServiceInstance.GetRawMilkFromBlockchain(fullHashCommon)
	if err != nil {
		log.Println("❌ Failed to fetch raw milk data:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch raw milk data"})
	}

	// ✅ ส่ง response กลับไป (UI ใช้ 16-char ID)
	return c.JSON(fiber.Map{
		"rawMilkID":   rawMilkID, // ✅ UI ใช้ 16-char ID
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

// GenerateQRCodeHandler - API สำหรับสร้าง QR Code
func GenerateQRCodeHandler(c *fiber.Ctx) error {
	rawMilkID := c.Params("id") // ✅ รับ rawMilkID จาก URL

	// ✅ ตรวจสอบว่า Blockchain Service ทำงานอยู่หรือไม่
	if services.BlockchainServiceInstance == nil {
		log.Println("❌ Blockchain service is not initialized")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain service is not initialized"})
	}

	// ✅ แปลง rawMilkID เป็น Hash สำหรับดึงข้อมูลจาก Blockchain
	rawMilkHash := utils.GenerateHash(rawMilkID)

	// ✅ ดึงข้อมูลจาก Blockchain
	rawMilk, err := services.BlockchainServiceInstance.GetRawMilkFromBlockchain(rawMilkHash)
	if err != nil {
		log.Println("❌ Failed to fetch raw milk data:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch raw milk data"})
	}

	// ✅ สร้าง JSON Data สำหรับ QR Code
	qrData := map[string]interface{}{
		"rawMilkID":   rawMilkID,
		"farmID":      rawMilk.FarmWallet, // ❌ ต้องแก้เป็น FarmID ถ้ามีใน Blockchain
		"milkTankNum": rawMilk.MilkTankNum,
		"ipfsCid":     rawMilk.IPFSCid,
	}

	qrJSON, err := json.Marshal(qrData)
	if err != nil {
		log.Println("❌ Failed to create QR Code JSON:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create QR Code JSON"})
	}

	// ✅ ใช้ QR Code Service สร้าง QR Code (ยังไม่เพิ่ม ฟังก์ชัน)
	qrCodeImage, err := services.GenerateQRCode(string(qrJSON))
	if err != nil {
		log.Println("❌ Failed to generate QR Code:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate QR Code"})
	}

	// ✅ ส่ง QR Code กลับไปเป็น Base64
	return c.JSON(fiber.Map{
		"message": "QR Code generated successfully",
		"qrCode":  qrCodeImage,
	})
}
