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

func CreateCertification(c *fiber.Ctx) error {
	fmt.Println("📌 CreateCertification API called...")

	type CertRequest struct {
		EntityType       string `json:"entity_type"`
		EntityID         string `json:"entity_id"`
		CertificationCID string `json:"certification_cid"`
		IssuedDate       string `json:"issued_date"`
		ExpiryDate       string `json:"expiry_date"`
	}

	var req CertRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	fmt.Println("📌 Received Certification Request:", req)

	// ✅ ตรวจสอบค่า CID
	if req.CertificationCID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certification CID is required"})
	}

	// ✅ แปลงวันที่จาก `string` → `time.Time`
	issuedDate, err := time.Parse("2006-01-02", req.IssuedDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid issued date format. Use YYYY-MM-DD"})
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid expiry date format. Use YYYY-MM-DD"})
	}

	// ✅ แปลง `time.Time` → `*big.Int`
	issuedDateBigInt := big.NewInt(issuedDate.Unix())
	expiryDateBigInt := big.NewInt(expiryDate.Unix())

	// ✅ ตรวจสอบว่ามีใบเซอร์อยู่แล้วหรือไม่ (จาก Blockchain)
	eventID := fmt.Sprintf("EVENT-%s", req.EntityID)
	existingCert, err := services.BlockchainServiceInstance.GetCertificationFromBlockchain(eventID)

	if err == nil {
		// ✅ ถ้าใบเซอร์เก่ามีอยู่ และยังไม่หมดอายุ → อัปเดตแทน
		if existingCert.IssuedDate.After(time.Now()) {
			fmt.Println("📌 Updating existing certification on Blockchain...")

			txHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(
				eventID,
				req.EntityType,
				req.EntityID,
				req.CertificationCID,
				issuedDateBigInt,
				expiryDateBigInt,
			)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update certification on blockchain"})
			}

			fmt.Println("✅ Certification updated on Blockchain:", txHash)
			return c.JSON(fiber.Map{
				"message":       "Certification updated successfully",
				"event_id":      eventID,
				"cid":           req.CertificationCID,
				"blockchain_tx": txHash,
			})
		}
	}

	// ✅ ถ้าไม่มีใบเซอร์เก่าหรือหมดอายุ → สร้างใหม่
	txHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(
		eventID,
		req.EntityType,
		req.EntityID,
		req.CertificationCID,
		issuedDateBigInt,
		expiryDateBigInt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store certification on blockchain"})
	}

	// ✅ บันทึกข้อมูลลง PostgreSQL
	certification := models.Certification{
		CertificationID:   eventID,
		EntityType:        req.EntityType,
		EntityID:          req.EntityID,
		CertificationType: "Organic",
		CertificationCID:  req.CertificationCID,
		IssuedDate:        issuedDate,
		EffectiveDate:     expiryDate,
		BlockchainTxHash:  txHash,
		CreatedOn:         time.Now(),
	}

	if err := database.DB.Create(&certification).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save certification to database"})
	}

	fmt.Println("✅ Certification saved to PostgreSQL:", certification)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Certification event saved successfully",
		"event_id":      eventID,
		"cid":           req.CertificationCID,
		"blockchain_tx": txHash,
	})
}

func GetCertificationByEntity(c *fiber.Ctx) error {
	entityID := c.Params("entityID")
	if entityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing entity ID"})
	}

	// ✅ ค้นหาข้อมูลใบเซอร์จาก Blockchain
	eventID := fmt.Sprintf("EVENT-%s", entityID)
	certification, err := services.BlockchainServiceInstance.GetCertificationFromBlockchain(eventID)
	if err != nil {
		fmt.Println("❌ [GetCertification] Failed to fetch from Blockchain, trying database...")

		// ✅ ถ้า Blockchain ไม่มีข้อมูล ลองดึงจาก PostgreSQL
		var cert models.Certification
		if err := database.DB.Where("entity_id = ?", entityID).First(&cert).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No certification found"})
		}

		// ✅ ส่งข้อมูลใบเซอร์จาก Database กลับไป
		return c.JSON(fiber.Map{
			"event_id":      cert.CertificationID,
			"entity_type":   cert.EntityType,
			"entity_id":     cert.EntityID,
			"cid":           cert.CertificationCID,
			"issued_date":   cert.IssuedDate.Format("2006-01-02"),
			"expiry_date":   cert.EffectiveDate.Format("2006-01-02"),
			"blockchain_tx": cert.BlockchainTxHash,
		})
	}

	// ✅ ส่งข้อมูลใบเซอร์จาก Blockchain กลับไป
	return c.JSON(certification)
}
