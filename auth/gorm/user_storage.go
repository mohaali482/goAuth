package gorm

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Username  string `gorm:"index"`
	Phone     string
	Password  string
	Role      string
	IsAdmin   bool `gorm:"default:false"`
	IsActive  bool `gorm:"default:true index"`
}
