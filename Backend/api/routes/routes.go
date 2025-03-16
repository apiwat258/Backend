package routes

import (
	"finalyearproject/Backend/api/controllers"
	"finalyearproject/Backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up API routes
func SetupRoutes(app *fiber.App, rmc *controllers.RawMilkController, pc *controllers.ProductController, plc *controllers.ProductLotController, tc *controllers.TrackingController) {
	api := app.Group("/api/v1")                              // ใช้ Prefix "/api/v1" สำหรับ API ทั้งหมด
	api.Get("/tracking-details", tc.GetTrackingDetailsByLot) // ✅ ต้องอยู่ใต้ `api`

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
	retailer.Get("/list", controllers.GetAllRetailers) // ดึง retailerID และชื่อร้านค้าทั้งหมด
	retailer.Get("/usernames", controllers.GetRetailerUsernames)
	retailer.Get("/:id", controllers.GetRetailerByID) // ดึงข้อมูลร้านค้าตาม retailerID

	certification := api.Group("/certifications")
	certification.Post("/upload", controllers.UploadCertificate)
	certification.Get("/me", middleware.AuthMiddleware(), controllers.GetCertificationByUser)
	certification.Delete("/", middleware.AuthMiddleware(), controllers.DeleteCertification) // ✅ API ต้องการ Auth
	certification.Get("/check/:certCID", controllers.CheckCertificationCID)
	certification.Post("/store", middleware.AuthMiddleware(), controllers.StoreCertification)

	// ✅ Milk Tank Routes
	milk := api.Group("/farm/milk", middleware.AuthMiddleware())
	milk.Post("/create", rmc.CreateMilkTank)
	milk.Get("/list", rmc.GetFarmRawMilkTanks)
	milk.Get("/details/:tankId", rmc.GetRawMilkTankDetails)
	milk.Get("/qrcode/:tankId", rmc.GetQRCodeByTankID)

	// ✅ Milk Tank Routes สำหรับโรงงาน
	factoryMilk := api.Group("/factory/milk", middleware.AuthMiddleware())
	factoryMilk.Get("/list", rmc.GetFactoryRawMilkTanks)
	factoryMilk.Post("/update-status", rmc.UpdateMilkTankStatus)

	// ✅ Product Routes
	product := api.Group("/products", middleware.AuthMiddleware())
	product.Post("/create", pc.CreateProduct)
	product.Get("/list", pc.GetFactoryProducts)
	product.Get("/:productId", pc.GetProductDetails)

	// ✅ Product Lot Routes (ใหม่)
	productLot := api.Group("/product-lots", middleware.AuthMiddleware())
	productLot.Post("/create", plc.CreateProductLot)    // ✅ โรงงานสร้าง Product Lot ใหม่
	productLot.Get("/list", plc.GetFactoryProductLots)  // ✅ ดึงรายการ Product Lot ของโรงงาน
	productLot.Get("/:lotId", plc.GetProductLotDetails) // ✅ ดึงรายละเอียด Product Lot โดยใช้ lotId

	// ✅ กลุ่ม Routing ของ Tracking
	tracking := api.Group("/tracking", middleware.AuthMiddleware())
	tracking.Get("/ids", plc.GetAllTrackingIds)                                     // ✅ ดึง Tracking ID ทั้งหมด
	tracking.Post("/logistics", plc.UpdateLogisticsCheckpoint)                      // ✅ อัปเดตจุดตรวจโลจิสติกส์
	tracking.Get("/logistics/checkpoints", plc.GetLogisticsCheckpointsByTrackingId) // ✅ ดึงข้อมูล Checkpoint ตาม Tracking ID
	tracking.Get("/retailer", plc.GetRetailerTracking)                              // ✅ ดึง Tracking IDs ตาม Retailer ID
	tracking.Post("/retailer/receive", plc.RetailerReceiveProduct)                  // ✅ Retailer รับสินค้า
	tracking.Get("/retailer/received", plc.GetRetailerReceivedProduct)              // ✅ ดึงข้อมูลสินค้าที่ Retailer รับแล้ว

}
