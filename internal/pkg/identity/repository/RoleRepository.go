package repository

import (
	"errors"
	"time"

	"github.com/StindCo/smart_ispt/internal/entities"
	"gorm.io/gorm"
)

type RoleGORM struct {
	gorm.Model
	ID            string `gorm:"primarykey"`
	Name          string
	ApplicationID string
	Description   string
	UserID        string
	Tag           string
	CreatedAt     time.Time
}

// Set tablename (GORM)
func (RoleGORM) TableName() string {
	return "roles"
}

func (r RoleGORM) ToEntitiesRole() *entities.Role {
	return &entities.Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Tag:         r.Tag,
		CreatedAt:   r.CreatedAt,
	}
}

func NewRoleGORM(entityRole *entities.Role) *RoleGORM {
	r := RoleGORM{}

	r.ID = entityRole.ID
	r.Name = entityRole.Name
	r.Description = entityRole.Description
	r.CreatedAt = entityRole.CreatedAt
	r.Tag = entityRole.Tag
	return &r
}

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleGORMRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		DB: db,
	}
}

func (r *RoleRepository) Create(entityRole *entities.Role) error {
	role := NewRoleGORM(entityRole)

	err := r.DB.Create(&role).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *RoleRepository) List() ([]*entities.Role, error) {
	var roles []RoleGORM

	err := r.DB.Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// TODO: Refactor. maybe inefficient.
	result := make([]*entities.Role, 0, len(roles))
	for _, user := range roles {
		result = append(result, user.ToEntitiesRole())
	}

	return result, nil
}

func (r *RoleRepository) GetByTag(tag string) (*entities.Role, error) {
	var role RoleGORM

	r.DB.Find(&role, "tag = ?", tag)
	// If no such user present return an error
	if role.ID == "" {
		return nil, errors.New("user does not exists")
	}

	return role.ToEntitiesRole(), nil
}

func (r *RoleRepository) Get(roleID string) (*entities.Role, error) {
	var role RoleGORM

	r.DB.Find(&role, "id = ?", roleID)

	// If no such user present return an error
	if role.ID == "" {
		return nil, errors.New("user does not exists")
	}

	return role.ToEntitiesRole(), nil
}

func (r *RoleRepository) Update(roleID string, entityRole *entities.Role) (*entities.Role, error) {
	role, err := r.Get(roleID)
	if err != nil {
		return nil, err
	}

	if role.ID == "" {
		return nil, errors.New("error à la modification")
	}
	return entityRole, nil
}

func (r *RoleRepository) Delete(roleId string) error {
	return r.DB.Delete(&RoleGORM{}, roleId).Error
}
