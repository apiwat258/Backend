package routes

import (
	"finalyearproject/Backend/api/controllers"
	"finalyearproject/Backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up API routes
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1") // ใช้ Prefix "/api/v1" สำหรับ API ทั้งหมด

	// ✅ Authentication Routes
	auth := api.Group("/auth")
	auth.Post("/update-role", controllers.UpdateUserRole)
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)       // ✅ เพิ่มเส้นทาง Logout
	auth.Get("/get-role", controllers.GetUserRole) // ✅ เพิ่ม API ดึง Role ของผู้ใช้

	// ✅ Protected Routes (ใช้ Middleware JWT)
	protected := api.Group("/protected", middleware.AuthMiddleware())
	protected.Get("/route", controllers.ProtectedRoute)

	// ✅ Farmer Routes
	farmer := api.Group("/farmers")
	farmer.Post("/", controllers.CreateFarmer)
	farmer.Get("/:id", controllers.GetFarmerByID) // ✅ เพิ่ม API สำหรับดึงข้อมูล Farmer ตาม ID
	// ✅ Route สำหรับอัปเดตข้อมูลฟาร์ม
	farmer.Put("/update", middleware.AuthMiddleware(), controllers.UpdateFarmer)
	farmer.Get("/me", middleware.AuthMiddleware(), controllers.GetFarmerByUser) // ✅ ใช้ Middleware

	// ✅ Raw Milk Routes (เกษตรกรใช้เพิ่มข้อมูลน้ำนมดิบ)
	rawMilk := api.Group("/rawmilk")
	rawMilk.Post("/", middleware.AuthMiddleware(), controllers.AddRawMilkHandler) // 🔐 ใช้ Middleware เฉพาะ POST
	rawMilk.Get("/:id", controllers.GetRawMilkHandler)                            // ✅ เอา Middleware ออก

	// ✅ Factory Routes
	factory := api.Group("/factories")
	factory.Post("/", controllers.CreateFactory)

	// ✅ Logistics Routes
	logistics := api.Group("/logistics")
	logistics.Post("/", controllers.CreateLogistics)

	// ✅ Retailer Routes
	retailer := api.Group("/retailers")
	retailer.Post("/", controllers.CreateRetailer)

	// ✅ Certification Routes
	certification := api.Group("/certifications")
	certification.Post("/upload", controllers.UploadCertificate)
	certification.Post("/create", controllers.CreateCertification)
}
