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

	// ✅ กำหนด CORS Origins → รองรับหลาย Origin
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		// Default เผื่อกรณีไม่มีตั้งค่า (Local dev)
		allowedOrigins = "http://127.0.0.1:3000, http://localhost:3000"
	}
	fmt.Println("📌 DEBUG - ALLOWED_ORIGINS:", allowedOrigins)

	fmt.Println("📌 DEBUG - ALLOWED_ORIGINS:", allowedOrigins)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins, // ✅ ENV
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	// ✅ เชื่อมต่อฐานข้อมูล
	database.Connect()

	// ✅ Migrate Tables
	database.DB.AutoMigrate(&models.User{}, &models.Farmer{}, &models.Logistics{}, &models.Factory{}, &models.Retailer{})

	// ✅ เริ่มต้น Blockchain Service
	err = services.InitBlockchainService()
	if err != nil {
		log.Fatalf("❌ Blockchain Service Error: %v", err)
	}
	fmt.Println("✅ Blockchain Service Started Successfully!")

	// ✅ เริ่มต้น IPFS Service
	services.InitIPFSService()

	if services.IPFSServiceInstance == nil {
		log.Fatal("❌ IPFS Service failed to initialize. Exiting...")
	}

	services.QRCodeServiceInstance = &services.QRCodeService{
		IPFSService: services.IPFSServiceInstance,
	}

	rawMilkController := controllers.NewRawMilkController(
		database.DB,
		services.BlockchainServiceInstance,
		services.IPFSServiceInstance,
		services.QRCodeServiceInstance,
	)

	productController := controllers.NewProductController(
		database.DB,
		services.BlockchainServiceInstance,
		services.IPFSServiceInstance,
	)

	productLotController := controllers.NewProductLotController(
		database.DB,
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

	// ✅ Setup Routes
	routes.SetupRoutes(app, rawMilkController, productController, productLotController, trackingController)

	// ✅ Serve Static (Frontend)
	app.Static("/", "./frontend")

	// ✅ Start Server
	log.Fatal(app.Listen(":8081"))
}
