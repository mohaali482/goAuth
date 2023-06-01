package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Username  string
	Phone     string
	Password  string
	Role      string
	IsAdmin   bool
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Users []User

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type UseCase interface {
	Create(user User) (User, error)
	GetAll() (Users, error)
	GetByID(id int) (User, error)
	GetByUsername(username string) (User, error)
	GetByPhone(phone string) (User, error)
	Update(id int, user User) (User, error)
	Delete(id int) (User, error)
	Login(username string, password string) (User, error)
	Logout() error
}

type Repository interface {
	Create(user User) (User, error)
	GetAll() (Users, error)
	GetByID(id int) (User, error)
	GetByUsername(username string) (User, error)
	GetByPhone(phone string) (User, error)
	Update(id int, user User) (User, error)
	Delete(id int) (User, error)
}
