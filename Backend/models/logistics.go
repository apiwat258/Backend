package models

import (
	"database/sql"
	"time"
)

type Logistics struct {
	LogisticsID  string         `gorm:"primaryKey;column:logisticsid"`
	UserID       string         `gorm:"column:userid;unique"`
	Username     string         `gorm:"column:username"` // ✅ ใช้ `Username`
	CompanyName  string         `gorm:"column:companyname"`
	Address      string         `gorm:"column:address"`
	City         string         `gorm:"column:city"`
	Province     string         `gorm:"column:province"`
	Country      string         `gorm:"column:country;default:Thailand"`
	PostCode     string         `gorm:"column:postcode"`
	Telephone    string         `gorm:"column:telephone"`
	LineID       sql.NullString `gorm:"column:lineid"`
	Facebook     sql.NullString `gorm:"column:facebook"`
	LocationLink sql.NullString `gorm:"column:location_link"`
	CreatedOn    time.Time      `gorm:"column:createdon;autoCreateTime"`
	Email        string         `gorm:"column:email;unique"`
}

// ✅ ใช้ตาราง `logisticsprovider`
func (Logistics) TableName() string {
	return "logisticsprovider"
}
