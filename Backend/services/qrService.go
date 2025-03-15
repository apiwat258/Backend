package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // ✅ รองรับ JPEG
	_ "image/png"
	"os"

	gozxing "github.com/makiuchi-d/gozxing"
	gozxingqrcode "github.com/makiuchi-d/gozxing/qrcode" // 📌 ตั้ง alias เพื่อหลีกเลี่ยงชื่อซ้ำ
	goqrcode "github.com/skip2/go-qrcode"
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
	qrCode, err := goqrcode.Encode(tankID, goqrcode.Medium, 256)
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

// ✅ แก้ให้ใช้ JSON แทน `tankID`
func (qrs *QRCodeService) GenerateQRCodeforFactory(qrData string) (string, error) {
	fmt.Println("📌 Generating QR Code for:", qrData)

	// ✅ ตรวจสอบว่า IPFSService ถูกกำหนดค่าหรือไม่
	if qrs.IPFSService == nil {
		fmt.Println("❌ IPFS Service is not initialized")
		return "", fmt.Errorf("IPFS Service is not initialized")
	}

	// ✅ แปลง JSON เป็น Struct เพื่อดึง `trackingId`
	var qrDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(qrData), &qrDataMap); err != nil {
		fmt.Println("❌ Failed to parse QR data JSON:", err)
		return "", fmt.Errorf("Failed to parse QR data JSON: %v", err)
	}

	trackingID, ok := qrDataMap["trackingId"].(string)
	if !ok || trackingID == "" {
		fmt.Println("❌ Missing trackingId in QR data")
		return "", errors.New("missing trackingId in QR data")
	}

	// ✅ สร้าง QR Code เป็น PNG (ใช้ JSON Object แทน `tankID`)
	qrCode, err := goqrcode.Encode(qrData, goqrcode.Medium, 256)
	if err != nil {
		fmt.Println("❌ Failed to generate QR Code:", err)
		return "", fmt.Errorf("Failed to generate QR Code: %v", err)
	}

	// ✅ ใช้ `trackingId` เป็นชื่อไฟล์ QR Code
	tempFilePath := fmt.Sprintf("/tmp/qrcode_%s.png", trackingID)
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

	// ✅ อัปโหลดไฟล์ QR Code ไปยัง IPFS
	qrCodeCID, err := qrs.IPFSService.UploadFile(file)
	if err != nil {
		fmt.Println("❌ Failed to upload QR Code to IPFS:", err)
		return "", fmt.Errorf("Failed to upload QR Code to IPFS: %v", err)
	}

	fmt.Println("✅ QR Code uploaded to IPFS with CID:", qrCodeCID)
	return qrCodeCID, nil
}

// ✅ ฟังก์ชันอ่าน QR Code จาก Base64
func (qrs *QRCodeService) ReadQRCodeFromBase64(base64Image string) (string, error) {
	// ✅ Decode Base64 เป็น Byte
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to decode Base64 image: %v", err)
	}

	// ✅ แปลง Byte เป็น Image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("❌ Failed to decode image: %v", err)
	}

	// ✅ ใช้ gozxing อ่าน QR Code
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create binary bitmap: %v", err)
	}

	qrReader := gozxingqrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to read QR Code: %v", err)
	}

	return result.GetText(), nil
}

func (qrs *QRCodeService) ReadQRCodeFromCID(qrCodeCID string) (map[string]interface{}, error) {
	fmt.Println("📌 Fetching and Decoding QR Code from IPFS CID:", qrCodeCID)

	// ✅ ดึง QR Code จาก IPFS
	qrBase64, err := qrs.IPFSService.GetImageBase64FromIPFS(qrCodeCID)
	if err != nil {
		fmt.Println("❌ Failed to fetch QR Code from IPFS:", err)
		return nil, fmt.Errorf("Failed to retrieve QR Code")
	}

	// ✅ ถอดรหัส QR Code
	qrData, err := qrs.ReadQRCodeFromBase64(qrBase64)
	if err != nil {
		fmt.Println("❌ Failed to decode QR Code:", err)
		return nil, fmt.Errorf("Failed to decode QR Code")
	}

	var qrDataMap map[string]interface{}
	err = json.Unmarshal([]byte(qrData), &qrDataMap)
	if err != nil {
		fmt.Println("❌ Failed to parse QR Code JSON:", err)
		return nil, fmt.Errorf("Failed to parse QR Code JSON")
	}

	return qrDataMap, nil

}
