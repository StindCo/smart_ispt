package repository

import (
	"errors"
	"time"

	"github.com/StindCo/smart_ispt/internal/entities"
	"github.com/StindCo/smart_ispt/pkg/id"
	"gorm.io/gorm"
)

type UserGORM struct {
	gorm.Model
	ID        id.ID `gorm:"primarykey"`
	Username  string
	Password  string
	CreatedAt time.Time
}

// Set tablename (GORM)
func (UserGORM) TableName() string {
	return "users"
}

func (u UserGORM) toEntitiesUser() *entities.User {
	return &entities.User{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

func NewUserGORM(entityUser *entities.User) *UserGORM {
	u := UserGORM{}
	u.ID = entityUser.ID
	u.Username = entityUser.Username
	u.Password = entityUser.Password
	u.CreatedAt = entityUser.CreatedAt
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
		result = append(result, user.toEntitiesUser())
	}

	return result, nil
}

func (r *UserRepository) GetByUsername(username string) (*entities.User, error) {
	var user UserGORM

	r.DB.Find(&user, "username = ?", username)
	// If no such user present return an error
	if id.UUIDIsNil(user.ID) {
		return nil, errors.New("user does not exists")
	}

	return user.toEntitiesUser(), nil
}

func (r *UserRepository) Get(userID string) (*entities.User, error) {
	var user UserGORM

	r.DB.Find(&user, "id = ?", userID)

	// If no such user present return an error
	if id.UUIDIsNil(user.ID) {
		return nil, errors.New("User does not exists.")
	}

	return user.toEntitiesUser(), nil
}
