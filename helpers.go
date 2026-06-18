package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

// helpers
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func generateVoterID() string {
	counter++
	return fmt.Sprintf("VTR-%05d", counter)
}

func (req *RegisterVoterRequest) validate() map[string]string {
	errs := map[string]string{}

	req.FullName = strings.TrimSpace(req.FullName)
	switch {
	case req.FullName == "":
		errs["full_name"] = "required"
	case utf8.RuneCountInString(req.FullName) < 3:
		errs["full_name"] = "minimum 3 characters"
	case utf8.RuneCountInString(req.FullName) > 100:
		errs["full_name"] = "maximum 100 characters"
	}

	req.NIN = strings.TrimSpace(req.NIN)
	if req.NIN == "" {
		errs["nin"] = "required"
	} else if !ninPattern.MatchString(req.NIN) {
		errs["nin"] = "must be exactly 11 digits"
	}

	req.DOB = strings.TrimSpace(req.DOB)
	if req.DOB == "" {
		errs["dob"] = "required"
	} else if dob, err := time.Parse("2006-01-02", req.DOB); err != nil {
		errs["dob"] = "must be YYYY-MM-DD"
	} else if time.Now().Year()-dob.Year() < 18 {
		errs["dob"] = "must be at least 18 years old"
	}

	req.State = strings.ToLower(strings.TrimSpace(req.State))
	if req.State == "" {
		errs["state"] = "required"
	} else if !validStates[req.State] {
		errs["state"] = "not a valid Nigerian state"
	}

	req.Lga = strings.ToLower(strings.TrimSpace(req.Lga))
	if req.Lga == "" {
		errs["lga"] = "required"
	}

	req.Phone = strings.TrimSpace(req.Phone)
	if req.Phone == "" {
		errs["phone"] = "required"
	} else if !phonePattern.MatchString(req.Phone) {
		errs["phone"] = "must be a valid Nigerian number"
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}
