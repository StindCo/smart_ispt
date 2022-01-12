package entities

import (
	"errors"
	"time"

	"github.com/StindCo/smart_ispt/pkg/id"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        id.ID     `json="id"`
	Username  string    `json="username"`
	Password  string    `json="password"`
	CreatedAt time.Time `json="createdAt"`
	Role      Role
}

func (u User) IsValid() (*User, error) {
	if u.Username != "" && u.Password != "" {
		return nil, errors.New("username or password is invalid")
	}
	return &u, nil
}

func (u User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func newPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func NewUser(username string, password string) (*User, error) {
	pwd, err := newPasswordHash(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        id.NewID(),
		Username:  username,
		Password:  pwd,
		CreatedAt: time.Now(),
	}
	_, err = user.IsValid()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u User) SetRole(role Role) {
	u.Role = role
}
