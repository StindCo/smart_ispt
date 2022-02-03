package repository

import (
	"errors"
	"time"

	"github.com/StindCo/smart_ispt/internal/entities"
	"gorm.io/gorm"
)

type UserGORM struct {
	gorm.Model
	ID            string `gorm:"primarykey"`
	Username      string
	ApplicationID string
	Fullname      string
	Password      string
	CreatedAt     time.Time
	RoleID        string `gorm:"size:60"`
	IsAdmin       int
	IsDevelopper  int
}

// Set tablename (GORM)
func (UserGORM) TableName() string {
	return "users"
}

func (u UserGORM) ToEntitiesUser() *entities.User {
	return &entities.User{
		ID:           u.ID,
		Fullname:     u.Fullname,
		Username:     u.Username,
		Password:     u.Password,
		CreatedAt:    u.CreatedAt,
		RoleID:       u.RoleID,
		IsAdmin:      u.IsAdmin,
		IsDevelopper: u.IsDevelopper,
	}
}

func NewUserGORM(entityUser *entities.User) *UserGORM {
	u := UserGORM{}
	u.ID = entityUser.ID
	u.Fullname = entityUser.Fullname
	u.Username = entityUser.Username
	u.Password = entityUser.Password
	u.CreatedAt = entityUser.CreatedAt
	u.RoleID = entityUser.RoleID
	u.IsAdmin = entityUser.IsAdmin
	u.IsDevelopper = entityUser.IsDevelopper
	return &u
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserGORMRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(entityUser *entities.User) error {
	u := NewUserGORM(entityUser)

	err := r.DB.Create(&u).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) List() ([]*entities.User, error) {
	var users []UserGORM

	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}

	// TODO: Refactor. maybe inefficient.
	result := make([]*entities.User, 0, len(users))
	for _, user := range users {
		result = append(result, user.ToEntitiesUser())
	}

	return result, nil
}

func (r *UserRepository) GetUsersByRoleID(roleId string) ([]*entities.User, error) {
	var users []UserGORM

	err := r.DB.Find(&users, "role_id = ?", roleId).Error
	if err != nil {
		return nil, err
	}

	// TODO: Refactor. maybe inefficient.
	result := make([]*entities.User, 0, len(users))
	for _, user := range users {
		result = append(result, user.ToEntitiesUser())
	}

	return result, nil
}

func (r *UserRepository) GetAdminsUsers() ([]*entities.User, error) {
	var users []UserGORM

	err := r.DB.Find(&users, "is_admin = ?", 1).Error
	if err != nil {
		return nil, err
	}

	// TODO: Refactor. maybe inefficient.
	result := make([]*entities.User, 0, len(users))
	for _, user := range users {
		result = append(result, user.ToEntitiesUser())
	}

	return result, nil
}

func (r *UserRepository) GetDeveloppersUsers() ([]*entities.User, error) {
	var users []UserGORM

	err := r.DB.Find(&users, "is_developper = ?", 1).Error
	if err != nil {
		return nil, err
	}

	// TODO: Refactor. maybe inefficient.
	result := make([]*entities.User, 0, len(users))
	for _, user := range users {
		result = append(result, user.ToEntitiesUser())
	}

	return result, nil
}

func (r *UserRepository) GetByUsername(username string) (*entities.User, error) {
	var user UserGORM

	r.DB.Find(&user, "username = ?", username)
	// If no such user present return an error
	if user.ID == "" {
		return nil, errors.New("user does not exists")
	}

	return user.ToEntitiesUser(), nil
}

func (r *UserRepository) GetByUsernameAndPassword(username string, password string) (*entities.User, error) {
	var user UserGORM

	r.DB.Find(&user, "username = ?", username, "password = ?", password)
	// If no such user present return an error
	if user.ID == "" {
		return nil, errors.New("user does not exists")
	}

	return user.ToEntitiesUser(), nil
}

func (r *UserRepository) Get(userID string) (*entities.User, error) {
	var user UserGORM

	r.DB.Find(&user, "id = ?", userID)

	// If no such user present return an error
	if user.ID == "" {
		return nil, errors.New("user does not exists")
	}

	return user.ToEntitiesUser(), nil
}

func (r *UserRepository) Update(userID string, entityUser *entities.User) (*entities.User, error) {
	var user UserGORM

	r.DB.Find(&user, "id = ?", userID)
	user.Password = entityUser.Password
	r.DB.Save(&user)
	if user.ID == "" {
		return nil, errors.New("error à la modification")
	}
	return user.ToEntitiesUser(), nil
}

func (r *UserRepository) UpdateRole(userID, roleId string) (*entities.User, error) {
	var user UserGORM
	r.DB.Find(&user, "id = ?", userID)
	user.RoleID = roleId
	r.DB.Save(&user)
	if user.ID == "" {
		return nil, errors.New("error à la modification")
	}
	return user.ToEntitiesUser(), nil
}

func (r *UserRepository) UpdateAdmin(userID string, value int) (*entities.User, error) {
	var user UserGORM
	r.DB.Find(&user, "id = ?", userID)
	user.IsAdmin = value
	r.DB.Save(&user)
	if user.ID == "" {
		return nil, errors.New("error à la modification")
	}
	return user.ToEntitiesUser(), nil
}

func (r *UserRepository) UpdateDevelopper(userID string, value int) (*entities.User, error) {
	var user UserGORM
	r.DB.Find(&user, "id = ?", userID)
	user.IsDevelopper = value
	r.DB.Save(&user)
	if user.ID == "" {
		return nil, errors.New("error à la modification")
	}
	return user.ToEntitiesUser(), nil
}

// TODO: Une fonction à créer
func (r *UserRepository) Delete(userId string) error {
	return r.DB.Where("id = ?", userId).Delete(&UserGORM{}).Error
}
