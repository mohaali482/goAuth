package auth

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidUsername  = errors.New("invalid username")
	ErrInvalidPhone     = errors.New("invalid phone number")
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrUsernameExists   = errors.New("username already exists")
	ErrPhoneExists      = errors.New("phone already exists")
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username" validate:"required"`
	Phone     string    `json:"phone" validate:"required,e164"`
	Password  string    `json:"password" validate:"required"`
	Role      string    `json:"role"`
	IsAdmin   bool      `json:"is_admin"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Users []User

var (
	Access  = "access"
	Refresh = "refresh"
)

type JWTClaim struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Token    string `json:"token"`
	jwt.RegisteredClaims
}

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
	GenerateJWT(user User) (map[string]string, error)
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
