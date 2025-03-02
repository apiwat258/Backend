package controllers

import (
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
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
	// ✅ ดึง `userID` จาก JWT Token ที่ AuthMiddleware กำหนดไว้
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [GetCertificationByUser] Fetching certification for userID:", userID)

	// ✅ ค้นหา `entityID` ของผู้ใช้ที่ล็อกอินอยู่
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		fmt.Println("❌ [GetCertificationByUser] User ID not found:", userID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.EntityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User has no associated entity"})
	}

	entityID := user.EntityID
	fmt.Println("✅ [GetCertificationByUser] Found entityID:", entityID)

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

func DeleteCertification(c *fiber.Ctx) error {
	entityID := c.Params("entityID")
	eventID := c.Query("eventID") // ✅ รับ eventID จาก Query Parameter

	if entityID == "" || eventID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing entityID or eventID"})
	}

	fmt.Println("📌 [DeleteCertification] Deleting certification for entityID:", entityID, "eventID:", eventID)

	// ✅ ตรวจสอบว่า `entityID` มีอยู่จริงในฐานข้อมูลหรือไม่
	var farmer models.Farmer
	if err := database.DB.Select("farmerid").Where("entityid = ?", entityID).First(&farmer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Entity ID not found in database"})
	}

	// ✅ ดึงข้อมูลใบเซอร์ทั้งหมดของ `entityID` จาก Blockchain
	certifications, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(entityID)
	if err != nil {
		fmt.Println("❌ [DeleteCertification] Failed to fetch certifications from Blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch certifications"})
	}

	// ✅ ตรวจสอบว่า `eventID` มีอยู่จริงใน Blockchain หรือไม่
	certMap := make(map[string]bool)
	for _, cert := range certifications {
		if cert.IsActive {
			certMap[cert.EventID] = true
		}
	}
	if !certMap[eventID] {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Certification event not found or already inactive"})
	}

	// ✅ ปิดใช้งานใบเซอร์ที่ถูกเลือก
	txHash, err := services.BlockchainServiceInstance.DeactivateCertificationOnBlockchain(eventID)
	if err != nil {
		fmt.Println("❌ [DeleteCertification] Failed to deactivate certification:", eventID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deactivate certification"})
	}

	if txHash == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
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

func UpdateCertification(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string) // ✅ ดึง UserID จาก Middleware
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("🔍 [UpdateCertification] Updating certification for userID:", userID)

	// ✅ ตรวจสอบว่า User มีฟาร์มหรือไม่
	var user models.User
	if err := database.DB.Where("userid = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if user.EntityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User does not have a registered farm"})
	}

	// ✅ ดึงฟาร์มจาก EntityID ของ User
	var farmer models.Farmer
	if err := database.DB.Where("farmerid = ?", user.EntityID).First(&farmer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Farm not found"})
	}

	// ✅ รับค่า `cert_cid` ใหม่จาก FormData
	newCertCID := strings.TrimSpace(c.FormValue("cert_cid"))
	if newCertCID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certification CID is required"})
	}

	// ✅ ดึงใบเซอร์เก่าของฟาร์มจาก Blockchain
	oldCerts, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(user.EntityID)
	if err != nil {
		fmt.Println("❌ Failed to fetch existing certifications:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch existing certifications"})
	}

	// ✅ เช็คว่าใบเซอร์ใหม่ซ้ำกับของเก่าหรือไม่
	for _, cert := range oldCerts {
		if cert.CertificationCID == newCertCID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "New certification CID is the same as the existing one"})
		}
	}

	// ✅ สร้าง `eventID` ใหม่
	eventID := fmt.Sprintf("EVENT-%s-%s", user.EntityID, uuid.New().String())

	// ✅ วันที่ออกใบเซอร์ และวันหมดอายุ (1 ปี)
	issuedDate := time.Now()
	expiryDate := issuedDate.AddDate(1, 0, 0)

	// ✅ แปลงวันที่เป็น *big.Int
	issuedDateBigInt := big.NewInt(issuedDate.Unix())
	expiryDateBigInt := big.NewInt(expiryDate.Unix())

	// ✅ บันทึกใบเซอร์ใหม่ลง Blockchain (โดยไม่ปิดใบเก่า)
	certTxHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(
		farmer.WalletAddress, // ✅ ใช้ Wallet Address ของฟาร์ม
		eventID,
		"farmer",
		user.EntityID,
		newCertCID,
		issuedDateBigInt,
		expiryDateBigInt,
	)

	if err != nil {
		fmt.Println("❌ Failed to store new certification on blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store new certification on blockchain"})
	}

	fmt.Println("✅ Certification updated on Blockchain. Transaction Hash:", certTxHash)

	// ✅ ส่งข้อมูลกลับให้ Frontend
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Certification uploaded successfully",
		"event_id":      eventID,
		"cert_cid":      newCertCID,
		"blockchain_tx": certTxHash,
	})
}
