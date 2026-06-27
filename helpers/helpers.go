package helpers

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/google/uuid"
)

func GenerateVoterID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}

func GeneratePollingStationCode() (string, error) {
	return gonanoid.New()
}