package util

import "github.com/google/uuid"

func NewNullUUID(id uuid.UUID) uuid.NullUUID {
	if id == uuid.Nil {
		return uuid.NullUUID{}
	}

	return uuid.NullUUID{
		UUID:  id,
		Valid: true,
	}
}
