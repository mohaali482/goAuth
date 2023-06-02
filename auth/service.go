package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohaali/goAuth/config"
)

type UserService struct {
	repo   Repository
	Config *config.Config
}

func NewUserService(r Repository, c *config.Config) *UserService {
	return &UserService{
		repo:   r,
		Config: c,
	}
}

func (s *UserService) Create(u User) (User, error) {
	if err := u.Validate(); err != nil {
		return User{}, err
	}
	if _, err := s.repo.GetByUsername(u.Username); err == nil {
		return User{}, ErrUsernameExists
	}
	if _, err := s.repo.GetByPhone(u.Phone); err == nil {
		return User{}, ErrPhoneExists
	}
	return s.repo.Create(u)
}

func (s *UserService) GetAll() (Users, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetByID(id int) (User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetByUsername(username string) (User, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) GetByPhone(phone string) (User, error) {
	return s.repo.GetByUsername(phone)
}

func (s *UserService) Update(id int, u User) (User, error) {
	if err := u.Validate(); err != nil {
		return User{}, err
	}
	user, _ := s.repo.GetByUsername(u.Username)
	if user.ID != id {
		return User{}, ErrUsernameExists
	}
	user, _ = s.repo.GetByPhone(u.Phone)
	if user.ID != id {
		return User{}, ErrPhoneExists
	}

	return s.repo.Update(id, u)
}

func (s *UserService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *UserService) Login(username string, password string) (User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return User{}, err
	}

	err = user.CheckPassword(password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserService) GenerateJWT(user User) (map[string]string, error) {
	accessTokenExpirationTime := time.Now().Add(time.Duration(s.Config.AccessExpTime) * time.Minute)
	refreshTokenExpirationTime := time.Now().Add(time.Duration(s.Config.RefreshExpTime) * time.Minute)
	jwtClaim := JWTClaim{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Token:    Access,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpirationTime),
		},
	}
	refreshJwtClaim := JWTClaim{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Token:    Refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpirationTime),
		},
	}

	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim).SignedString(s.Config.Secret)
	if err != nil {
		return nil, err
	}

	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshJwtClaim).SignedString(s.Config.Secret)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access":  t,
		"refresh": rt,
	}, nil
}
