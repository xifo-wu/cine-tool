package core

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open("database/data.db"), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database")
	}

	return DB
}
