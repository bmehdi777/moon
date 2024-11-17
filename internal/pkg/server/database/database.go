package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDBConn() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("moon.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
