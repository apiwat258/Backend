package models

import (
	"database/sql"
	"time"
)

type Farmer struct {
	FarmerID      string         `gorm:"primaryKey;column:farmerid"`
	CompanyName   string         `gorm:"column:companyname"`
	Address       string         `gorm:"column:address"`
	District      string         `gorm:"column:district"`
	SubDistrict   string         `gorm:"column:subdistrict"`
	Province      string         `gorm:"column:province"`
	Country       string         `gorm:"column:country;default:Thailand"`
	PostCode      string         `gorm:"column:postcode"`
	Telephone     string         `gorm:"column:telephone"`
	Email         string         `gorm:"column:email;unique;not null"` // ✅ เพิ่ม Email ของฟาร์ม
	LineID        sql.NullString `gorm:"column:lineid"`
	Facebook      sql.NullString `gorm:"column:facebook"`
	LocationLink  string         `gorm:"column:location_link"`
	CreatedOn     time.Time      `gorm:"column:createdon;autoCreateTime"`
	WalletAddress string         `gorm:"column:wallet_address;not null"`
}

// ✅ ใช้ตาราง `farmer`
func (Farmer) TableName() string {
	return "farmer"
}
