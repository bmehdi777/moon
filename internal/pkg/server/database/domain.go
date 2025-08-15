package database

import "gorm.io/gorm"

type DomainName struct {
	gorm.Model
	FQDN string
}

func CreateDomainName(name string, db *gorm.DB) {
	domain := DomainName{FQDN: name}
	db.Create(&domain)
}

func FindDomainNameById(id uint, db *gorm.DB) (*DomainName, *gorm.DB) {
	var domain DomainName
	res := db.First(&domain, "id = ?", id)
	return &domain, res
}

func FindDomainNameByName(fqdn string, db *gorm.DB) (*DomainName, *gorm.DB) {
	var domain DomainName
	res := db.First(&domain, "fqdn = ?", fqdn)
	return &domain, res
}
