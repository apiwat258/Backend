package main

import (
	"fmt"
	"log"
	"os"

	"finalyearproject/Backend/api/routes"
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv" // ✅ เพิ่มการโหลดไฟล์ .env
)

func main() {
	// ✅ โหลดไฟล์ `.env` ก่อน
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// ✅ Debug ตรวจสอบว่า `BLOCKCHAIN_RPC_URL` โหลดถูกต้องหรือไม่
	fmt.Println("📌 DEBUG - BLOCKCHAIN_RPC_URL:", os.Getenv("BLOCKCHAIN_RPC_URL"))

	fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))

	// ✅ เริ่มต้น Fiber App
	app := fiber.New()

	// ✅ อนุญาตให้ Frontend เชื่อมต่อ API
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://0.0.0.0:8081",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
	}))

	// ✅ เชื่อมต่อฐานข้อมูล
	database.Connect()

	// ✅ Migrate Tables (Users, Farmers, Certifications)
	database.DB.AutoMigrate(&models.User{}, &models.Farmer{}, &models.Certification{}, &models.Logistics{}, &models.Factory{}, &models.Retailer{})

	// ✅ เริ่มต้น Blockchain Service
	err = services.InitBlockchainService()
	if err != nil {
		log.Fatalf("❌ Blockchain Service Error: %v", err)
	}
	fmt.Println("✅ Blockchain Service Started Successfully!")

	// ✅ กำหนดเส้นทาง API ทั้งหมด
	routes.SetupRoutes(app)

	// ✅ ให้บริการไฟล์ Static (Frontend)
	app.Static("/", "./frontend")

	// ✅ เริ่มเซิร์ฟเวอร์
	log.Fatal(app.Listen(":8080"))
}
