package models

import (
	"time"
)

type Certification struct {
	CertificationID   string    `gorm:"primaryKey;column:certificationid"`
	EntityID          string    `gorm:"column:entityid"`
	EntityType        string    `gorm:"column:entitytype"`
	CertificationType string    `gorm:"column:certificationtype"`
	CertificationCID  string    `gorm:"column:certificationcid;unique"`
	EffectiveDate     time.Time `gorm:"column:effective_date"`
	IssuedDate        time.Time `gorm:"column:issued_date"`
	CreatedOn         time.Time `gorm:"column:created_on;autoCreateTime"`
	BlockchainTxHash  string    `gorm:"column:blockchain_tx"`
}

// ✅ บังคับ GORM ให้ใช้ตาราง `organiccertification` ที่มีอยู่
func (Certification) TableName() string {
	return "organiccertification"
}
