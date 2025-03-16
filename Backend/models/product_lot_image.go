package models

import (
	"time"
)

// ProductLotImage - โครงสร้างตารางสำหรับเก็บ Image CID, Tracking ID และ Person In Charge
type ProductLotImage struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LotID          string    `gorm:"uniqueIndex;not null" json:"lot_id"`
	ImageCID       string    `gorm:"not null" json:"image_cid"`
	TrackingIDs    string    `gorm:"not null" json:"tracking_ids"` // ✅ เพิ่มฟิลด์นี้
	PersonInCharge string    `gorm:"not null" json:"person_in_charge"`
	CreatedAt      time.Time `json:"created_at"`
}
