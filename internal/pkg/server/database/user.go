package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	KCUserId string
	DomainRecordID int
	DomainRecord DomainRecord
}

type DomainRecord struct {
	gorm.Model
	DNSRecord      string
	ConnectionOpen bool
}


func FindDomainRecordByName(name string, db *gorm.DB) *DomainRecord {
	var record DomainRecord
	db.First(&record, "dns_record = ?", name)
	return &record
}
