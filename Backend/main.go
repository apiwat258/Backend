package main

import (
	"fmt"
	"log"
	"os"

	"finalyearproject/Backend/api/controllers"
	"finalyearproject/Backend/api/routes"
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// ✅ โหลดไฟล์ `.env`
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// ✅ Debug ตรวจสอบค่า Environment Variables
	fmt.Println("📌 DEBUG - BLOCKCHAIN_RPC_URL:", os.Getenv("BLOCKCHAIN_RPC_URL"))
	fmt.Println("📌 DEBUG - JWT_SECRET:", os.Getenv("JWT_SECRET"))

	// ✅ เริ่มต้น Fiber App
	app := fiber.New()

	// ✅ แก้ CORS Policy (ลบ `/` ท้าย `AllowOrigins`)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://192.168.43.218:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	// ✅ เชื่อมต่อฐานข้อมูล
	database.Connect()

	// ✅ Migrate Tables
	database.DB.AutoMigrate(&models.User{}, &models.Farmer{}, &models.Certification{}, &models.Logistics{}, &models.Factory{}, &models.Retailer{})

	// ✅ เริ่มต้น Blockchain Service
	err = services.InitBlockchainService()
	if err != nil {
		log.Fatalf("❌ Blockchain Service Error: %v", err)
	}
	fmt.Println("✅ Blockchain Service Started Successfully!")

	// ✅ เริ่มต้น IPFS Service
	services.InitIPFSService()

	// ✅ ตรวจสอบว่า IPFSServiceInstance ถูกกำหนดค่าก่อนใช้งาน
	if services.IPFSServiceInstance == nil {
		log.Fatal("❌ IPFS Service failed to initialize. Exiting...")
	}

	// ✅ กำหนดค่าให้ QRCodeServiceInstance ใช้ IPFSServiceInstance ที่ถูกต้อง
	services.QRCodeServiceInstance = &services.QRCodeService{
		IPFSService: services.IPFSServiceInstance,
	}

	// ✅ สร้าง RawMilkController
	rawMilkController := controllers.NewRawMilkController(
		database.DB,
		services.BlockchainServiceInstance,
		services.IPFSServiceInstance,
		services.QRCodeServiceInstance,
	)

	// ✅ สร้าง ProductController
	productController := controllers.NewProductController(
		database.DB,
		services.BlockchainServiceInstance,
		services.IPFSServiceInstance,
	)

	// ✅ สร้าง ProductLotController
	productLotController := controllers.NewProductLotController(
		database.DB, // ✅ ใช้ database.DB ที่เป็น *gorm.DB
		services.BlockchainServiceInstance,
		services.IPFSServiceInstance,
		services.QRCodeServiceInstance,
	)

	trackingController := controllers.NewTrackingController(
		database.DB,
		services.BlockchainServiceInstance,
		services.IPFSServiceInstance,
		services.QRCodeServiceInstance,
	)

	// ✅ ส่ง Controller ไปที่ `SetupRoutes`
	routes.SetupRoutes(app, rawMilkController, productController, productLotController, trackingController)
	// ✅ ให้บริการไฟล์ Static (Frontend)
	app.Static("/", "./frontend")

	// ✅ เริ่มเซิร์ฟเวอร์
	log.Fatal(app.Listen(":8080"))
}
