package models

import "time"

// TrackingStatus - บันทึกสถานะแทรคกิ้งไอดี
type TrackingStatus struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TrackingID string    `gorm:"unique;not null" json:"tracking_id"`
	Status     int       `gorm:"not null;default:0" json:"status"` // 0 = Pending, 1 = Received by Logistics, 2 = Received by Farm
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (TrackingStatus) TableName() string {
	return "tracking_status" // <--- ตรงตาม DB ที่คุณสร้างไว้
}
