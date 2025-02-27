package controllers

import (
	"finalyearproject/Backend/services"
	certification "finalyearproject/Backend/services/certification_event"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ✅ สร้างอินสแตนซ์ของ IPFSService (ที่หายไป)
var ipfsService = services.NewIPFSService()

// ✅ API: อัปโหลดไฟล์ไปยัง IPFS และส่ง CID กลับไปยัง Frontend
func UploadCertificate(c *fiber.Ctx) error {
	fmt.Println("📌 UploadCertificate API called...")

	// ✅ ตรวจสอบ `Content-Type`
	contentType := c.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// ✅ รับไฟล์จาก `multipart/form-data`
		file, err := c.FormFile("file")
		if err != nil {
			fmt.Println("❌ No file received")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File not received"})
		}

		fmt.Println("✅ File received:", file.Filename)

		// ✅ ตรวจสอบประเภทไฟล์
		allowedExtensions := []string{".pdf", ".jpg", ".png"}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !contains(allowedExtensions, ext) {
			fmt.Println("❌ Unsupported file type:", ext)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported file type"})
		}

		// ✅ เปิดไฟล์
		src, err := file.Open()
		if err != nil {
			fmt.Println("❌ Failed to open file:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
		}
		defer src.Close()

		// ✅ อัปโหลดไปยัง IPFS
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

	} else if strings.Contains(contentType, "application/json") {
		// ✅ รองรับ `Base64`
		var req struct {
			Base64Data string `json:"base64"`
			Filename   string `json:"filename"`
		}

		if err := c.BodyParser(&req); err != nil {
			fmt.Println("❌ Invalid JSON format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		fmt.Println("✅ Base64 file received:", req.Filename)

		// ✅ ตรวจสอบประเภทไฟล์
		allowedExtensions := []string{".pdf", ".jpg", ".png"}
		ext := strings.ToLower(filepath.Ext(req.Filename))
		if !contains(allowedExtensions, ext) {
			fmt.Println("❌ Unsupported file type:", ext)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported file type"})
		}

		// ✅ อัปโหลดไปยัง IPFS
		cid, err := ipfsService.UploadBase64File(req.Base64Data)
		if err != nil {
			fmt.Println("❌ Failed to upload Base64 file to IPFS:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload Base64 file to IPFS"})
		}

		fmt.Println("✅ Base64 file uploaded to IPFS with CID:", cid)

		// ✅ ส่ง CID กลับไปยัง Frontend
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "File uploaded successfully",
			"cid":     cid,
		})
	}

	// ❌ ไม่รองรับ `Content-Type` อื่น ๆ
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported Content-Type"})
}

// ✅ ฟังก์ชันตรวจสอบประเภทไฟล์
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// ✅ API: สร้างใบเซอร์ใหม่ใน Blockchain
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

	// ✅ ตรวจสอบว่า `EntityID` มีอยู่จริงในฐานข้อมูลหรือไม่
	var farmer models.Farmer
	if err := database.DB.Where("entityid = ?", req.EntityID).First(&farmer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Entity ID not found in database"})
	}

	// ✅ ตรวจสอบค่า CID
	if req.CertificationCID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certification CID is required"})
	}

	// ✅ ตรวจสอบว่า `certificationCID` มีอยู่ใน Blockchain แล้วหรือไม่
	cidExists, err := services.BlockchainServiceInstance.CheckUserCertification(req.CertificationCID)
	if err != nil {
		fmt.Println("❌ Failed to check existing certification:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check certification CID"})
	}

	if !cidExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certification CID already exists in Blockchain"})
	}

	// ✅ แปลงวันที่จาก string → time.Time
	issuedDate, err := time.Parse("2006-01-02", req.IssuedDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid issued date format. Use YYYY-MM-DD"})
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid expiry date format. Use YYYY-MM-DD"})
	}

	// ✅ แปลง time.Time → *big.Int
	issuedDateBigInt := big.NewInt(issuedDate.Unix())
	expiryDateBigInt := big.NewInt(expiryDate.Unix())

	// ✅ สร้าง `eventID` ใหม่โดยใช้ `UUID` เพื่อป้องกันการซ้ำกัน
	eventID := fmt.Sprintf("EVENT-%s-%s", req.EntityID, uuid.New().String())

	// ✅ บันทึกใบเซอร์ใหม่ลง Blockchain
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

	fmt.Println("✅ Certification Event stored on Blockchain:", txHash)
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

	// ✅ ตรวจสอบว่า `entityID` มีอยู่ในระบบหรือไม่
	var farmer models.Farmer
	if err := database.DB.Where("entityid = ?", entityID).First(&farmer).Error; err != nil {
		fmt.Println("❌ [GetCertificationByEntity] Entity ID not found:", entityID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Entity ID not found"})
	}

	// ✅ ดึงข้อมูลใบเซอร์ทั้งหมดของ entityID จาก Blockchain
	certifications, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(entityID)
	if err != nil {
		fmt.Println("❌ [GetCertificationByEntity] Failed to fetch certifications from Blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch certifications"})
	}

	// ✅ ตรวจสอบ Query Parameter `?includeExpired=true`
	includeExpired := c.Query("includeExpired") == "true"

	// ✅ กรองเฉพาะใบเซอร์ที่ยัง Active หรือถ้า `includeExpired=true` ให้แสดงทั้งหมด
	var filteredCerts []certification.CertificationEventCertEvent
	for _, cert := range certifications {
		if includeExpired || cert.IsActive {
			filteredCerts = append(filteredCerts, cert)
		}
	}

	// ✅ เรียงลำดับ Certifications ตาม `expiryDate` (มาก → น้อย)
	sort.Slice(filteredCerts, func(i, j int) bool {
		return filteredCerts[i].ExpiryDate.Cmp(filteredCerts[j].ExpiryDate) > 0
	})

	// ✅ ส่งข้อมูลกลับ
	return c.JSON(fiber.Map{
		"entity_id":      entityID,
		"certifications": filteredCerts,
	})
}

// ✅ API: ปิดใช้งานใบเซอร์ที่ถูกเลือกโดย `eventID`
func DeleteCertification(c *fiber.Ctx) error {
	entityID := c.Params("entityID")
	eventID := c.Query("eventID") // ✅ รับ eventID จาก Query Parameter

	if entityID == "" || eventID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing entityID or eventID"})
	}

	fmt.Println("📌 [DeleteCertification] Deleting certification for entityID:", entityID, "eventID:", eventID)

	// ✅ ตรวจสอบว่า `entityID` มีอยู่จริงในฐานข้อมูลหรือไม่
	var farmer models.Farmer
	if err := database.DB.Where("entityid = ?", entityID).First(&farmer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Entity ID not found in database"})
	}

	// ✅ ดึงข้อมูลใบเซอร์ทั้งหมดของ `entityID` จาก Blockchain
	certifications, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(entityID)
	if err != nil {
		fmt.Println("❌ [DeleteCertification] Failed to fetch certifications from Blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch certifications"})
	}

	// ✅ ตรวจสอบว่า `eventID` มีอยู่จริงใน Blockchain หรือไม่
	var certExists bool
	for _, cert := range certifications {
		if cert.EventID == eventID && cert.IsActive {
			certExists = true
			break
		}
	}

	if !certExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Certification event not found or already inactive"})
	}

	// ✅ ปิดใช้งานใบเซอร์ที่ถูกเลือก
	txHash, err := services.BlockchainServiceInstance.DeactivateCertificationOnBlockchain(eventID)
	if err != nil {
		fmt.Println("❌ [DeleteCertification] Failed to deactivate certification:", eventID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deactivate certification"})
	}

	fmt.Println("✅ Certification Event deactivated on Blockchain:", txHash)
	return c.JSON(fiber.Map{
		"message":       "Certification deactivated successfully",
		"event_id":      eventID,
		"entity_id":     entityID,
		"blockchain_tx": txHash,
	})
}


func CheckCertificationCID(c *fiber.Ctx) error {
	certCID := c.Params("certCID")
	if certCID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing certCID"})
	}

	fmt.Println("🔍 [CheckCertificationCID] Checking CID in Blockchain:", certCID)

	// ✅ ใช้ `CheckUserCertification` แทน
	exists, err := services.BlockchainServiceInstance.CheckUserCertification(certCID)
	if err != nil {
		fmt.Println("❌ [CheckCertificationCID] Failed to check CID in Blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check certification CID"})
	}

	// ✅ `CheckUserCertification` คืนค่า `true` ถ้า CID **ยังไม่ถูกใช้** → ต้องกลับค่าก่อนส่ง
	return c.JSON(fiber.Map{
		"certCID": certCID,
		"exists":  !exists, // ✅ กลับค่าก่อนส่ง (true = มีอยู่แล้ว, false = ไม่มี)
	})
}
