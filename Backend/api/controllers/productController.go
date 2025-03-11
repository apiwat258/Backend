package controllers

import (
	"encoding/json"
	"finalyearproject/Backend/models"
	"finalyearproject/Backend/services"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ Struct สำหรับ Controller
type ProductController struct {
	DB                *gorm.DB
	BlockchainService *services.BlockchainService
	IPFSService       *services.IPFSService
	Mutex             sync.Mutex
}

// ✅ ฟังก์ชัน Constructor
func NewProductController(
	db *gorm.DB,
	blockchainService *services.BlockchainService,
	ipfsService *services.IPFSService,
) *ProductController {
	return &ProductController{
		DB:                db,
		BlockchainService: blockchainService,
		IPFSService:       ipfsService,
		Mutex:             sync.Mutex{},
	}
}

// ✅ ฟังก์ชันสร้าง Category Code ถ้าไม่มีอยู่ในระบบ
func (pc *ProductController) getOrCreateCategory(categoryName string) (uint, error) {
	pc.Mutex.Lock()
	defer pc.Mutex.Unlock()

	var category models.Category
	err := pc.DB.Where("name = ?", categoryName).First(&category).Error
	if err != nil {
		// ✅ ถ้ายังไม่มี Category ให้สร้างใหม่
		if err == gorm.ErrRecordNotFound {
			newCategory := models.Category{Name: categoryName}
			if err := pc.DB.Create(&newCategory).Error; err != nil {
				return 0, err
			}
			fmt.Println("✅ Created New Category:", newCategory.CategoryID)
			return newCategory.CategoryID, nil
		}
		return 0, err
	}

	return category.CategoryID, nil
}

// ✅ ฟังก์ชันแปลงเลขเป็น Base36
func toBase36(n int) string {
	const base36 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if n < 36 {
		return string(base36[n])
	}
	return toBase36(n/36) + string(base36[n%36])
}

// ✅ ฟังก์ชันสร้าง Product ID
func (pc *ProductController) generateProductID(factoryID string, categoryID uint) string {
	// ✅ ใช้ 3 หลักท้ายของ Factory ID
	shortFactoryID := factoryID[len(factoryID)-3:] // เช่น "FAC2500005" → "005"

	// ✅ แปลง Category ID เป็น Base36
	categoryBase36 := toBase36(int(categoryID)) // เช่น "12" → "C"

	// ✅ ใช้ Timestamp YYMMDD + รหัสสุ่ม 2 ตัว
	timestamp := time.Now().Format("060102") // YYMMDD
	randomCode := rand.Intn(36)              // Base36 (0-9, A-Z)
	randomBase36 := toBase36(randomCode)     // เช่น "X"

	// ✅ รวมเป็น Product ID
	return fmt.Sprintf("%s%s-%s%s", shortFactoryID, categoryBase36, timestamp, randomBase36)
}

// ✅ ฟังก์ชันสร้างสินค้า
func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Create Product")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)
	factoryID := c.Locals("entityID").(string)          // ✅ ใช้ Factory ID สำหรับ Product ID
	walletAddress := c.Locals("walletAddress").(string) // ✅ ใช้ Wallet สำหรับ Blockchain

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Factory เท่านั้นที่สามารถสร้างสินค้าได้)
	if role != "factory" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only factories can create products"})
	}

	// ✅ รับข้อมูล JSON ที่ส่งมา
	var request struct {
		GeneralInfo struct {
			ProdName    string `json:"prodName"`
			Category    string `json:"category"`
			Description string `json:"description"`
			Quantity    string `json:"quantity"`
		} `json:"GeneralInfo"`
		Nutrition json.RawMessage `json:"Nutrition"` // ✅ เก็บโภชนาการเป็น Raw JSON
	}

	// ✅ แปลง JSON ที่รับเข้ามา
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		fmt.Println("❌ Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ ตรวจสอบว่า Category ไม่เป็นค่าว่าง
	if strings.TrimSpace(request.GeneralInfo.Category) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Category cannot be empty"})
	}

	// ✅ ค้นหา/สร้าง Category ID
	categoryID, err := pc.getOrCreateCategory(request.GeneralInfo.Category)
	if err != nil {
		fmt.Println("❌ Failed to create/find category:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process category"})
	}

	// ✅ ใช้ Factory ID + Category ID สร้าง `Product ID`
	productID := pc.generateProductID(factoryID, categoryID)

	// ✅ รวมข้อมูลทั้งหมด (GeneralInfo + Nutrition) สำหรับ IPFS
	productMetadata := map[string]interface{}{
		"prodName":    request.GeneralInfo.ProdName,
		"category":    request.GeneralInfo.Category,
		"description": request.GeneralInfo.Description,
		"quantity":    request.GeneralInfo.Quantity,
		"nutrition":   json.RawMessage(request.Nutrition),
	}

	// ✅ อัปโหลดข้อมูลไป IPFS
	productCID, err := pc.IPFSService.UploadDataToIPFS(productMetadata)
	if err != nil {
		fmt.Println("❌ Failed to upload product data to IPFS:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload product data"})
	}

	// ✅ ส่งธุรกรรมไปที่ Blockchain
	txHash, err := pc.BlockchainService.CreateProduct(
		walletAddress, // ✅ ใช้ Wallet Address สำหรับ Blockchain
		productID,     // ✅ ใช้ Factory ID + Category ID สำหรับ Product ID
		request.GeneralInfo.ProdName,
		productCID,
		request.GeneralInfo.Category, // ✅ ส่ง Category เข้า Blockchain
	)
	if err != nil {
		fmt.Println("❌ Blockchain transaction failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Blockchain transaction failed"})
	}

	// ✅ ส่ง Response กลับไปให้ Frontend
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":   "Product created successfully",
		"productId": productID, // ✅ ใช้ Factory ID สร้าง Product ID
		"txHash":    txHash,
		"ipfsCID":   productCID,
		"category":  request.GeneralInfo.Category, // ✅ ส่ง Category ใน Response ด้วย
	})
}

func (pc *ProductController) GetFactoryProducts(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Factory Products")

	// ✅ ดึงข้อมูลจาก JWT Token
	role := c.Locals("role").(string)

	// ✅ ตรวจสอบสิทธิ์ (เฉพาะ Factory เท่านั้น)
	if role != "factory" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Only factories can view products"})
	}

	// ✅ ดึง `entityID` จาก JWT Token ที่ AuthMiddleware กำหนดไว้
	factoryID, ok := c.Locals("entityID").(string)
	if !ok || factoryID == "" {
		fmt.Println("❌ Factory ID is missing in Context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - Factory ID is missing"})
	}
	fmt.Println("✅ Factory ID from Context:", factoryID)

	// ✅ ดึง Wallet Address จาก JWT Token
	walletAddress, ok := c.Locals("walletAddress").(string)
	if !ok || walletAddress == "" {
		fmt.Println("❌ Wallet Address is missing in Context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - Wallet Address is missing"})
	}
	fmt.Println("✅ Wallet Address from Context:", walletAddress)

	// ✅ ดึงค่าที่พิมพ์ในช่องค้นหา (Search Query)
	searchQuery := strings.ToLower(c.Query("search", ""))

	// ✅ ดึงข้อมูลสินค้าในโรงงานจาก Blockchain (ใช้ Wallet Address แทน Factory ID)
	products, err := pc.BlockchainService.GetProductsByFactory(walletAddress)
	if err != nil {
		fmt.Println("❌ Failed to fetch factory products:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch factory products"})
	}

	// ✅ กรองข้อมูลตาม Search Query
	var filteredProducts []map[string]interface{}
	for _, product := range products {
		productId := product["productId"].(string)
		productName := product["productName"].(string)
		category := product["category"].(string)

		if searchQuery == "" ||
			strings.Contains(strings.ToLower(productId), searchQuery) ||
			strings.Contains(strings.ToLower(productName), searchQuery) ||
			strings.Contains(strings.ToLower(category), searchQuery) {

			filteredProducts = append(filteredProducts, map[string]interface{}{
				"productId":   strings.TrimRight(productId, "\x00"),
				"productName": productName,
				"category":    category,
				"detailsLink": fmt.Sprintf("/Factory/ProductDetails/%s", productId),
			})
		}
	}

	// ✅ ส่ง Response กลับไปที่ Frontend ให้มีมาตรฐานเดียวกัน
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"displayedProducts": filteredProducts,
	})
}

func (pc *ProductController) GetProductDetails(c *fiber.Ctx) error {
	fmt.Println("📌 Request received: Get Product Details")

	// ✅ ดึง `productId` จาก URL Parameter
	productID := c.Params("productId")
	if productID == "" {
		fmt.Println("❌ Product ID is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product ID is required"})
	}
	fmt.Println("✅ Product ID:", productID)

	// ✅ ดึงข้อมูล Product จาก Smart Contract
	productData, err := pc.BlockchainService.GetProductDetails(productID)
	if err != nil {
		fmt.Println("❌ Failed to fetch product from blockchain:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product details"})
	}

	// ✅ ดึงข้อมูล JSON จาก IPFS
	ipfsData, err := pc.IPFSService.GetJSONFromIPFS(productData["productCID"].(string))
	if err != nil {
		fmt.Println("❌ Failed to fetch product data from IPFS:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product data"})
	}

	// ✅ จัดรูปแบบ Response ให้ตรงกับ Frontend
	response := fiber.Map{
		"GeneralInfo": fiber.Map{
			"productName": productData["productName"],
			"category":    productData["category"],
			"description": ipfsData["description"],
			"quantity":    ipfsData["quantity"],
		},
		"Nutrition": ipfsData["nutrition"], // ✅ ดึง Nutrition จาก IPFS
	}

	// ✅ ส่งข้อมูลให้ Frontend
	return c.Status(http.StatusOK).JSON(response)
}
