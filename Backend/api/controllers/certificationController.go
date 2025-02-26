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

	// ✅ รับไฟล์จาก `multipart/form-data`
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

	// ✅ อัปโหลดไปยัง IPFS
	cid, err := ipfsService.UploadFile(src) // ✅ เปลี่ยนจาก UploadBase64File → UploadFile
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

	// ตรวจสอบค่า CID
	if req.CertificationCID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certification CID is required"})
	}

	// แปลงวันที่จาก string → time.Time
	issuedDate, err := time.Parse("2006-01-02", req.IssuedDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid issued date format. Use YYYY-MM-DD"})
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid expiry date format. Use YYYY-MM-DD"})
	}

	// แปลง time.Time → *big.Int
	issuedDateBigInt := big.NewInt(issuedDate.Unix())
	expiryDateBigInt := big.NewInt(expiryDate.Unix())

	// ✅ ดึงใบเซอร์ทั้งหมดของ Entity จาก Blockchain
	existingCerts, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(req.EntityID)
	if err != nil {
		fmt.Println("❌ Failed to fetch existing certifications:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch existing certifications"})
	}

	// ✅ ค้นหาใบเซอร์ที่ยัง Active และหมดอายุไกลที่สุด
	var latestActiveCert *certification.CertificationEventCertEvent
	for _, cert := range existingCerts {
		if cert.IsActive {
			if latestActiveCert == nil || cert.ExpiryDate.Cmp(latestActiveCert.ExpiryDate) > 0 {
				latestActiveCert = &cert
			}
		}
	}

	// ✅ ถ้ามีใบเซอร์ที่ยัง Active และยังไม่หมดอายุ → อัปเดตแทนที่จะสร้างใหม่
	if latestActiveCert != nil {
		expiryDateTime := time.Unix(latestActiveCert.ExpiryDate.Int64(), 0)
		if expiryDateTime.After(time.Now()) {
			fmt.Println("📌 Updating existing active certification on Blockchain...")
			txHash, err := services.BlockchainServiceInstance.StoreCertificationOnBlockchain(
				latestActiveCert.EventID,
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
				"event_id":      latestActiveCert.EventID,
				"cid":           req.CertificationCID,
				"blockchain_tx": txHash,
			})
		}
	}

	// ✅ ถ้าไม่มีใบเซอร์ที่ยัง Active → สร้างใหม่
	eventID := fmt.Sprintf("EVENT-%s", req.EntityID)
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

	// ✅ ดึงข้อมูลใบเซอร์ทั้งหมดของ entityID จาก Blockchain
	certifications, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(entityID)
	if err != nil {
		fmt.Println("❌ [GetCertificationByEntity] Failed to fetch certifications from Blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch certifications"})
	}

	// ✅ กรองเฉพาะใบเซอร์ที่ยัง Active
	var activeCerts []certification.CertificationEventCertEvent
	for _, cert := range certifications {
		if cert.IsActive {
			activeCerts = append(activeCerts, cert)
		}
	}

	// ✅ เรียงลำดับ Active Certifications ตาม `expiryDate` (มาก → น้อย)
	sort.Slice(activeCerts, func(i, j int) bool {
		return activeCerts[i].ExpiryDate.Cmp(activeCerts[j].ExpiryDate) > 0
	})

	// ✅ ส่งข้อมูลกลับ (เฉพาะใบเซอร์ที่ยัง Active)
	return c.JSON(fiber.Map{
		"entity_id":      entityID,
		"certifications": activeCerts,
	})
}

func DeleteCertification(c *fiber.Ctx) error {
	entityID := c.Params("entityID")
	if entityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing entity ID"})
	}

	// ✅ ดึงข้อมูลใบเซอร์ทั้งหมดของ entityID จาก Blockchain
	certifications, err := services.BlockchainServiceInstance.GetAllCertificationsForEntity(entityID)
	if err != nil {
		fmt.Println("❌ [DeleteCertification] Failed to fetch certifications from Blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch certifications"})
	}

	// ✅ กรองเฉพาะใบเซอร์ที่ยัง Active
	var activeCerts []certification.CertificationEventCertEvent
	for _, cert := range certifications {
		if cert.IsActive {
			activeCerts = append(activeCerts, cert)
		}
	}

	// ✅ ถ้าไม่มีใบเซอร์ที่ยัง Active → แจ้งว่าไม่มีอะไรให้ลบ
	if len(activeCerts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No active certifications found for this entity"})
	}

	// ✅ ลบทุกใบเซอร์ที่ยัง Active
	var txHashes []string
	for _, cert := range activeCerts {
		eventID := cert.EventID
		txHash, err := services.BlockchainServiceInstance.DeactivateCertificationOnBlockchain(eventID)
		if err != nil {
			fmt.Println("❌ [DeleteCertification] Failed to deactivate certification:", eventID, err)
			continue // ข้ามไปใบถัดไป
		}
		txHashes = append(txHashes, txHash)
		fmt.Println("✅ Certification deactivated on Blockchain:", txHash)
	}

	// ✅ ส่งคืน `txHashes` ของทุกใบเซอร์ที่ถูกลบ
	return c.JSON(fiber.Map{
		"message":       "Certifications deactivated successfully",
		"entity_id":     entityID,
		"blockchain_tx": txHashes,
	})
}
