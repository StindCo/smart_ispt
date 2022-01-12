package interfaces

import "github.com/StindCo/smart_ispt/internal/entities"

type UserService interface {
	CreateUser(username string, password string) (*entities.User, error)
}
