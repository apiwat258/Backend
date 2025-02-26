package models

import (
	"database/sql"
	"time"
)

type Factory struct {
	FactoryID     string         `gorm:"primaryKey;column:factoryid"`
	CompanyName   string         `gorm:"column:companyname"`
	Address       string         `gorm:"column:address"`
	District      string         `gorm:"column:district"`
	SubDistrict   string         `gorm:"column:subdistrict"`
	Province      string         `gorm:"column:province"`
	Country       string         `gorm:"column:country;default:Thailand"`
	PostCode      string         `gorm:"column:postcode"`
	Telephone     string         `gorm:"column:telephone"`
	LineID        sql.NullString `gorm:"column:lineid"`
	Facebook      sql.NullString `gorm:"column:facebook"`
	LocationLink  sql.NullString `gorm:"column:location_link"`
	CreatedOn     time.Time      `gorm:"column:createdon;autoCreateTime"`
	WalletAddress string         `gorm:"column:wallet_address;not null"`
}

// ✅ ใช้ตาราง `dairyfactory`
func (Factory) TableName() string {
	return "dairyfactory"
}
