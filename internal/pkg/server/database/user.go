package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email   string
	Records DomainRecord
}

type DomainRecord struct {
	gorm.Model
	DNSRecord string
	ConnectionOpen  bool
}


func FindDomainRecordByName(name string, db *gorm.DB) *DomainRecord {
	var record DomainRecord
	db.First(&record, "DNSRecord = ?", name)
	return &record
}
