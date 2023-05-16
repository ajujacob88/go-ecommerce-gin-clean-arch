package domain

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Phone     string `json:"phone_no" gorm:"uniqueIndex;not null"`
	Password  string `json:"password" gorm:"not null"`
}
