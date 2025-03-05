package controllers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"finalyearproject/Backend/services"

	"github.com/gofiber/fiber/v2"
)

type RawMilkController struct {
	BlockchainService *services.BlockchainService
	QRCodeService     *services.QRCodeService
	IPFSService       *services.IPFSService
	MilkTankCounter   map[string]int
	Mutex             sync.Mutex
}

// ✅ ฟังก์ชันสร้าง Tank ID (FarmID + วันที่ + Running Number)
func (rmc *RawMilkController) generateTankID(farmID string) string {
	rmc.Mutex.Lock()
	defer rmc.Mutex.Unlock()

	// ✅ ดึงวันที่ปัจจุบันในรูปแบบ YYYYMMDD
	currentDate := time.Now().Format("20060102")

	// ✅ คีย์สำหรับเก็บ Running Number (FarmID + วันที่)
	key := farmID + "_" + currentDate

	// ✅ ถ้าไม่มีข้อมูลเก่า หรือเป็นวันใหม่ ให้รีเซ็ตเลขลำดับ
	if _, exists := rmc.MilkTankCounter[key]; !exists {
		rmc.MilkTankCounter[key] = 1
	} else {
		rmc.MilkTankCounter[key]++
	}

	// ✅ สร้าง Tank ID => FarmID + วันที่ + Running Number (3 หลัก)
	tankID := fmt.Sprintf("%s-%s-%03d", farmID, currentDate, rmc.MilkTankCounter[key])

	fmt.Println("✅ Generated Tank ID:", tankID)
	return tankID
}

// ✅ ฟังก์ชันสร้างแท็งก์นมดิบใหม่
func (rmc *RawMilkController) CreateMilkTank(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Create Milk Tank")

	// ✅ ดึงข้อมูลจาก JWT Token ที่อยู่ใน Cookie
	role := c.Locals("role").(string)
	farmID := c.Locals("entityID").(string)             // ✅ ฟาร์มไอดี
	walletAddress := c.Locals("walletAddress").(string) // ✅ ที่อยู่กระเป๋าเงินของเกษตรกร

	// ✅ ตรวจสอบสิทธิ์
	if role != "farmer" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only farmers can create milk tanks"})
	}

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		FarmName        string `json:"farmName"`
		PersonInCharge  string `json:"personInCharge"`
		Quantity        uint64 `json:"quantity"`
		QuantityUnit    string `json:"quantityUnit"`
		Temp            uint64 `json:"temp"`
		TempUnit        string `json:"tempUnit"`
		PH              uint64 `json:"pH"`
		Fat             uint64 `json:"fat"`
		Protein         uint64 `json:"protein"`
		Bacteria        bool   `json:"bacteria"`
		BacteriaInfo    string `json:"bacteriaInfo"`
		Contaminants    bool   `json:"contaminants"`
		ContaminantInfo string `json:"contaminantInfo"`
		AbnormalChar    bool   `json:"abnormalChar"`
		AbnormalType    struct {
			SmellBad      bool `json:"smellBad"`
			SmellNotFresh bool `json:"smellNotFresh"`
			AbnormalColor bool `json:"abnormalColor"`
			Sour          bool `json:"sour"`
			Bitter        bool `json:"bitter"`
			Cloudy        bool `json:"cloudy"`
			Lumpy         bool `json:"lumpy"`
			Separation    bool `json:"separation"`
		} `json:"abnormalType"`
		ShippingAddress struct {
			CompanyName string `json:"companyName"`
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			Email       string `json:"email"`
			AreaCode    string `json:"areaCode"`
			PhoneNumber string `json:"phoneNumber"`
			Address     string `json:"address"`
			Province    string `json:"province"`
			District    string `json:"district"`
			SubDistrict string `json:"subDistrict"`
			PostalCode  string `json:"postalCode"`
			Location    string `json:"location"`
		} `json:"shippingAddress"`
	}

	if err := c.BodyParser(&request); err != nil {
		fmt.Println("❌ Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ สร้าง `tankId` ตามรูปแบบที่กำหนด (ฟาร์มไอดี + วันเดือนปี + หมายเลขแท็งก์ที่เพิ่มขึ้น)
	tankId := rmc.generateTankID(farmID)

	fmt.Println("BlockchainService instance:", rmc.BlockchainService)

	valid, validationMsg := rmc.BlockchainService.ValidateMilkData(
		request.Quantity,
		request.Temp*100,    // ✅ คูณ 100 ก่อนส่ง
		request.PH*100,      // ✅ คูณ 100 ก่อนส่ง
		request.Fat*100,     // ✅ คูณ 100 ก่อนส่ง
		request.Protein*100, // ✅ คูณ 100 ก่อนส่ง
		request.Bacteria,
		request.Contaminants,
	)
	if !valid {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationMsg})
	}

	// ✅ รวมข้อมูล `bacteriaInfo`, `contaminantInfo`, `abnormalType`, และ `shippingAddress` ลงไฟล์เดียวแล้วอัปโหลดไป IPFS
	rawMilkData := map[string]interface{}{
		"bacteriaInfo":    request.BacteriaInfo,
		"contaminantInfo": request.ContaminantInfo,
		"abnormalType":    request.AbnormalType,
		"shippingAddress": request.ShippingAddress,
	}
	// ✅ แปลง ShippingAddress struct เป็น map[string]interface{}
	shippingAddressMap := map[string]interface{}{
		"companyName": request.ShippingAddress.CompanyName,
		"firstName":   request.ShippingAddress.FirstName,
		"lastName":    request.ShippingAddress.LastName,
		"email":       request.ShippingAddress.Email,
		"areaCode":    request.ShippingAddress.AreaCode,
		"phoneNumber": request.ShippingAddress.PhoneNumber,
		"address":     request.ShippingAddress.Address,
		"province":    request.ShippingAddress.Province,
		"district":    request.ShippingAddress.District,
		"subDistrict": request.ShippingAddress.SubDistrict,
		"postalCode":  request.ShippingAddress.PostalCode,
		"location":    request.ShippingAddress.Location,
	}

	qualityReportCID, err := rmc.IPFSService.UploadMilkDataToIPFS(rawMilkData, shippingAddressMap)
	if err != nil {
		fmt.Println("❌ Failed to upload quality report to IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload quality report"})
	}

	// ✅ สร้าง QR Code สำหรับแท็งก์นม
	qrCodeCID, err := rmc.QRCodeService.GenerateQRCode(tankId)
	if err != nil {
		fmt.Println("❌ Failed to generate QR Code:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate QR Code"})
	}

	// ✅ บันทึกแท็งก์นมดิบบน Blockchain
	txHash, err := rmc.BlockchainService.CreateMilkTank(walletAddress, tankId, request.PersonInCharge, qrCodeCID)
	if err != nil {
		fmt.Println("❌ Blockchain transaction failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
	}

	// ✅ ส่ง Response กลับ
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":          "Milk tank created successfully",
		"tankId":           tankId,
		"txHash":           txHash,
		"qrCodeCID":        qrCodeCID,
		"qualityReportCID": qualityReportCID,
	})
}
