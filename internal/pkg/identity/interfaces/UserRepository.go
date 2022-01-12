package interfaces

import "github.com/StindCo/smart_ispt/internal/entities"

type UserRepositoryInterface interface {
	Create(entityUser *entities.User) error
	// Delete(userID string) error
	// Update(userId string, entityUser *entities.User) (*entities.User, error)

	List() ([]*entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	Get(userID string) (*entities.User, error)
}
