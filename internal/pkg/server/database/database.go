package database

import (
	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDBConn() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.GlobalConfig.Database.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&DomainRecord{})
	db.AutoMigrate(&User{})

	test := User{
		Email: "mehdi.bentouati",
	}
	db.Create(&test)
	return db, nil
}
