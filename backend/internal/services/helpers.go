package services

import "github.com/google/uuid"

// toUUID converts any supported type to uuid.UUID
func toUUID(v interface{}) uuid.UUID {
	switch val := v.(type) {
	case uuid.UUID:
		return val
	case string:
		id, _ := uuid.Parse(val)
		return id
	}
	return uuid.Nil
}
