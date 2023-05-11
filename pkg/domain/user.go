package domain

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	// ID          uint      `json:"id" gorm:"primaryKey;unique"`  //gorm.Model is used instead of id, created at, deleted at
	FirstName   string `json:"first_name" gorm:"not null" binding:"required,min=3,max=18"`
	LastName    string `json:"last_name" binding:"required,max=15"`
	Email       string `json:"email" gorm:"unique,not null" binding:"required,email"`
	Phone       string `json:"phone_no" gorm:"unique" binding:"required,min=10,max=10"`
	Password    string `json:"password" gorm:"not null" binding:"required"`
	BlockStatus bool   `json:"block_status" gorm:"not null;default:false"`
	// CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	// UpdatedAt   time.Time `json:"updated_at"`
}
