package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

var IPFSServiceInstance = &IPFSService{}

// ✅ IPFSService ใช้สำหรับอัปโหลดและดึงไฟล์จาก IPFS
type IPFSService struct {
	shell *shell.Shell // ✅ คงเป็นตัวพิมพ์เล็ก
}

// ✅ ฟังก์ชันสร้างอินสแตนซ์ของ IPFSService
func NewIPFSService() *IPFSService {
	ipfsURL := os.Getenv("IPFS_API_URL")
	if ipfsURL == "" {
		ipfsURL = "http://localhost:5001" // ✅ ใช้ค่า default ถ้าไม่มี .env
	}

	return &IPFSService{
		shell: shell.NewShell(ipfsURL),
	}
}

func InitIPFSService() {
	ipfsURL := os.Getenv("IPFS_API_URL")
	if ipfsURL == "" {
		ipfsURL = "http://localhost:5001"
	}

	shellInstance := shell.NewShell(ipfsURL)
	if shellInstance == nil {
		fmt.Println("❌ Failed to initialize IPFS Shell")
		return
	}

	IPFSServiceInstance = &IPFSService{
		shell: shellInstance,
	}

	// ✅ ตรวจสอบการเชื่อมต่อ IPFS
	_, _, err := IPFSServiceInstance.shell.Version()
	if err != nil {
		fmt.Println("❌ Failed to connect to IPFS:", err)
		IPFSServiceInstance = nil // ✅ ป้องกันการใช้ instance ที่เชื่อมต่อไม่ได้
	} else {
		fmt.Println("✅ Connected to IPFS at:", ipfsURL)
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

// ✅ ฟังก์ชันดึงข้อมูล JSON จาก IPFS
func (s *IPFSService) GetFromIPFS(cid string) ([]byte, error) {
	fmt.Println("📌 Retrieving file from IPFS... CID:", cid)

	// ✅ ตรวจสอบว่า IPFS Daemon ทำงานอยู่หรือไม่
	if !s.shell.IsUp() {
		fmt.Println("❌ IPFS Daemon is not running!")
		return nil, fmt.Errorf("IPFS node is not available")
	}

	// ✅ ดึงไฟล์จาก IPFS
	reader, err := s.shell.Cat(cid)
	if err != nil {
		fmt.Println("❌ Failed to retrieve file from IPFS:", err)
		return nil, fmt.Errorf("failed to retrieve file from IPFS")
	}
	defer reader.Close()

	// ✅ อ่านข้อมูลจาก reader
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("❌ Error reading file content:", err)
		return nil, fmt.Errorf("failed to read file content")
	}

	fmt.Println("✅ File retrieved from IPFS successfully")
	return data, nil // ✅ คืนค่าเป็น `[]byte`
}

// ✅ ฟังก์ชันดึงข้อมูล JSON และแปลงเป็น Map
func (s *IPFSService) GetJSONFromIPFS(cid string) (map[string]interface{}, error) {
	data, err := s.GetFromIPFS(cid)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Println("❌ Error parsing JSON:", err)
		return nil, fmt.Errorf("invalid JSON format from IPFS")
	}

	return jsonData, nil
}

// ✅ ฟังก์ชันดึงไฟล์จาก IPFS
func (s *IPFSService) GetFileFromIPFS(cid string) ([]byte, error) {
	fmt.Println("📌 Retrieving file from IPFS... CID:", cid)

	// ✅ ตรวจสอบว่า IPFS Daemon ทำงานอยู่หรือไม่
	if !s.shell.IsUp() {
		fmt.Println("❌ IPFS Daemon is not running!")
		return nil, fmt.Errorf("IPFS node is not available")
	}

	// ✅ ดึงไฟล์จาก IPFS
	reader, err := s.shell.Cat(cid)
	if err != nil {
		fmt.Println("❌ Failed to retrieve file from IPFS:", err)
		return nil, fmt.Errorf("failed to retrieve file from IPFS")
	}
	defer reader.Close()

	// ✅ อ่านข้อมูลจาก reader
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("❌ Error reading file content:", err)
		return nil, fmt.Errorf("failed to read file content")
	}

	fmt.Println("✅ File retrieved from IPFS successfully")
	return data, nil
}

// ✅ ฟังก์ชันดึงรูปภาพจาก IPFS และแปลงเป็น Base64
func (s *IPFSService) GetImageBase64FromIPFS(cid string) (string, error) {
	data, err := s.GetFileFromIPFS(cid)
	if err != nil {
		return "", err
	}

	// ✅ แปลงไฟล์เป็น Base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	return base64Data, nil
}

// ✅ ฟังก์ชันอัปโหลดไฟล์จาก Base64 ไปยัง IPFS
func (s *IPFSService) UploadBase64File(base64Str string) (string, error) {
	fmt.Println("📌 Uploading Base64 file to IPFS...")

	// ✅ ตรวจสอบว่า Base64 String ไม่ว่าง
	if base64Str == "" {
		fmt.Println("❌ Base64 string is empty")
		return "", fmt.Errorf("empty base64 string")
	}

	// ✅ ตัด `data:image/png;base64,` หรือ `data:application/pdf;base64,` ออกถ้ามี
	if strings.Contains(base64Str, ",") {
		parts := strings.SplitN(base64Str, ",", 2)
		base64Str = parts[1] // ใช้เฉพาะส่วนข้อมูล
	}

	// ✅ เติม Padding (`=`) ให้ครบถ้าขาดไป
	padding := len(base64Str) % 4
	if padding > 0 {
		base64Str += strings.Repeat("=", 4-padding)
	}

	// ✅ แปลง Base64 เป็นไบต์ (ใช้ `RawStdEncoding` รองรับข้อมูลแบบไม่มี padding)
	data, err := base64.RawStdEncoding.DecodeString(base64Str)
	if err != nil {
		fmt.Println("❌ Failed to decode Base64:", err)
		return "", fmt.Errorf("failed to decode Base64")
	}

	// ✅ แปลงเป็น Buffer
	buf := bytes.NewReader(data)

	// ✅ ตรวจสอบว่า IPFS Daemon ทำงานอยู่หรือไม่
	if !s.shell.IsUp() {
		fmt.Println("❌ IPFS Daemon is not running!")
		return "", fmt.Errorf("IPFS node is not available")
	}

	// ✅ อัปโหลดไฟล์ไปยัง IPFS
	cid, err := s.shell.Add(buf)
	if err != nil {
		fmt.Println("❌ Failed to upload to IPFS:", err)
		return "", fmt.Errorf("failed to upload to IPFS")
	}

	fmt.Println("✅ File uploaded to IPFS with CID:", cid)
	return cid, nil
}

// ✅ ฟังก์ชันอัปโหลดข้อมูลนมดิบ + Shipping Address ขึ้น IPFS
func (s *IPFSService) UploadMilkDataToIPFS(rawMilkData map[string]interface{}, shippingAddress map[string]interface{}) (string, error) {
	fmt.Println("📌 Uploading Milk Data to IPFS...")

	// ✅ สร้าง JSON ที่รวม Raw Milk Data และ Shipping Address
	data := map[string]interface{}{
		"rawMilkData":     rawMilkData,
		"shippingAddress": shippingAddress,
	}

	// ✅ แปลง JSON เป็น bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return "", fmt.Errorf("Failed to encode JSON: %v", err)
	}

	// ✅ แปลงเป็น Buffer
	buf := bytes.NewReader(jsonData)

	// ✅ ตรวจสอบว่า IPFS Daemon ทำงานอยู่หรือไม่
	if !s.shell.IsUp() {
		fmt.Println("❌ IPFS Daemon is not running!")
		return "", fmt.Errorf("IPFS node is not available")
	}

	// ✅ อัปโหลด JSON ไปยัง IPFS
	cid, err := s.shell.Add(buf)
	if err != nil {
		fmt.Println("❌ Failed to upload Milk Data to IPFS:", err)
		return "", fmt.Errorf("Failed to upload to IPFS: %v", err)
	}

	fmt.Println("✅ Milk Data uploaded to IPFS with CID:", cid)
	return cid, nil
}

// ✅ ฟังก์ชันอัปโหลดข้อมูล `Product` ไปยัง IPFS
func (s *IPFSService) UploadDataToIPFS(productData map[string]interface{}) (string, error) {
	fmt.Println("📌 Uploading Product Data to IPFS...")

	// ✅ แปลง JSON เป็น bytes
	jsonData, err := json.Marshal(productData)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return "", fmt.Errorf("Failed to encode JSON: %v", err)
	}

	// ✅ แปลงเป็น Buffer
	buf := bytes.NewReader(jsonData)

	// ✅ ตรวจสอบว่า IPFS Daemon ทำงานอยู่หรือไม่
	if !s.shell.IsUp() {
		fmt.Println("❌ IPFS Daemon is not running!")
		return "", fmt.Errorf("IPFS node is not available")
	}

	// ✅ อัปโหลด JSON ไปยัง IPFS
	cid, err := s.shell.Add(buf)
	if err != nil {
		fmt.Println("❌ Failed to upload Product Data to IPFS:", err)
		return "", fmt.Errorf("Failed to upload to IPFS: %v", err)
	}

	fmt.Println("✅ Product Data uploaded to IPFS with CID:", cid)
	return cid, nil
}
