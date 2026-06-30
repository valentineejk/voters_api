package helpers

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateVoterID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}

var validNigerianStates = map[string]bool{
	"abia": true, "adamawa": true, "akwa ibom": true, "anambra": true,
	"bauchi": true, "bayelsa": true, "benue": true, "borno": true,
	"cross river": true, "delta": true, "ebonyi": true, "edo": true,
	"ekiti": true, "enugu": true, "fct": true, "gombe": true,
	"imo": true, "jigawa": true, "kaduna": true, "kano": true,
	"katsina": true, "kebbi": true, "kogi": true, "kwara": true,
	"lagos": true, "nasarawa": true, "niger": true, "ogun": true,
	"ondo": true, "osun": true, "oyo": true, "plateau": true,
	"rivers": true, "sokoto": true, "taraba": true, "yobe": true,
	"zamfara": true,
}

func ValidateState(state string) bool {
	return validNigerianStates[strings.ToLower(strings.TrimSpace(state))]
}
