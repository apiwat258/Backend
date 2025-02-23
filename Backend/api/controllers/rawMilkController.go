package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"finalyearproject/Backend/services"
	"finalyearproject/Backend/utils"

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
	request.IPFSCid = "" // ตั้งค่าเริ่มต้นให้ว่างก่อนอัปโหลด

	// ✅ Debug: Log ค่าที่ได้รับ
	fmt.Printf("Received Raw Milk Data: %+v\n", request)

	// ✅ แปลง JSON เป็นไฟล์แล้วอัปโหลดไป IPFS
	requestData, err := json.Marshal(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate JSON"})
	}

	ipfsService := services.NewIPFSService()
	ipfsCid, err := ipfsService.UploadFile(bytes.NewReader(requestData))
	if err != nil {
		log.Println("❌ Failed to upload file to IPFS:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file to IPFS"})
	}

	fmt.Printf("✅ Raw Milk JSON uploaded to IPFS: %s\n", ipfsCid)

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
		ipfsCid, // ✅ ใช้ CID ที่อัปโหลดจาก IPFS
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
		"ipfsCid":   ipfsCid,        // ✅ ส่ง CID กลับให้ผู้ใช้
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
		"rawMilkID": rawMilkID,
		"farmID":    rawMilk.FarmWallet, // ❌ ต้องแก้เป็น FarmID ถ้ามีใน Blockchain
		"ipfsCid":   rawMilk.IPFSCid,
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

// UploadRawMilkFileHandler - API สำหรับอัปโหลดไฟล์ Raw Milk ไปยัง IPFS
func UploadRawMilkFileHandler(c *fiber.Ctx) error {
	var request RawMilkRequest

	// ✅ แปลง JSON request เป็น struct
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// ✅ สร้าง RawMilkID
	rawMilkShortID, _ := generateRawMilkID(request.FarmWallet)
	request.IPFSCid = "" // ตั้งค่าเริ่มต้นให้ว่างก่อนอัปโหลด

	// ✅ เพิ่ม RawMilkID ลงใน JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate JSON"})
	}

	// ✅ อัปโหลด JSON ไปยัง IPFS
	ipfsService := services.NewIPFSService()
	ipfsCid, err := ipfsService.UploadFile(bytes.NewReader(requestData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file to IPFS"})
	}

	fmt.Printf("✅ Raw Milk JSON uploaded to IPFS: %s\n", ipfsCid)

	// ✅ คืนค่า `CID` และ `RawMilkID`
	return c.JSON(fiber.Map{
		"message":   "Raw milk JSON uploaded to IPFS",
		"rawMilkID": rawMilkShortID,
		"ipfsCid":   ipfsCid,
	})
}

// GetRawMilkFromIPFSHandler - ดึงข้อมูล Raw Milk JSON จาก IPFS
func GetRawMilkFromIPFSHandler(c *fiber.Ctx) error {
	ipfsCid := c.Params("cid") // ✅ รับ CID จาก URL

	if ipfsCid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "IPFS CID is required"})
	}

	// ✅ ดึงข้อมูลจาก IPFS
	ipfsService := services.NewIPFSService()
	rawMilkJSON, err := ipfsService.GetFile(ipfsCid)
	if err != nil {
		log.Println("❌ Failed to retrieve file from IPFS:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve file from IPFS"})
	}

	fmt.Printf("✅ Raw Milk JSON retrieved from IPFS: %s\n", ipfsCid)

	// ✅ คืน JSON ให้ผู้ใช้
	return c.JSON(fiber.Map{
		"message": "Raw milk JSON retrieved successfully",
		"ipfsCid": ipfsCid,
		"data":    json.RawMessage(rawMilkJSON),
	})
}
