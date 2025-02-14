package controllers

import (
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"
	"fmt"
	"math/big"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ✅ สร้างอินสแตนซ์ของ IPFSService (ที่หายไป)
var ipfsService = services.NewIPFSService()

// ✅ API: อัปโหลดไฟล์ไปยัง IPFS และส่ง CID กลับไปยัง Frontend
func UploadCertificate(c *fiber.Ctx) error {
	fmt.Println("📌 UploadCertificate API called...")

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("❌ No file received")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File not received"})
	}

	fmt.Println("✅ File received:", file.Filename)

	// ✅ เปิดไฟล์
	src, err := file.Open()
	if err != nil {
		fmt.Println("❌ Failed to open file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	// ✅ อัปโหลดไปยัง IPFS ผ่าน `ipfsService`
	cid, err := ipfsService.UploadFile(src)
	if err != nil {
		fmt.Println("❌ Failed to upload to IPFS:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload to IPFS"})
	}

	fmt.Println("✅ Uploaded file to IPFS with CID:", cid)

	// ✅ ส่ง CID กลับไปยัง Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"cid":     cid,
	})
}

// ✅ ฟังก์ชันสร้าง Certification และบันทึกลง Blockchain + PostgreSQL
func CreateCertification(c *fiber.Ctx) error {
	fmt.Println("📌 CreateCertification API called...")

	type CertRequest struct {
		EntityType       string `json:"entity_type"` // Farmer, Factory, Retailer, Logistics
		EntityID         string `json:"entity_id"`   // ID ของหน่วยงาน
		CertificationCID string `json:"certification_cid"`
		IssuedDate       string `json:"issued_date"`
		ExpiryDate       string `json:"expiry_date"`
	}

	var req CertRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("❌ Error parsing request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	fmt.Println("📌 Received Certification Request:", req)

	// ✅ ตรวจสอบค่า CID
	if req.CertificationCID == "" {
		fmt.Println("❌ Certification CID is missing!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certification CID is required"})
	}

	// ✅ แปลงวันที่จาก `string` → `time.Time`
	issuedDate, err := time.Parse("2006-01-02", req.IssuedDate)
	if err != nil {
		fmt.Println("❌ Invalid issued date format:", req.IssuedDate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid issued date format. Use YYYY-MM-DD"})
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		fmt.Println("❌ Invalid expiry date format:", req.ExpiryDate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid expiry date format. Use YYYY-MM-DD"})
	}

	// ✅ แปลง `time.Time` → `*big.Int` (Unix Timestamp)
	issuedDateBigInt := big.NewInt(issuedDate.Unix())
	expiryDateBigInt := big.NewInt(expiryDate.Unix())

	// ✅ สร้าง Certification Event ID ใหม่
	eventID := fmt.Sprintf("EVENT-%d", time.Now().Unix())

	// ✅ บันทึกลง Blockchain (ส่งค่าเป็น *big.Int)
	txHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(eventID, req.EntityType, req.EntityID, req.CertificationCID, issuedDateBigInt, expiryDateBigInt)
	if err != nil {
		fmt.Println("❌ Error storing certification event on blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store certification on blockchain"})
	}

	fmt.Println("✅ Certification Event stored on Blockchain:", txHash)

	// ✅ บันทึกข้อมูลลง PostgreSQL (supplychain_db)
	certification := models.Certification{
		CertificationID:   eventID,
		EntityType:        req.EntityType,
		EntityID:          req.EntityID,
		CertificationType: "Organic", // สมมติว่าเป็น Organic Certification
		CertificationCID:  req.CertificationCID,
		IssuedDate:        issuedDate,
		EffectiveDate:     expiryDate,
		BlockchainTxHash:  txHash,
		CreatedOn:         time.Now(),
	}

	// ✅ เช็คว่าบันทึกลงฐานข้อมูลสำเร็จหรือไม่
	if err := database.DB.Create(&certification).Error; err != nil {
		fmt.Println("❌ Error saving certification to database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save certification to database"})
	}

	fmt.Println("✅ Certification saved to PostgreSQL:", certification)

	// ✅ ส่งข้อมูลกลับไปยัง Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Certification event saved successfully",
		"event_id":      eventID,
		"cid":           req.CertificationCID,
		"blockchain_tx": txHash,
	})
}
