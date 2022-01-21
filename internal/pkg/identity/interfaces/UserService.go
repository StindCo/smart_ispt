package interfaces

import "github.com/StindCo/smart_ispt/internal/entities"

type UserService interface {
	CreateUser(username string, password string, fullname string) (*entities.User, error)
	GetUser(id string) (*entities.User, error)
	List() ([]*entities.User, error)
	GetUsersWhoAreAdmin() ([]*entities.User, error)
	GetUsersWhoAreDevelopper() ([]*entities.User, error)
	SetRole(userId, roleId string) (*entities.User, error)
	SetAdminPermission(userId string) (*entities.User, error)
	SetDevelopperPermission(userId string) (*entities.User, error)
	RemoveAdminPermission(userId string) (*entities.User, error)
	RemoveDevelopperPermission(userId string) (*entities.User, error)
	UpdatePassword(id string, oldPassword string, newPassword string) (*entities.User, error)
	Delete(id string) error
}
