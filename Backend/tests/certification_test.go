package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"finalyearproject/Backend/api/controllers" // ✅ นำเข้า controllers

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUploadCertificateCertification(t *testing.T) {
	app := fiber.New()
	app.Post("/upload", controllers.UploadCertificate)

	// ✅ จำลองไฟล์อัปโหลด
	fileContent := []byte("test file content")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write(fileContent)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
