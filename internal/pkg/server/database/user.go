package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email   string
	DomainRecordID int
	DomainRecord DomainRecord
}

type DomainRecord struct {
	gorm.Model
	DNSRecord      string
	ConnectionOpen bool
}

func UpdateDomainRecord(domainRec DomainRecord, db *gorm.DB) {
}

func FindDomainRecordByUserEmail(email string, db *gorm.DB) *DomainRecord {
	var record DomainRecord
	db.InnerJoins("User").First(&record, "email = ?", email)
	return &record
}

func FindDomainRecordByName(name string, db *gorm.DB) *DomainRecord {
	var record DomainRecord
	db.First(&record, "dns_record = ?", name)
	return &record
}
