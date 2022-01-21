package service

import (
	"errors"

	"github.com/StindCo/smart_ispt/internal/entities"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/interfaces"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
)

type RoleServiceImpl struct {
	RoleRepository repository.RoleRepository
	UserRepository repository.UserRepository
}

func NewRoleService(r repository.RoleRepository, u repository.UserRepository) interfaces.RoleService {
	return &RoleServiceImpl{
		RoleRepository: r,
		UserRepository: u,
	}
}

func (r RoleServiceImpl) CreateRole(name string, tag string, description string) (*entities.Role, error) {
	_, err := r.RoleRepository.GetByTag(tag)
	if err == nil {
		return nil, errors.New("il existe déjà un role avec ce tag")
	}
	role, err := entities.NewRole(name, tag, description)
	if err != nil {
		return nil, err
	}
	role.Users, _ = r.UserRepository.GetUsersByRoleID(role.ID.String())
	return role, r.RoleRepository.Create(role)
}

func (r RoleServiceImpl) GetRole(id string) (*entities.Role, error) {
	role, err := r.RoleRepository.Get(id)
	if err != nil {
		return nil, errors.New("ce role n'existe pas")
	}
	role.Users, _ = r.UserRepository.GetUsersByRoleID(role.ID.String())
	return role, nil
}

func (r RoleServiceImpl) List() ([]*entities.Role, error) {
	roles, err := r.RoleRepository.List()
	if err != nil {
		return nil, errors.New("désolé, il y' a erreur")
	}
	var rolesResult []*entities.Role
	for _, role := range roles {
		role.Users, _ = r.UserRepository.GetUsersByRoleID(role.ID.String())
		rolesResult = append(rolesResult, role)
	}

	return rolesResult, nil
}

func (r RoleServiceImpl) UpdateRole(id string, entityRole *entities.Role) (*entities.Role, error) {
	role, err := r.RoleRepository.Get(id)
	if err != nil {
		return nil, err
	}
	if entityRole.Name != "" {
		role.Name = entityRole.Name
	}
	if entityRole.Name != "" {
		role.Description = entityRole.Description
	}

	err = r.RoleRepository.DB.Save(role).Error
	return role, err
}

func (r RoleServiceImpl) Delete(id string) error {
	_, err := r.RoleRepository.Get(id)
	if err != nil {
		return err
	}
	if (r.RoleRepository.Delete(id)) != nil {
		return errors.New("erreur lors de la suppression de ce role, peut-être qu'il n'existe pas")
	}
	return nil
}

func (r RoleServiceImpl) GetUsers(roleId string) ([]*entities.User, error) {
	_, err := r.RoleRepository.Get(roleId)
	if err != nil {
		return nil, errors.New("ce role n'existe pas")
	}
	return r.UserRepository.GetUsersByRoleID(roleId)
}
