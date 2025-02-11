package tests

import (
	"bytes"
	"encoding/json"
	"finalyearproject/Backend/api/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// ✅ Set up Fiber app
	app := fiber.New()
	app.Post("/auth/register", controllers.Register)

	// ✅ Mock User Registration Request
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "testuser@example.com",
		"password": "TestPassword123",
	})

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	// ✅ ตรวจสอบผลลัพธ์
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
