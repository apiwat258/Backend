package controllers

import (
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"
	"fmt"
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
		FarmerID          string `json:"farmerid"`
		CertificationType string `json:"certificationtype"`
		CertificationCID  string `json:"certificationcid"`
		IssuedDate        string `json:"issued_date"`
	}

	var req CertRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("❌ Error parsing request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	fmt.Println("📌 Received Certification Request:", req)

	// ✅ ตรวจสอบว่า Farmer ID มีอยู่จริงหรือไม่
	var farmer models.Farmer
	if err := database.DB.Where("farmerid = ?", req.FarmerID).First(&farmer).Error; err != nil {
		fmt.Println("❌ Farmer ID not found:", req.FarmerID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farmer ID not found"})
	}

	// ✅ ตรวจสอบค่า CID
	if req.CertificationCID == "" {
		fmt.Println("❌ Certification CID is missing!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certification CID is required"})
	}

	// ✅ แปลงวันที่จาก `string` → `time.Time`
	var issuedDate time.Time
	var err error
	if req.IssuedDate != "" {
		issuedDate, err = time.Parse("2006-01-02", req.IssuedDate)
		if err != nil {
			fmt.Println("❌ Invalid date format:", req.IssuedDate)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format. Use YYYY-MM-DD"})
		}
	} else {
		issuedDate = time.Time{} // ใช้ `zero time` (NULL)
	}

	// ✅ สร้าง Certification ID ใหม่
	certID := fmt.Sprintf("CERT-%d", time.Now().Unix())

	// ✅ บันทึกลง Database
	certification := models.Certification{
		CertificationID:   certID,
		FarmerID:          req.FarmerID,
		CertificationType: req.CertificationType,
		CertificationCID:  req.CertificationCID, // ✅ ใช้ CID ที่ได้จาก IPFS
		EffectiveDate:     time.Now(),
		IssuedDate:        issuedDate,
		CreatedOn:         time.Now(),
	}

	if err := database.DB.Create(&certification).Error; err != nil {
		fmt.Println("❌ Error saving certification:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save certification"})
	}

	fmt.Println("✅ Certification saved:", certification)

	// ✅ ส่งข้อมูลกลับไปยัง Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Certification saved successfully",
		"certification_id": certification.CertificationID,
		"cid":              certification.CertificationCID,
	})
}
