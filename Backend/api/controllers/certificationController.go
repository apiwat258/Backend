package controllers

import (
	"finalyearproject/Backend/services"
	certification "finalyearproject/Backend/services/certification_event"
	"fmt"
	"math/big"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		allowedExtensions := []string{".pdf", ".jpg", ".jpeg", ".png"}
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
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GetCertificationByUser(c *fiber.Ctx) error {
	// ✅ ดึง `entityID` จาก JWT Token ที่ AuthMiddleware กำหนดไว้
	entityID, ok := c.Locals("entityID").(string)
	if !ok || entityID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - entityID missing"})
	}

	fmt.Println("🔍 [GetCertificationByUser] Fetching certification for entityID:", entityID)

	// ✅ ดึงข้อมูลใบเซอร์ทั้งหมดของ entityID จาก Blockchain
	certifications, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(entityID)
	if err != nil {
		fmt.Println("❌ [GetCertificationByUser] Failed to fetch certifications from Blockchain:", err)
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

// DeleteCertification - ปิดใช้งานใบเซอร์จาก Blockchain
func DeleteCertification(c *fiber.Ctx) error {
	// ✅ ดึง `walletAddress` และ `entityID` จาก JWT Token
	walletAddress, ok := c.Locals("walletAddress").(string)
	if !ok || walletAddress == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - walletAddress missing"})
	}

	entityID, ok := c.Locals("entityID").(string)
	if !ok || entityID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - entityID missing"})
	}

	fmt.Println("📌 [DeleteCertification] Deleting certification for Wallet:", walletAddress, "EntityID:", entityID)

	// ✅ รับ eventID จาก Query Parameter
	eventID := c.Query("eventID")
	if eventID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing eventID"})
	}

	fmt.Println("📌 [DeleteCertification] Using entityID:", entityID, "and eventID:", eventID)

	// ✅ ดึงข้อมูลใบเซอร์ทั้งหมดของ `entityID` จาก Blockchain
	certifications, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(entityID)
	if err != nil {
		fmt.Println("❌ [DeleteCertification] Failed to fetch certifications from Blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch certifications"})
	}

	fmt.Println("📌 Retrieved", len(certifications), "certifications from blockchain for entity:", entityID)

	// ✅ ตรวจสอบว่า `eventID` มีอยู่จริงหรือไม่
	var certExists bool
	for _, cert := range certifications {
		fmt.Println("🔍 Checking eventID:", cert.EventID, "→ Active:", cert.IsActive) // ✅ Debug
		if cert.EventID == eventID && cert.IsActive {
			certExists = true
			break
		}
	}

	if !certExists {
		fmt.Println("❌ Certification event not found or already inactive:", eventID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Certification event not found or already inactive"})
	}

	// ✅ ปิดใช้งานใบเซอร์ใน Blockchain โดยใช้ `walletAddress`
	txHash, err := services.BlockchainServiceInstance.DeactivateCertificationOnBlockchain(walletAddress, eventID)
	if err != nil {
		fmt.Println("❌ [DeleteCertification] Failed to deactivate certification:", eventID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deactivate certification"})
	}

	fmt.Println("✅ Certification Event deactivated on Blockchain:", txHash)
	return c.JSON(fiber.Map{
		"message":       "Certification deactivated successfully",
		"event_id":      eventID,
		"entity_id":     entityID,
		"wallet":        walletAddress,
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

// StoreCertification - บันทึกใบรับรองลง Blockchain
func StoreCertification(c *fiber.Ctx) error {
	// ✅ ดึง `walletAddress` และ `entityID` จาก JWT Token
	walletAddress, ok := c.Locals("walletAddress").(string)
	if !ok || walletAddress == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - walletAddress missing"})
	}

	entityID, ok := c.Locals("entityID").(string)
	if !ok || entityID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - entityID missing"})
	}

	role, ok := c.Locals("role").(string)
	if !ok || role == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - role missing"})
	}

	// ✅ แปลง Role เป็น EntityType
	entityType := map[string]string{
		"farmer":    "Farmer",
		"factory":   "Factory",
		"logistics": "Logistics",
		"retailer":  "Retailer",
	}[role]

	if entityType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role type"})
	}

	// ✅ รับค่า `certCID` จาก Body
	var request struct {
		CertCID string `json:"certCID"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if request.CertCID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing certCID"})
	}

	// ✅ สร้าง `eventID`
	eventID := fmt.Sprintf("EVENT-%s-%s", entityID, uuid.New().String())

	// ✅ กำหนด `issuedDate` และ `expiryDate`
	issuedDate := big.NewInt(time.Now().Unix())
	expiryDate := big.NewInt(time.Now().AddDate(1, 0, 0).Unix()) // หมดอายุในอีก 1 ปี

	fmt.Println("📌 [StoreCertification] Storing Certification on Blockchain...")
	fmt.Println("   - Wallet Address:", walletAddress)
	fmt.Println("   - Entity Type:", entityType)
	fmt.Println("   - Entity ID:", entityID)
	fmt.Println("   - Cert CID:", request.CertCID)
	fmt.Println("   - Issued Date:", issuedDate)
	fmt.Println("   - Expiry Date:", expiryDate)

	// ✅ ส่งธุรกรรมไปยัง Blockchain
	txHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(
		walletAddress, eventID, entityType, entityID, request.CertCID, issuedDate, expiryDate,
	)
	if err != nil {
		fmt.Println("❌ [StoreCertification] Failed to store certification:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store certification"})
	}

	fmt.Println("✅ [StoreCertification] Certification stored successfully. TX Hash:", txHash)

	return c.JSON(fiber.Map{
		"message":       "Certification stored successfully",
		"event_id":      eventID,
		"entity_id":     entityID,
		"wallet":        walletAddress,
		"blockchain_tx": txHash,
	})
}
