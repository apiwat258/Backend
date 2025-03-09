package routes

import (
	"finalyearproject/Backend/api/controllers"
	"finalyearproject/Backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up API routes
func SetupRoutes(app *fiber.App, rmc *controllers.RawMilkController, pc *controllers.ProductController) {
	api := app.Group("/api/v1") // ใช้ Prefix "/api/v1" สำหรับ API ทั้งหมด

	// ✅ Authentication Routes
	auth := api.Group("/auth")
	auth.Post("/update-role", controllers.UpdateUserRole)
	auth.Post("/register", controllers.Register)
	auth.Get("/check-email", controllers.CheckEmailAvailability)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)
	auth.Get("/get-role", controllers.GetUserRole)
	auth.Get("/user-info", middleware.AuthMiddleware(), controllers.GetUserInfo)
	auth.Put("/update-user", middleware.AuthMiddleware(), controllers.UpdateUserInfo)
	api.Post("/refresh-token", controllers.RefreshTokenHandler)

	protected := api.Group("/protected", middleware.AuthMiddleware())
	protected.Get("/route", controllers.ProtectedRoute)

	farmer := api.Group("/farmers")
	farmer.Post("/create", middleware.AuthMiddleware(), controllers.CreateFarmer)
	farmer.Get("/me", middleware.AuthMiddleware(), controllers.GetFarmByUser)
	farmer.Put("/update", middleware.AuthMiddleware(), controllers.UpdateFarmer)

	factory := api.Group("/factories")
	factory.Post("/", middleware.AuthMiddleware(), controllers.CreateFactory)
	factory.Get("/", middleware.AuthMiddleware(), controllers.GetFactoryByUser)
	factory.Put("/", middleware.AuthMiddleware(), controllers.UpdateFactory)

	logistics := api.Group("/logistics")
	logistics.Post("/", middleware.AuthMiddleware(), controllers.CreateLogistics)
	logistics.Get("/", middleware.AuthMiddleware(), controllers.GetLogisticsByUser)
	logistics.Put("/", middleware.AuthMiddleware(), controllers.UpdateLogistics)

	retailer := api.Group("/retailers")
	retailer.Post("/", middleware.AuthMiddleware(), controllers.CreateRetailer)
	retailer.Get("/", middleware.AuthMiddleware(), controllers.GetRetailerByUser)
	retailer.Put("/", middleware.AuthMiddleware(), controllers.UpdateRetailer)

	certification := api.Group("/certifications")
	certification.Post("/upload", controllers.UploadCertificate)
	certification.Get("/me", middleware.AuthMiddleware(), controllers.GetCertificationByUser)
	certification.Delete("/", middleware.AuthMiddleware(), controllers.DeleteCertification) // ✅ API ต้องการ Auth
	certification.Get("/check/:certCID", controllers.CheckCertificationCID)
	certification.Post("/store", middleware.AuthMiddleware(), controllers.StoreCertification)

	// ✅ Milk Tank Routes (แก้ให้เรียกผ่าน `rmc.CreateMilkTank`)
	milk := api.Group("/farm/milk", middleware.AuthMiddleware())
	milk.Post("/create", rmc.CreateMilkTank)                // ✅ ฟาร์มสร้างแท็งก์นมดิบใหม่
	milk.Get("/list", rmc.GetFarmRawMilkTanks)              // ✅ ฟาร์มดึงรายการแท็งก์นมดิบของตัวเอง
	milk.Get("/details/:tankId", rmc.GetRawMilkTankDetails) // ✅ ดึงรายละเอียดแท็งก์นมดิบตาม Tank ID
	milk.Get("/qrcode/:tankId", rmc.GetQRCodeByTankID)      // ✅ ดึง QR Code ของแท็งก์นมดิบ

	// ✅ Milk Tank Routes สำหรับโรงงาน
	factoryMilk := api.Group("/factory/milk", middleware.AuthMiddleware())
	factoryMilk.Get("/list", rmc.GetFactoryRawMilkTanks) // ✅ โรงงานดึงรายการแท็งก์นมดิบที่ได้รับ
	factoryMilk.Post("/update-status", rmc.UpdateMilkTankStatus)

	// ✅ Product Routes (เพิ่มเข้ามา)
	product := api.Group("/products", middleware.AuthMiddleware())
	product.Post("/create", pc.CreateProduct) // ✅ โรงงานสร้างสินค้าใหม่
	//product.Get("/list", pc.GetFactoryProducts)  // ✅ ดึงสินค้าของโรงงาน
	//product.Get("/:productId", pc.GetProductByID) // ✅ ดึงรายละเอียดสินค้าโดยใช้ Product ID
}
