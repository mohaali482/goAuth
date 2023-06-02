package auth

type UserService struct {
	repo Repository
}

func NewUserService(r Repository) *UserService {
	return &UserService{
		repo: r,
	}
}
