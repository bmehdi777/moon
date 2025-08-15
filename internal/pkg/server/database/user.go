package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	KCUserID string
}

func FindUserByKCUID(kcUid string, db *gorm.DB) (*User, *gorm.DB) {
	var user User
	res := db.First(&user, "kc_user_id = ?", kcUid)
	return &user, res
}
