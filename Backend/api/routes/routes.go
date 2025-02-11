package routes

import (
	"finalyearproject/Backend/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1") // กำหนด Prefix "/api/v1" ให้ทุก API

	// ✅ Authentication Routes
	auth := api.Group("/auth")
	auth.Post("/update-role", controllers.UpdateUserRole)
	auth.Post("/register", controllers.Register) // ✅ เส้นทาง /auth/register

	// ✅ Farmer Routes
	farmer := api.Group("/farmers")
	farmer.Post("/", controllers.CreateFarmer)

	// ✅ Certification Routes
	certification := api.Group("/certifications")
	certification.Post("/upload", controllers.UploadCertificate) // ✅ อัปโหลดไฟล์ไป IPFS
	certification.Post("/create", controllers.CreateCertification)
}
