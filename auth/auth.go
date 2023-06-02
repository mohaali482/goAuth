package auth

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidUsername  = errors.New("invalid username")
	ErrInvalidPhone     = errors.New("invalid phone number")
	ErrWrongCredentials = errors.New("wrong credentials")
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Username  string `json:"username" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Password  string `json:"password" validate:"required"`
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

func (u *User) Validate() error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate.Struct(u)
}

type UseCase interface {
	Create(user User) (User, error)
	GetAll() (Users, error)
	GetByID(id int) (User, error)
	GetByUsername(username string) (User, error)
	GetByPhone(phone string) (User, error)
	Update(id int, user User) (User, error)
	Delete(id int) error
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
	Delete(id int) error
}
