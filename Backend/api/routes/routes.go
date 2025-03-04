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
	auth.Get("/check-email", controllers.CheckEmailAvailability)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)       // ✅ เพิ่มเส้นทาง Logout
	auth.Get("/get-role", controllers.GetUserRole) // ✅ เพิ่ม API ดึง Role ของผู้ใช้
	auth.Get("/user-info", middleware.AuthMiddleware(), controllers.GetUserInfo)
	auth.Put("/update-user", middleware.AuthMiddleware(), controllers.UpdateUserInfo)
	api.Post("/refresh-token", controllers.RefreshTokenHandler)

	// ✅ Protected Routes (ใช้ Middleware JWT)
	protected := api.Group("/protected", middleware.AuthMiddleware())
	protected.Get("/route", controllers.ProtectedRoute)

	// ✅ Farmer Routes
	farmer := api.Group("/farmers")
	farmer.Post("/create", middleware.AuthMiddleware(), controllers.CreateFarmer)
	farmer.Get("/me", middleware.AuthMiddleware(), controllers.GetFarmByUser)
	farmer.Put("/update", middleware.AuthMiddleware(), controllers.UpdateFarmer)

	// ✅ Raw Milk Routes (เกษตรกรใช้เพิ่มข้อมูลน้ำนมดิบ)
	rawMilk := api.Group("/rawmilk")
	rawMilk.Post("/", middleware.AuthMiddleware(), controllers.AddRawMilkHandler)              // 🔐 ใช้ Middleware เฉพาะ POST
	rawMilk.Get("/:id", controllers.GetRawMilkHandler)                                         // ✅ เอา Middleware ออก
	rawMilk.Post("/upload", middleware.AuthMiddleware(), controllers.UploadRawMilkFileHandler) // ✅ ใหม่: อัปโหลดไฟล์ JSON ไป IPFS
	rawMilk.Get("/ipfs/:cid", controllers.GetRawMilkFromIPFSHandler)                           // ✅ ใหม่: ดึง JSON จาก IPFS

	// ✅ Factory Routes
	factory := api.Group("/factories")
	factory.Post("/", middleware.AuthMiddleware(), controllers.CreateFactory)
	factory.Get("/", middleware.AuthMiddleware(), controllers.GetFactoryByUser)
	factory.Put("/", middleware.AuthMiddleware(), controllers.UpdateFactory)

	// ✅ Logistics Routes
	logistics := api.Group("/logistics")
	logistics.Post("/", middleware.AuthMiddleware(), controllers.CreateLogistics)
	logistics.Get("/", middleware.AuthMiddleware(), controllers.GetLogisticsByUser)
	logistics.Put("/", middleware.AuthMiddleware(), controllers.UpdateLogistics)

	retailer := api.Group("/retailers")
	retailer.Post("/", middleware.AuthMiddleware(), controllers.CreateRetailer)
	retailer.Get("/", middleware.AuthMiddleware(), controllers.GetRetailerByUser)
	retailer.Put("/", middleware.AuthMiddleware(), controllers.UpdateRetailer)

	// ✅ Certification Routes (เพิ่มเส้นทางสำหรับลบใบเซอร์)
	certification := api.Group("/certifications")
	certification.Post("/upload", controllers.UploadCertificate)
	certification.Get("/me", middleware.AuthMiddleware(), controllers.GetCertificationByUser)
	certification.Delete("/", middleware.AuthMiddleware(), controllers.DeleteCertification) // ✅ API ต้องการ Auth
	certification.Get("/check/:certCID", controllers.CheckCertificationCID)
	certification.Post("/store", middleware.AuthMiddleware(), controllers.StoreCertification)

	// ✅ QR Code Routes (ใหม่)
	qr := api.Group("/qr")
	qr.Get("/rawmilk/:id", controllers.GenerateQRCodeHandler) // ✅ ใช้ "/api/v1/qr/rawmilk/:id"

}
