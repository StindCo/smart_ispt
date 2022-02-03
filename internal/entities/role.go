package entities

import (
	"time"

	"github.com/StindCo/smart_ispt/pkg/id"
)

type Role struct {
	ID          id.ID  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	CreatedAt   time.Time
	Users       []*User
}

func NewRole(name string, tag string, description string) (*Role, error) {

	role := &Role{
		ID:          id.NewID(),
		Name:        name,
		Description: description,
		Tag:         tag,
		CreatedAt:   time.Now(),
	}

	return role, nil
}
