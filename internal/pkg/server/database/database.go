package database

import (
	"fmt"
	"moon/internal/pkg/server/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDBConn() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch config.GlobalConfig.Database.Driver {
	case config.DRIVER_SQLITE:
		db, err = gorm.Open(sqlite.Open(config.GlobalConfig.Database.SqlitePath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		break
	case config.DRIVER_POSTGRES:
		dsn := fmt.Sprintf("user=%v password=%v dbname=%v port=%v sslmode=disable",
			config.GlobalConfig.Database.PostgresUser,
			config.GlobalConfig.Database.PostgresPassword,
			config.GlobalConfig.Database.PostgresDbName,
			config.GlobalConfig.Database.PostgresPort)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		break
	}

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&DomainRecord{})
	db.AutoMigrate(&User{})

	//test := User{
	//	KCUserId: "7b98b00f-fd7e-4d41-81d3-447a36c59200",
	//}
	//db.Create(&test)
	return db, nil
}
