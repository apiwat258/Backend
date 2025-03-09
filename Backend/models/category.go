package models

import "gorm.io/gorm"

// ✅ Category Model
type Category struct {
	CategoryID uint   `gorm:"primaryKey;autoIncrement"` // ✅ ใช้ auto-increment
	Name       string `gorm:"unique;not null"`
}

// ✅ ฟังก์ชัน Migrate Category Table
func MigrateCategory(db *gorm.DB) error {
	return db.AutoMigrate(&Category{})
}
