package interfaces

import "github.com/StindCo/smart_ispt/internal/entities"

type RoleService interface {
	CreateRole(name string, tag string, description string) (*entities.Role, error)
	GetRole(id string) (*entities.Role, error)
	List() ([]*entities.Role, error)
	UpdateRole(id string, entityRole *entities.Role) (*entities.Role, error)
	Delete(id string) error
	GetUsers(roleId string) ([]*entities.User, error)
}
