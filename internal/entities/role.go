package entities

import (
	"errors"

	"github.com/StindCo/smart_ispt/pkg/id"
)

type Role struct {
	ID   id.ID  `json="id"`
	Name string `json="name"`
}

func (r Role) IsValid() (*Role, error) {
	if r.Name != "" {
		return nil, errors.New("username or password is invalid")
	}
	return &r, nil
}

func NewRole(name string) (*Role, error) {

	role := &Role{
		ID:   id.NewID(),
		Name: name,
	}
	_, err := role.IsValid()
	if err != nil {
		return nil, err
	}

	return role, nil
}
