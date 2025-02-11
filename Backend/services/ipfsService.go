package services

import (
	"bytes"
	"fmt"
	"io"

	shell "github.com/ipfs/go-ipfs-api"
)

// IPFSService ใช้สำหรับอัปโหลดไฟล์ไปยัง IPFS
type IPFSService struct {
	shell *shell.Shell
}

// NewIPFSService คืนค่าอินสแตนซ์ของ IPFSService
func NewIPFSService() *IPFSService {
	return &IPFSService{
		shell: shell.NewShell("localhost:5001"), // ✅ เชื่อมต่อ IPFS Daemon
	}
}

// UploadFile อัปโหลดไฟล์ไปยัง IPFS และคืนค่า CID
func (s *IPFSService) UploadFile(file io.Reader) (string, error) {
	fmt.Println("📌 Uploading file to IPFS...")

	// ✅ ตรวจสอบว่า IPFS Daemon ทำงานอยู่หรือไม่
	if !s.shell.IsUp() {
		fmt.Println("❌ IPFS Daemon is not running!")
		return "", fmt.Errorf("IPFS node is not available")
	}

	// ✅ อ่านไฟล์เป็น bytes
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		fmt.Println("❌ Error copying file content:", err)
		return "", fmt.Errorf("failed to read file content")
	}

	// ✅ อัปโหลดไปยัง IPFS
	cid, err := s.shell.Add(buf)
	if err != nil {
		fmt.Println("❌ Failed to upload to IPFS:", err)
		return "", fmt.Errorf("failed to upload to IPFS")
	}

	fmt.Println("✅ File uploaded to IPFS with CID:", cid)
	return cid, nil
}
