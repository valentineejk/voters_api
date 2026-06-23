package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mu       sync.RWMutex
	store    = map[string]Voter{}
	ninIndex = map[string]string{}
	counter  int
)

var (
	ninPattern   = regexp.MustCompile(`^\d{11}$`)
	phonePattern = regexp.MustCompile(`^(\+234|0)[7-9][0-1]\d{8}$`)
	voterIDPat   = regexp.MustCompile(`^VTR-\d{5}$`)

	validStatuses = map[string]bool{
		"pending": true, "verified": true, "rejected": true,
	}
	validStates = map[string]bool{
		"lagos": true, "abuja": true, "kano": true,
		"rivers": true, "oyo": true, // add all 36 states
	}
)

// handlers
func get_all_voters(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	state := strings.ToLower(q.Get("state"))
	status := strings.ToLower(q.Get("status"))

	if status != "" && !validStatuses[status] {
		writeError(w, http.StatusBadRequest,
			"status must be: pending, verified or rejected")
		return
	}

	page, limit := 1, 50
	if p := q.Get("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}
	if l := q.Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	mu.RLock()
	var filtered []Voter
	for _, v := range store {
		if state != "" && v.State != state {
			continue
		}
		if status != "" && v.Status != status {
			continue
		}
		filtered = append(filtered, v)
	}
	mu.RUnlock()

	total := len(filtered)
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": filtered[start:end],
		"meta": PaginatedMeta{
			Page: page, Limit: limit, Total: total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	})

}

func get_one_voter(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id") //id from param

	if !voterIDPat.MatchString(id) {
		writeError(w, http.StatusBadRequest, "invalid voter id format")
		return
	}

	mu.RLock()
	voter, ok := store[id]
	mu.RUnlock()

	if !ok {
		writeJSON(w, http.StatusNotFound, "voters with this id not found")
		return
	}

	writeJSON(w, http.StatusOK, voter)

}

func register_voter(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req RegisterVoterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnsupportedMediaType, "Content type must be app/json")
		return
	}

	////////////////////////////////////////////////////////
	if errs := req.validate(); errs != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "validation failed", "fields": errs,
		})
		return
	}

	////////////////////////////////////////////////////////
	mu.Lock()
	defer mu.Unlock()
	if _, exists := ninIndex[req.NIN]; exists {
		writeError(w, http.StatusConflict, "NIN already registered")
		return
	}
	////////////////////////////////////////////////////////
	vo := Voter{
		ID:       generateVoterID(),
		FullName: req.FullName,
		NIN:      req.NIN,
		DOB:      req.DOB,
		State:    req.State,
		Lga:      req.Lga,
		Phone:    req.Phone,
		Status:   "pending",
	}

	store[vo.ID] = vo
	ninIndex[req.NIN] = vo.ID

	writeJSON(w, http.StatusCreated, vo)

}

func delete_voter(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	//voters id check
	if !voterIDPat.MatchString(id) {
		writeError(w, http.StatusBadRequest, "invalid voter id format")
		return
	}
	mu.Lock()
	defer mu.Unlock()
	v, ok := store[id]

	if !ok {
		writeError(w, http.StatusNotFound, "voters not found")
		return
	}

	delete(store, id)
	delete(ninIndex, v.NIN)
	w.WriteHeader(http.StatusNoContent)

}

func update_voter_status(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))
	if req.Status == "" {
		writeError(w, http.StatusUnprocessableEntity, "status is required")
		return
	}
	if !validStatuses[req.Status] {
		writeError(w, http.StatusUnprocessableEntity,
			"status must be: pending, verified or rejected")
		return
	}

	mu.Lock()
	defer mu.Unlock()
	vo, ok := store[id]
	if !ok {
		writeError(w, http.StatusNotFound, "voter not found")
	}

	vo.Status = req.Status
	store[id] = vo
	writeJSON(w, http.StatusOK, vo)

}

// healthcheck
func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "service is running fine",
	})
}

func main() {

	PORT := ":8000"
	// mux := http.NewServeMux()

	//health
	// mux.HandleFunc("GET /api/v1/health", healthCheck)

	// //voter routes
	// mux.HandleFunc("GET /api/v1/voters", get_all_voters)

	// mux.HandleFunc("POST /api/v1/voters", register_voter)

	// mux.HandleFunc("GET /api/v1/voters/{id}", get_one_voter)

	// mux.HandleFunc("DELETE /api/v1/voters/{id}", delete_voter)

	// mux.HandleFunc("PUT /api/v1/voters/{id}/status", update_voter_status)

	//server
	fmt.Println("server started")
	// http.ListenAndServe(PORT, mux)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.GET("/health")
	r.Run(PORT)

}
