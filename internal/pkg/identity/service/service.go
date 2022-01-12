package service

import (
	"errors"
	"fmt"

	"github.com/StindCo/smart_ispt/internal/entities"
	repository "github.com/StindCo/smart_ispt/internal/pkg/identity/Repository"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/interfaces"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewUserService(r repository.UserRepository) interfaces.UserService {
	return &UserServiceImpl{
		UserRepository: r,
	}
}

func (u UserServiceImpl) CreateUser(username string, password string) (*entities.User, error) {
	fmt.Println(username)
	_, err := u.UserRepository.GetByUsername(username)
	if err == nil {
		fmt.Println("1")
		return nil, errors.New("votre username exist déjà")
	}
	user, err := entities.NewUser(username, password)
	if err != nil {
		fmt.Println("2")
		return nil, err
	}
	return user, u.UserRepository.Create(user)
}
