package gorm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormUser struct {
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

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(dbURL string) (*GormRepository, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&GormUser{})

	return &GormRepository{db: db}, nil
}
