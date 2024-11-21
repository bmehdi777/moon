package database

import (
	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDBConn() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.GlobalConfig.DatabasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
