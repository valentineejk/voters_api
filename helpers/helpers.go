package helpers

import (
	"github.com/google/uuid"
)

func GenerateVoterID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}
