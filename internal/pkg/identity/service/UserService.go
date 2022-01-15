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
	RoleRepository repository.RoleRepository
}

func NewUserService(r repository.UserRepository, rr repository.RoleRepository) interfaces.UserService {
	return &UserServiceImpl{
		UserRepository: r,
		RoleRepository: rr,
	}
}

func (u UserServiceImpl) CreateUser(username string, password string) (*entities.User, error) {
	_, err := u.UserRepository.GetByUsername(username)
	if err == nil {
		return nil, errors.New("votre username exist déjà")
	}
	user, err := entities.NewUser(username, password)
	if err != nil {
		return nil, err
	}
	user.Role, _ = u.RoleRepository.Get(user.RoleID)
	return user, u.UserRepository.Create(user)
}

func (u UserServiceImpl) GetUser(id string) (*entities.User, error) {
	user, err := u.UserRepository.Get(id)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas")
	}
	user.Role, _ = u.RoleRepository.Get(user.RoleID)
	return user, nil
}

func (u UserServiceImpl) List() ([]*entities.User, error) {
	users, err := u.UserRepository.List()
	if err != nil {
		return nil, errors.New("désolé, il y' a erreur")
	}
	var usersResult []*entities.User
	for _, user := range users {
		user.Role, _ = u.RoleRepository.Get(user.RoleID)
		usersResult = append(usersResult, user)
	}
	return usersResult, nil
}

func (u UserServiceImpl) UpdatePassword(id string, oldPassword string, newPassword string) (*entities.User, error) {
	user, err := u.UserRepository.Get(id)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas")
	}
	if (user.ValidatePassword(oldPassword)) != nil {
		fmt.Println(user.ValidatePassword(oldPassword))
		return nil, errors.New("mot de passe incorrect")
	}
	if (user.NewPassword(newPassword)) != nil {
		return nil, errors.New("erreur interne, impossible d'hasher le mot de passe")
	}

	return u.UserRepository.Update(id, user)
}

func (u UserServiceImpl) SetRole(userID, roleId string) (*entities.User, error) {
	_, err := u.UserRepository.Get(userID)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas")
	}
	_, err = u.RoleRepository.Get(roleId)
	if err != nil {
		return nil, errors.New("ce role n'existe pas")
	}

	user, err := u.UserRepository.UpdateRole(userID, roleId)
	user.Role, _ = u.RoleRepository.Get(user.RoleID)
	return user, err
}

func (u UserServiceImpl) Delete(id string) error {
	_, err := u.UserRepository.Get(id)
	if err != nil {
		return err	
	}
	if (u.UserRepository.Delete(id)) != nil {
		return errors.New("erreur lors de la suppression de cet utilisateur, peut-être qu'il n'existe pas")
	}
	return nil
}
