package interfaces

import (
	"github.com/StindCo/smart_ispt/internal/entities"
	dto "github.com/StindCo/smart_ispt/internal/pkg/identity/Dto"
)

type UserService interface {
	CreateUser(userDTO dto.UserDTO) (*entities.User, error)
	GetUser(id string) (*entities.User, error)
	GetRoleOfUser(id string) (*entities.Role, error)
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
