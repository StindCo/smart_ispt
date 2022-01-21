package id

import "github.com/google/uuid"

//ID entity ID
type ID = uuid.UUID

//NewID create a new entity ID
func NewID() ID {
	return ID(uuid.New())
}

func StringToUuid(param string) (uuid.UUID, error) {
	return uuid.Parse(NewID().String())
}
func UUIDIsNil(id uuid.UUID) bool {
	return id == uuid.Nil
}
