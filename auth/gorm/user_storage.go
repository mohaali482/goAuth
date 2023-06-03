package gorm

import (
	"github.com/mohaali/goAuth/auth"
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

func NewFromAuthUser(u auth.User) GormUser {
	return GormUser{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Phone:     u.Phone,
		Password:  u.Password,
		Role:      u.Role,
		IsAdmin:   u.IsAdmin,
		IsActive:  u.IsActive,
		Model: gorm.Model{
			ID:        uint(u.ID),
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: gorm.DeletedAt{
				Time: u.DeletedAt,
			},
		},
	}
}

func (u GormUser) ToEntity() auth.User {
	return auth.User{
		ID:        int(u.ID),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Phone:     u.Phone,
		Password:  u.Password,
		Role:      u.Role,
		IsAdmin:   u.IsAdmin,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt.Time,
	}
}
