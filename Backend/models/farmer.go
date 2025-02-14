package models

import (
	"database/sql"
	"time"
)

type Farmer struct {
	FarmerID      string         `gorm:"primaryKey;column:farmerid"`
	UserID        string         `gorm:"column:userid;unique"`
	FarmerName    string         `gorm:"column:farmer_name"`
	CompanyName   string         `gorm:"column:companyname"` // ✅ ใช้ชื่อให้ตรงกับ DB
	Address       string         `gorm:"column:address"`
	City          string         `gorm:"column:city"`
	Province      string         `gorm:"column:province"`
	Country       string         `gorm:"column:country;default:Thailand"`
	PostCode      string         `gorm:"column:postcode"`
	Telephone     string         `gorm:"column:telephone"`
	LineID        sql.NullString `gorm:"column:lineid"`
	Facebook      sql.NullString `gorm:"column:facebook"`
	LocationLink  sql.NullString `gorm:"column:location_link"`
	CreatedOn     time.Time      `gorm:"column:createdon;autoCreateTime"`
	Email         string         `gorm:"column:email;unique"`
	WalletAddress string         `gorm:"not null"` // ✅ เพิ่มฟิลด์นี้
}

// ✅ บังคับ GORM ให้ใช้ตาราง `farmer` ที่มีอยู่
func (Farmer) TableName() string {
	return "farmer"
}
