package services

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	shell "github.com/ipfs/go-ipfs-api"
)

// IPFSService ใช้สำหรับอัปโหลดและดึงไฟล์จาก IPFS
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

// GetFile ดึงไฟล์จาก IPFS และคืนค่าเป็น JSON String
func (s *IPFSService) GetFile(cid string) (string, error) {
	fmt.Println("📌 Retrieving file from IPFS... CID:", cid)

	// ✅ ตรวจสอบว่า IPFS Daemon ทำงานอยู่หรือไม่
	if !s.shell.IsUp() {
		fmt.Println("❌ IPFS Daemon is not running!")
		return "", fmt.Errorf("IPFS node is not available")
	}

	// ✅ ดึงไฟล์จาก IPFS
	reader, err := s.shell.Cat(cid)
	if err != nil {
		fmt.Println("❌ Failed to retrieve file from IPFS:", err)
		return "", fmt.Errorf("failed to retrieve file from IPFS")
	}
	defer reader.Close()

	// ✅ อ่านข้อมูลจาก reader
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("❌ Error reading file content:", err)
		return "", fmt.Errorf("failed to read file content")
	}

	fmt.Println("✅ File retrieved from IPFS successfully")
	return string(data), nil
}
