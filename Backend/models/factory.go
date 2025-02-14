package models

import (
	"database/sql"
	"time"
)

type Factory struct {
	FactoryID    string         `gorm:"primaryKey;column:factoryid"`
	UserID       string         `gorm:"column:userid;unique;not null"`
	Username     string         `gorm:"column:username"` // ✅ ใช้ `Username` แทน `FactoryName`
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

// ✅ ใช้ตาราง `dairyfactory`
func (Factory) TableName() string {
	return "dairyfactory"
}
