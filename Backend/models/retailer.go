package models

import (
	"database/sql"
	"time"
)

type Retailer struct {
	RetailerID    string         `gorm:"primaryKey;column:retailerid" json:"retailer_id"`
	CompanyName   string         `gorm:"column:companyname" json:"company_name"`
	Address       string         `gorm:"column:address" json:"address"`
	District      string         `gorm:"column:district" json:"district"`
	SubDistrict   string         `gorm:"column:subdistrict" json:"subdistrict"`
	Province      string         `gorm:"column:province" json:"province"`
	Country       string         `gorm:"column:country;default:Thailand" json:"country"`
	PostCode      string         `gorm:"column:postcode" json:"post_code"`
	Telephone     string         `gorm:"column:telephone" json:"telephone"`
	Email         string         `gorm:"column:email;unique;not null" json:"email"`
	LineID        sql.NullString `gorm:"column:lineid" json:"line_id"`
	Facebook      sql.NullString `gorm:"column:facebook" json:"facebook"`
	LocationLink  string         `gorm:"column:location_link" json:"location_link"`
	CreatedOn     time.Time      `gorm:"column:createdon;autoCreateTime" json:"created_on"`
	WalletAddress string         `gorm:"column:wallet_address;not null" json:"wallet_address"`
}

// ✅ ใช้ตาราง `retailer`
func (Retailer) TableName() string {
	return "retailer"
}
