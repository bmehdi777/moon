package database

import (
	"moon/internal/pkg/server/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDBConn() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.GlobalConfig.Database.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&DomainRecord{})
	db.AutoMigrate(&User{})

	//test := User{
	//	KCUserId: "7e4ae9a4-efd8-44e1-8c26-53b0b9d10667",
	//}
	//db.Create(&test)
	return db, nil
}
