package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv โหลดไฟล์ .env เพื่อใช้ตั้งค่าระบบ
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables.")
	}
}
