package entities

import (
	"time"
)

type User struct {
	ID        id.ID
	username  string
	password  string
	createdAt time.Time
}
