package main

import (
	"finalyearproject/Backend/api/routes"
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// ✅ อนุญาตให้ Frontend เชื่อมต่อ API
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://0.0.0.0:8081", // ✅ ตรงกับ Frontend
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
	}))

	// ✅ เชื่อมต่อฐานข้อมูล
	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Farmer{}, &models.Certification{})

	// ✅ ให้บริการไฟล์ Static (Frontend)
	app.Static("/", "./frontend")

	// ✅ เรียกใช้ SetupRoutes() เพื่อกำหนดเส้นทาง API ทั้งหมด
	routes.SetupRoutes(app)

	// ✅ เริ่มเซิร์ฟเวอร์
	app.Listen(":8080")
}
