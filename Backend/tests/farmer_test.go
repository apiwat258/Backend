package tests

import (
	"bytes"
	"encoding/json"
	"finalyearproject/Backend/api/controllers"
	"finalyearproject/Backend/database"
	"finalyearproject/Backend/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateFarmer(t *testing.T) {
	// ✅ Mock Database Connection
	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Farmer{})

	// ✅ Set up Fiber app
	app := fiber.New()
	app.Post("/farmers/complete-profile", controllers.CreateFarmer)

	// ✅ Mock Farmer Registration Request
	requestBody, _ := json.Marshal(map[string]string{
		"userid":       "250001",
		"company_name": "Organic Farm Ltd",
		"firstname":    "John",
		"lastname":     "Doe",
		"email":        "johndoe@example.com",
		"address":      "123 Farm Road",
		"phone":        "+1234567890",
		"post":         "10000",
		"city":         "Bangkok",
		"province":     "Bangkok",
		"country":      "Thailand",
	})

	req := httptest.NewRequest("POST", "/farmers/complete-profile", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	// ✅ ตรวจสอบผลลัพธ์
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
