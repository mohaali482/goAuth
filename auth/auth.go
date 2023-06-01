package auth

import "time"

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
