package services

import (
	"github.com/DieegoAlves/CrypexGoAPI/src/repositories"
)

type UserService struct {
	repository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return UserService{
		repository: userRepository,
	}
}
