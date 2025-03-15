package controllers

import (
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services" // 📌 แก้ชื่อ package ให้ตรงกับโครงสร้างของโปรเจค
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ TrackingController Struct (ใช้โครงสร้างเดียวกับ `ProductLotController`)
type TrackingController struct {
	DB                *gorm.DB
	BlockchainService *services.BlockchainService
	IPFSService       *services.IPFSService
	QRService         *services.QRCodeService
}

// ✅ NewTrackingController - ฟังก์ชันสร้าง Controller
func NewTrackingController(db *gorm.DB, blockchainService *services.BlockchainService, ipfsService *services.IPFSService, qrService *services.QRCodeService) *TrackingController {
	return &TrackingController{
		DB:                db,
		BlockchainService: blockchainService,
		IPFSService:       ipfsService,
		QRService:         qrService,
	}
}

// ✅ GetTrackingDetails - ฟังก์ชันกลางสำหรับดึงข้อมูล Tracking
func (tc *TrackingController) GetTrackingDetails(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Full Tracking Details")

	// ✅ รับ Tracking ID จาก Query Parameter
	trackingId := c.Query("trackingId")
	if trackingId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tracking ID is required"})
	}

	// ✅ ใช้ Goroutines + WaitGroup เพื่อลด Latency
	var wg sync.WaitGroup
	var response sync.Map
	errorList := []string{}

	// ✅ Helper Function สำหรับดึงข้อมูลแต่ละ Role
	fetchTrackingData := func(role string, fetchFunc func(string) (fiber.Map, error)) {
		defer wg.Done()
		data, err := fetchFunc(trackingId)
		if err != nil {
			fmt.Printf("❌ Failed to fetch %s data: %v\n", role, err)
			errorList = append(errorList, fmt.Sprintf("%s: %v", role, err))
		} else if data != nil {
			response.Store(role, data) // ✅ ป้องกัน nil และใช้ sync.Map
		}
	}

	// ✅ ใช้ Goroutines ดึงข้อมูลจากแต่ละ Role
	wg.Add(4)
	go fetchTrackingData("retailer", tc.GetRetailerTrackingData)
	go fetchTrackingData("logistics", tc.GetLogisticsTrackingData)
	go fetchTrackingData("factory", tc.GetFactoryTrackingData)
	go fetchTrackingData("farm", tc.GetFarmTrackingData)

	// ✅ รอให้ทุก Goroutines ทำงานเสร็จ
	wg.Wait()

	// ✅ รวมผลลัพธ์จาก `sync.Map`
	finalResponse := fiber.Map{}
	response.Range(func(key, value interface{}) bool {
		finalResponse[key.(string)] = value
		return true
	})

	// ✅ ตรวจสอบ Error List และเพิ่มไปใน Response
	if len(errorList) > 0 {
		fmt.Println("⚠️ Some data could not be fetched:", errorList)
		finalResponse["errors"] = errorList
	}

	// ✅ ส่ง Response กลับไปที่ Frontend
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}

// ✅ ฟังก์ชันย่อย (ตอนนี้ยังเป็นค่าว่าง)
func (tc *TrackingController) GetRetailerTrackingData(trackingId string) (fiber.Map, error) {
	fmt.Println("📌 Fetching Retailer Tracking Data for Tracking ID:", trackingId)

	// ✅ ดึงข้อมูล Retailer Confirmation จาก Blockchain
	retailerData, err := tc.BlockchainService.GetRetailerConfirmation(trackingId)
	if err != nil {
		fmt.Println("❌ Failed to fetch retailer confirmation:", err)
		return nil, fmt.Errorf("Failed to fetch retailer confirmation")
	}

	retailerID, ok := retailerData["retailerId"].(string)
	if !ok || retailerID == "" {
		return nil, fmt.Errorf("Invalid retailer ID")
	}

	// ✅ ดึงข้อมูลจาก Database (ตาราง `retailer`)
	var retailer models.Retailer
	err = tc.DB.Where("retailerid = ?", retailerID).First(&retailer).Error
	if err != nil {
		fmt.Println("❌ Failed to fetch retailer general info:", err)
		return nil, fmt.Errorf("Failed to fetch retailer general info")
	}

	// ✅ ดึงข้อมูลจาก IPFS โดยใช้ Quality CID
	qualityCID, ok := retailerData["qualityCID"].(string)
	if !ok || qualityCID == "" {
		return nil, fmt.Errorf("Invalid or missing qualityCID")
	}

	qualityData, err := tc.IPFSService.GetJSONFromIPFS(qualityCID)
	if err != nil {
		fmt.Println("❌ Failed to fetch quality report from IPFS:", err)
		return nil, fmt.Errorf("Failed to retrieve quality data from IPFS")
	}

	fmt.Println("📌 Debug - Quality Data from IPFS:", qualityData)

	// ✅ ตรวจสอบโครงสร้าง `recipientInfo`
	recipientInfo, _ := qualityData["recipientInfo"].(map[string]interface{})

	// ✅ ดึงข้อมูล `Quantity`
	quantity, _ := qualityData["quantity"].(float64)
	quantityUnit, _ := qualityData["quantityUnit"].(string)
	temp, _ := qualityData["temperature"].(float64)
	tempUnit, _ := qualityData["tempUnit"].(string)
	pH, _ := qualityData["pH"].(float64)
	fat, _ := qualityData["fat"].(float64)
	protein, _ := qualityData["protein"].(float64)
	bacteria, _ := qualityData["bacteria"].(bool)
	bacteriaInfo, _ := qualityData["bacteriaInfo"].(string)
	contaminants, _ := qualityData["contaminants"].(bool)
	contaminantInfo, _ := qualityData["contaminantInfo"].(string)
	abnormalChar, _ := qualityData["abnormalChar"].(bool)
	abnormalType, _ := qualityData["abnormalType"].(map[string]interface{})

	// ✅ ป้องกันค่าที่เป็น `nil`
	if abnormalType == nil {
		abnormalType = map[string]interface{}{}
	}

	// ✅ จัดรูปแบบข้อมูลให้ตรงกับ JSON ที่ต้องการ
	response := fiber.Map{
		"generalInfo": fiber.Map{ // ✅ ข้อมูลทั่วไปของ Retailer จาก Database
			"retailerName": retailer.CompanyName,
			"address":      retailer.Address,
			"contact":      retailer.Telephone,
		},
		"receivedProduct": fiber.Map{ // ✅ ข้อมูลการรับสินค้าของ Retailer จาก Blockchain & IPFS
			"trackingId": trackingId,
			"recipientInfo": fiber.Map{
				"personInCharge": recipientInfo["personInCharge"],
				"location":       recipientInfo["location"],
				"pickUpTime":     recipientInfo["pickUpTime"],
			},
			"quantity": fiber.Map{
				"value":           quantity,
				"unit":            quantityUnit,
				"temperature":     temp,
				"tempUnit":        tempUnit,
				"pH":              pH,
				"fat":             fat,
				"protein":         protein,
				"bacteria":        bacteria,
				"bacteriaInfo":    bacteriaInfo,
				"contaminants":    contaminants,
				"contaminantInfo": contaminantInfo,
				"abnormalChar":    abnormalChar,
				"abnormalType":    abnormalType,
			},
		},
	}

	// ✅ ส่ง JSON ออกไป
	return response, nil
}

func (tc *TrackingController) GetLogisticsTrackingData(trackingId string) (fiber.Map, error) {
	fmt.Println("📌 Fetching Logistics Tracking Data for Tracking ID:", trackingId)

	// ✅ ดึงข้อมูลจาก Smart Contract ผ่าน BlockchainService
	beforeCheckpoints, duringCheckpoints, afterCheckpoints, err := tc.BlockchainService.GetLogisticsCheckpointsByTrackingId(trackingId)
	if err != nil {
		fmt.Println("❌ Failed to fetch logistics checkpoints:", err)
		return nil, fmt.Errorf("Failed to fetch logistics checkpoints")
	}

	// ✅ ใช้ `map[string]bool` แยกข้อมูลซ้ำ **แต่ละสถานะ**
	uniqueCheckpoints := make(map[string]bool)

	// ✅ ฟังก์ชันช่วยลบข้อมูลซ้ำซ้อนใน **แต่ละสถานะ**
	filterUniqueCheckpoints := func(checkpoints []services.LogisticsCheckpoint, checkType string) []map[string]interface{} {
		var filteredCheckpoints []map[string]interface{}

		for _, cp := range checkpoints {
			// ✅ ใช้ `CheckType` เพื่อแยก Before / During / After
			key := fmt.Sprintf("%s-%d-%d-%d", cp.TrackingId, cp.CheckType, cp.PickupTime, cp.DeliveryTime)

			// ✅ ถ้า Key นี้มีอยู่แล้ว แปลว่า Duplicate → ข้ามไป
			if uniqueCheckpoints[key] {
				continue
			}
			uniqueCheckpoints[key] = true

			// ✅ ดึง Receiver Info จาก IPFS
			fmt.Println("📡 Fetching Receiver Info from IPFS CID:", cp.ReceiverCID)
			ipfsData, err := tc.IPFSService.GetJSONFromIPFS(cp.ReceiverCID)
			if err != nil {
				fmt.Println("⚠️ Warning: Failed to fetch receiver info from IPFS:", err)
				continue
			}

			// ✅ แปลงข้อมูล ReceiverInfo
			receiverInfo := map[string]interface{}{
				"companyName": getStringFromMap(ipfsData, "companyName"),
				"firstName":   getStringFromMap(ipfsData, "firstName"),
				"lastName":    getStringFromMap(ipfsData, "lastName"),
				"email":       getStringFromMap(ipfsData, "email"),
				"phone":       getStringFromMap(ipfsData, "phone"),
				"address":     getStringFromMap(ipfsData, "address"),
				"province":    getStringFromMap(ipfsData, "province"),
				"district":    getStringFromMap(ipfsData, "district"),
				"subDistrict": getStringFromMap(ipfsData, "subDistrict"),
				"postalCode":  getStringFromMap(ipfsData, "postalCode"),
				"location":    getStringFromMap(ipfsData, "location"),
			}

			// ✅ เพิ่มข้อมูลเข้า Array ใหม่
			filteredCheckpoints = append(filteredCheckpoints, map[string]interface{}{
				"trackingId":        cp.TrackingId,
				"logisticsProvider": cp.LogisticsProvider,
				"pickupTime":        cp.PickupTime,
				"deliveryTime":      cp.DeliveryTime,
				"quantity":          cp.Quantity,
				"temperature":       cp.Temperature,
				"personInCharge":    cp.PersonInCharge,
				"checkType":         cp.CheckType,
				"receiverCID":       cp.ReceiverCID,
				"receiverInfo":      receiverInfo, // ✅ เพิ่มข้อมูลจาก IPFS
			})
		}

		return filteredCheckpoints
	}

	// ✅ Response JSON
	response := fiber.Map{
		"trackingId":        trackingId,
		"beforeCheckpoints": filterUniqueCheckpoints(beforeCheckpoints, "before"),
		"duringCheckpoints": filterUniqueCheckpoints(duringCheckpoints, "during"),
		"afterCheckpoints":  filterUniqueCheckpoints(afterCheckpoints, "after"),
	}

	// ✅ ส่งข้อมูลกลับไปที่ `GetTrackingDetails`
	return response, nil
}

// GetFarmTrackingData - ดึงข้อมูล Tracking ของฟาร์ม

func (tc *TrackingController) GetFarmTrackingData(trackingId string) (fiber.Map, error) {
	fmt.Println("📌 Fetching Farm Tracking Data for trackingId:", trackingId)

	// ✅ 1️⃣ ดึง Product Lot ID จาก Tracking ID
	productLotId, err := tc.BlockchainService.GetProductLotByTrackingId(trackingId)
	if err != nil {
		fmt.Println("❌ Failed to fetch Product Lot ID:", err)
		return nil, err
	}
	fmt.Println("✅ Found Product Lot ID:", productLotId)

	// ✅ 2️⃣ ดึงรายละเอียด Product Lot จาก Blockchain
	productLotData, err := tc.BlockchainService.GetProductLotByLotID(productLotId)
	if err != nil {
		fmt.Println("❌ Failed to fetch Product Lot details:", err)
		return nil, err
	}

	// ✅ 3️⃣ ตรวจสอบว่า Product Lot มี Milk Tank IDs หรือไม่
	if len(productLotData.MilkTankIDs) == 0 {
		fmt.Println("⚠️ No Milk Tanks found for this Product Lot")
		return fiber.Map{"farms": []fiber.Map{}}, nil
	}

	// ✅ 4️⃣ ใช้ map เพื่อแยก Milk Tank ตาม `farmID`
	farmMilkTanks := make(map[string][]string)
	for _, tankID := range productLotData.MilkTankIDs {
		parts := strings.Split(tankID, "-")
		if len(parts) < 2 {
			fmt.Println("⚠️ Invalid Milk Tank ID format:", tankID)
			continue
		}
		farmID := parts[0] // ✅ ดึง `farmID` จาก Tank ID
		farmMilkTanks[farmID] = append(farmMilkTanks[farmID], tankID)
	}

	// ✅ 5️⃣ ดึงข้อมูลฟาร์มจาก Database
	var farms []fiber.Map
	for farmID, tankIDs := range farmMilkTanks {
		var farm models.Farmer
		err := tc.DB.Where("farmerid = ?", farmID).First(&farm).Error
		if err != nil {
			fmt.Println("❌ Failed to fetch farm details for:", farmID, "Error:", err)
			continue
		}

		// ✅ เพิ่มข้อมูลฟาร์มเข้าไปในอาร์เรย์
		farms = append(farms, fiber.Map{
			"farmID":       farm.FarmerID,
			"companyName":  farm.CompanyName,
			"address":      farm.Address,
			"subDistrict":  farm.SubDistrict,
			"district":     farm.District,
			"province":     farm.Province,
			"country":      farm.Country,
			"postCode":     farm.PostCode,
			"telephone":    farm.Telephone,
			"email":        farm.Email,
			"locationLink": farm.LocationLink,
			"milkTankIDs":  tankIDs, // ✅ แต่ละฟาร์มมีแท้งค์ของตัวเอง
		})
	}

	// ✅ 6️⃣ สร้างโครงสร้างข้อมูล Response
	response := fiber.Map{
		"farms": farms, // ✅ ส่งข้อมูลฟาร์มทั้งหมดพร้อม Tank IDs
	}

	fmt.Println("✅ Farm Tracking Data fetched successfully")
	return response, nil
}

func (tc *TrackingController) GetFactoryTrackingData(trackingId string) (fiber.Map, error) {
	fmt.Println("📌 Fetching Factory Tracking Data for Tracking ID:", trackingId)

	// ✅ ดึงข้อมูล Product Lot ID จาก Smart Contract
	lotID, err := tc.BlockchainService.GetProductLotByTrackingId(trackingId)
	if err != nil {
		fmt.Println("❌ Failed to fetch Product Lot ID:", err)
		return nil, err
	}

	// ✅ ดึงข้อมูล Product Lot จาก Blockchain
	productLotData, err := tc.BlockchainService.GetProductLotByLotID(lotID)
	if err != nil {
		fmt.Println("❌ Failed to fetch product lot from blockchain:", err)
		return nil, fmt.Errorf("Failed to fetch product lot details")
	}

	// ✅ ดึงข้อมูล Product จาก Smart Contract
	productID := productLotData.ProductID
	productData, err := tc.BlockchainService.GetProductDetails(productID)
	if err != nil {
		fmt.Println("❌ Failed to fetch product from blockchain:", err)
		return nil, fmt.Errorf("Failed to fetch product details")
	}

	// ✅ ดึงข้อมูล QR Code จาก Tracking Event
	_, _, qrCodeCIDs, err := tc.BlockchainService.GetTrackingByLotId(lotID)
	if err != nil || len(qrCodeCIDs) == 0 {
		fmt.Println("❌ Failed to fetch QR Code from blockchain:", err)
		return nil, fmt.Errorf("Failed to fetch QR Code")
	}

	// ✅ ดึงข้อมูล QR Code Data จาก IPFS
	qrCodeData, err := tc.QRService.ReadQRCodeFromCID(qrCodeCIDs[0]) // ใช้ตัวแรก
	if err != nil {
		fmt.Println("❌ Failed to decode QR Code from CID:", qrCodeCIDs[0])
		return nil, fmt.Errorf("Failed to decode QR Code")
	}
	fmt.Println("📌 QR Code Data:", qrCodeData)
	factoryMap, ok := qrCodeData["factory"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Factory data structure is incorrect")
		return nil, fmt.Errorf("Factory data structure is incorrect")
	}

	factoryID, ok := factoryMap["factoryId"].(string)
	if !ok || factoryID == "" {
		fmt.Println("❌ FactoryID missing in QR Code Data")
		return nil, fmt.Errorf("FactoryID is missing in QR Code Data")
	}

	// ✅ ค้นหาข้อมูลโรงงานจาก PostgreSQL
	var factory models.Factory
	if err := tc.DB.Where("factoryid = ?", factoryID).First(&factory).Error; err != nil {
		fmt.Println("❌ Failed to fetch factory details:", err)
		return nil, fmt.Errorf("Failed to fetch factory details")
	}

	// ✅ ดึงข้อมูล JSON จาก IPFS ของ Product
	productIPFSCID := productData["productCID"].(string)
	productIPFSData, err := tc.IPFSService.GetJSONFromIPFS(productIPFSCID)
	if err != nil {
		fmt.Println("❌ Failed to fetch product data from IPFS:", err)
		return nil, fmt.Errorf("Failed to fetch product data")
	}

	// ✅ ตรวจสอบว่า Nutrition มีค่าหรือไม่
	nutritionData, ok := productIPFSData["nutrition"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Error: Nutrition data is missing or incorrect")
		return nil, fmt.Errorf("Nutrition data structure is incorrect")
	}

	// ✅ ดึงข้อมูล JSON จาก IPFS ของ Product Lot (Quality & Nutrition)
	ipfsCID := productLotData.QualityAndNutritionCID
	ipfsData, err := tc.IPFSService.GetJSONFromIPFS(ipfsCID)
	if err != nil {
		fmt.Println("❌ Failed to fetch quality & nutrition data from IPFS:", err)
		return nil, fmt.Errorf("Failed to fetch quality & nutrition data")
	}

	// ✅ ตรวจสอบว่า `qualityData` มีอยู่จริงหรือไม่
	qualityDataMap, ok := ipfsData["qualityData"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Error: qualityData is missing or incorrect")
		return nil, fmt.Errorf("qualityData structure is incorrect")
	}

	// ✅ ดึงข้อมูลคุณภาพ
	qualityData, ok := qualityDataMap["quality"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Error: Quality data is missing or incorrect")
		return nil, fmt.Errorf("Quality data structure is incorrect")
	}

	// ✅ แปลง `grade` เป็นข้อความ
	var gradeText string
	if productLotData.Grade {
		gradeText = "Passed"
	} else {
		gradeText = "Failed"
	}

	// ✅ แปลง `inspectionDate` เป็น `YYYY-MM-DD HH:mm:ss`
	inspectionTime := time.Unix(productLotData.InspectionDate.Unix(), 0).Format("2006-01-02 15:04:05")

	// ✅ สร้าง JSON Response
	response := fiber.Map{
		"GeneralInfo": fiber.Map{
			"productId":    productID,
			"productName":  productData["productName"],
			"category":     productData["category"],
			"description":  productIPFSData["description"],
			"quantity":     productIPFSData["quantity"],
			"quantityUnit": nutritionData["quantityUnit"],
		},
		"selectMilkTank": fiber.Map{
			"temp":            qualityData["temp"],
			"tempUnit":        qualityData["tempUnit"],
			"pH":              qualityData["pH"],
			"fat":             qualityData["fat"],
			"protein":         qualityData["protein"],
			"bacteria":        qualityData["bacteria"],
			"bacteriaInfo":    qualityData["bacteriaInfo"],
			"contaminants":    qualityData["contaminants"],
			"contaminantInfo": qualityData["contaminantInfo"],
			"abnormalChar":    qualityData["abnormalChar"],
			"abnormalType":    qualityData["abnormalType"],
		},
		"Quality": fiber.Map{
			"grade":          gradeText,
			"inspectionDate": inspectionTime,
			"inspector":      productLotData.Inspector,
		},
		"nutrition": nutritionData,
		"factoryInfo": fiber.Map{
			"factoryID":     factory.FactoryID,
			"companyName":   factory.CompanyName,
			"address":       factory.Address,
			"district":      factory.District,
			"subDistrict":   factory.SubDistrict,
			"province":      factory.Province,
			"country":       factory.Country,
			"postCode":      factory.PostCode,
			"telephone":     factory.Telephone,
			"email":         factory.Email,
			"lineID":        factory.LineID.String,
			"facebook":      factory.Facebook.String,
			"locationLink":  factory.LocationLink.String,
			"createdOn":     factory.CreatedOn.Format("2006-01-02 15:04:05"),
			"walletAddress": factory.WalletAddress,
		},
	}

	// ✅ ส่งข้อมูลกลับไปยัง Frontend
	return response, nil
}
