package repository

import "gorm.io/gorm"

type cartDatabase struct {
	DB *gorm.DB
}
