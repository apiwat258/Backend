package services

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

var QRCodeServiceInstance = &QRCodeService{}

// QRCodeService จัดการสร้าง QR Code และอัปโหลดไปยัง IPFS
type QRCodeService struct {
	IPFSService *IPFSService
}

// NewQRCodeService สร้างอินสแตนซ์ของ QRCodeService
func NewQRCodeService(ipfsService *IPFSService) *QRCodeService {
	return &QRCodeService{
		IPFSService: ipfsService,
	}
}

// GenerateQRCode สร้าง QR Code จาก Tank ID และอัปโหลดไปยัง IPFS
func (qrs *QRCodeService) GenerateQRCode(tankID string) (string, error) {
	fmt.Println("📌 Generating QR Code for Tank ID:", tankID)

	// ✅ ตรวจสอบว่า IPFSService ถูกกำหนดค่าหรือไม่
	if qrs.IPFSService == nil {
		fmt.Println("❌ IPFS Service is not initialized")
		return "", fmt.Errorf("IPFS Service is not initialized")
	}

	// ✅ สร้าง QR Code เป็น PNG
	qrCode, err := qrcode.Encode(tankID, qrcode.Medium, 256)
	if err != nil {
		fmt.Println("❌ Failed to generate QR Code:", err)
		return "", fmt.Errorf("Failed to generate QR Code: %v", err)
	}

	// ✅ สร้างไฟล์ชั่วคราว
	tempFilePath := fmt.Sprintf("/tmp/qrcode_%s.png", tankID)
	err = os.WriteFile(tempFilePath, qrCode, 0644)
	if err != nil {
		fmt.Println("❌ Failed to save QR Code to file:", err)
		return "", fmt.Errorf("Failed to save QR Code to file: %v", err)
	}
	defer os.Remove(tempFilePath)

	// ✅ เปิดไฟล์ที่บันทึกไว้
	file, err := os.Open(tempFilePath)
	if err != nil {
		fmt.Println("❌ Failed to open QR Code file:", err)
		return "", fmt.Errorf("Failed to open QR Code file: %v", err)
	}
	defer file.Close()

	// ✅ ใช้ IPFSServiceInstance ที่ถูกต้อง
	qrCodeCID, err := qrs.IPFSService.UploadFile(file)
	if err != nil {
		fmt.Println("❌ Failed to upload QR Code to IPFS:", err)
		return "", fmt.Errorf("Failed to upload QR Code to IPFS: %v", err)
	}

	fmt.Println("✅ QR Code uploaded to IPFS with CID:", qrCodeCID)
	return qrCodeCID, nil
}
