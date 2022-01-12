package id

import "github.com/google/uuid"

//ID entity ID
type ID = uuid.UUID

//NewID create a new entity ID
func NewID() ID {
	return ID(uuid.New())
}

func UUIDIsNil(id uuid.UUID) bool {
	return id == uuid.Nil
}
