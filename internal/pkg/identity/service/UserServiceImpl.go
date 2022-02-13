package service

import (
	"errors"
	"fmt"

	"github.com/StindCo/smart_ispt/internal/entities"
	dto "github.com/StindCo/smart_ispt/internal/pkg/identity/Dto"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/interfaces"
	repository "github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
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

func (u UserServiceImpl) CreateUser(userDTO dto.UserDTO) (*entities.User, error) {
	_, err := u.UserRepository.GetByUsername(userDTO.Username)
	if err == nil {
		return nil, errors.New("votre username exist déjà")
	}
	user, err := entities.NewUser(userDTO.Username, userDTO.Password)
	if err != nil {
		return nil, err
	}
	if userDTO.RoleTag != "" {
		user.Role, err = u.RoleRepository.GetByTag(userDTO.RoleTag)
		if err == nil {
			user.RoleID = user.Role.ID
		}
	}
	user.Status = 3
	user.Fullname = userDTO.Fullname
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

func (u UserServiceImpl) GetRoleOfUser(userId string) (*entities.Role, error) {
	user, err := u.UserRepository.Get(userId)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas")
	}
	user.Role, _ = u.RoleRepository.Get(user.RoleID)
	return user.Role, nil
}

func (u UserServiceImpl) GetUsersWhoAreAdmin() ([]*entities.User, error) {
	return u.UserRepository.GetAdminsUsers()
}

func (u UserServiceImpl) GetUsersWhoAreDevelopper() ([]*entities.User, error) {
	return u.UserRepository.GetDeveloppersUsers()
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

func (u UserServiceImpl) SetAdminPermission(id string) (*entities.User, error) {
	_, err := u.UserRepository.Get(id)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas")
	}
	return u.UserRepository.UpdateAdmin(id, 1)
}

func (u UserServiceImpl) SetDevelopperPermission(id string) (*entities.User, error) {
	_, err := u.UserRepository.Get(id)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas")
	}
	return u.UserRepository.UpdateDevelopper(id, 1)
}

func (u UserServiceImpl) RemoveAdminPermission(id string) (*entities.User, error) {
	_, err := u.UserRepository.Get(id)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas")
	}
	return u.UserRepository.UpdateAdmin(id, 0)
}

func (u UserServiceImpl) RemoveDevelopperPermission(id string) (*entities.User, error) {
	_, err := u.UserRepository.Get(id)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas")
	}
	return u.UserRepository.UpdateDevelopper(id, 0)
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
