package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode - ฟังก์ชันสร้าง QR Code เป็น Base64
func GenerateQRCode(data string) (string, error) {
	// ✅ สร้าง QR Code ด้วย go-qrcode
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		fmt.Println("❌ Failed to create QR Code:", err)
		return "", err
	}

	// ✅ แปลง QR Code เป็น PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, qr.Image(256))
	if err != nil {
		fmt.Println("❌ Failed to encode QR Code:", err)
		return "", err
	}

	// ✅ แปลงเป็น Base64
	qrBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return qrBase64, nil
}
