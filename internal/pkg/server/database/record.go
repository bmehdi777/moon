package database

import "gorm.io/gorm"

type Record struct {
	UserID         uint `gorm:"primaryKey"`
	DomainNameID   uint `gorm:"primaryKey"`
	ConnectionOpen bool

	User       User
	DomainName DomainName
}

func FindRecordByUserID(uid uint, db *gorm.DB) (*Record, *gorm.DB) {
	var record Record
	res := db.First(&record, "user_id = ?", uid)
	return &record, res
}

func FindRecordByDomainFQDN(name string, db *gorm.DB) (*Record, *gorm.DB) {
	domain, _ := FindDomainNameByName(name, db)
	var record Record
	res := db.First(&record, "domain_name_id = ?", domain.ID)
	return &record, res
}
