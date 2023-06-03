package gorm

import (
	"github.com/mohaali482/goAuth/auth"
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
	IsActive  bool `gorm:"index,default:true"`
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

func (r *GormRepository) Create(u auth.User) (auth.User, error) {
	user := NewFromAuthUser(u)
	err := r.db.Create(&user).Error
	if err != nil {
		return auth.User{}, err
	}
	return user.ToEntity(), nil
}

func (r *GormRepository) GetAll() (auth.Users, error) {
	var users []GormUser
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	var usersEntity auth.Users
	for _, u := range users {
		usersEntity = append(usersEntity, u.ToEntity())
	}
	return usersEntity, nil
}

func (r *GormRepository) GetByID(id int) (auth.User, error) {
	var user GormUser
	err := r.db.First(&user, id).Error
	if err != nil {
		return auth.User{}, err
	}
	return user.ToEntity(), nil
}

func (r *GormRepository) GetByUsername(username string) (auth.User, error) {
	var user GormUser
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return auth.User{}, err
	}
	return user.ToEntity(), nil
}

func (r *GormRepository) GetByPhone(phone string) (auth.User, error) {
	var user GormUser
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return auth.User{}, err
	}
	return user.ToEntity(), nil
}

func (r *GormRepository) Update(id int, u auth.User) (auth.User, error) {
	user := NewFromAuthUser(u)
	err := r.db.Model(&user).Where("id = ?", id).Updates(&user).Error
	if err != nil {
		return auth.User{}, err
	}
	return user.ToEntity(), nil
}

func (r *GormRepository) Delete(id int) error {
	var user GormUser
	err := r.db.Where("id = ?", id).Delete(&user).Error
	return err
}
