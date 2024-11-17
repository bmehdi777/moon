package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email                string
	DomainNameAttributed *string
	ConnectionOpen       bool
}

func FindUserByDomainName(name string, db *gorm.DB) (*User) {
	var user User
	db.First(&user, "DomainNameAttributed = ?", name)

	return &user
}
