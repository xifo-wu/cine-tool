package api

import "gorm.io/gorm"

type Api struct {
	DB *gorm.DB
}
