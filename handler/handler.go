package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbq "github.com/valentineejk/voters_api/database/sqlc"
	"github.com/valentineejk/voters_api/helpers"
	model "github.com/valentineejk/voters_api/modal"
)

var (
	voterIDPat = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// // handlers
// func get_all_voters(w http.ResponseWriter, r *http.Request) {

// 	q := r.URL.Query()
// 	state := strings.ToLower(q.Get("state"))
// 	status := strings.ToLower(q.Get("status"))

// 	if status != "" && !validStatuses[status] {
// 		writeError(w, http.StatusBadRequest,
// 			"status must be: pending, verified or rejected")
// 		return
// 	}

// 	page, limit := 1, 50
// 	if p := q.Get("page"); p != "" {
// 		page, _ = strconv.Atoi(p)
// 	}
// 	if l := q.Get("limit"); l != "" {
// 		limit, _ = strconv.Atoi(l)
// 	}
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 50
// 	}
// 	if limit > 100 {
// 		limit = 100
// 	}

// 	mu.RLock()
// 	var filtered []Voter
// 	for _, v := range store {
// 		if state != "" && v.State != state {
// 			continue
// 		}
// 		if status != "" && v.Status != status {
// 			continue
// 		}
// 		filtered = append(filtered, v)
// 	}
// 	mu.RUnlock()

// 	total := len(filtered)
// 	totalPages := int(math.Ceil(float64(total) / float64(limit)))
// 	start := (page - 1) * limit
// 	end := start + limit
// 	if start > total {
// 		start = total
// 	}
// 	if end > total {
// 		end = total
// 	}

// 	writeJSON(w, http.StatusOK, map[string]any{
// 		"data": filtered[start:end],
// 		"meta": PaginatedMeta{
// 			Page: page, Limit: limit, Total: total,
// 			TotalPages: totalPages,
// 			HasNext:    page < totalPages,
// 			HasPrev:    page > 1,
// 		},
// 	})

// }

type Handler struct {
	queries *dbq.Queries
}

func New(queries *dbq.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

func (h *Handler) Get_one_voter(c *gin.Context) {

	// when voter_id query param is present, skip path-based id extraction
	if voterID := c.Query("voter_id"); voterID != "" {
		voter, err := h.queries.GetVoterByVoterID(c.Request.Context(), voterID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "voter not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}
		c.JSON(http.StatusOK, voter)
		return
	}

	id := c.Param("id")

	if !voterIDPat.MatchString(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid voter id",
		})
		return
	}

	voter, err := h.queries.GetVoter(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "voter not found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "something went wrong",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, voter)

}

func (h *Handler) Delete_voter(c *gin.Context) {

	id := c.Param("id")

	if !voterIDPat.MatchString(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid voter id"})
		return
	}

	err := h.queries.DeleteVoter(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "voter not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.Status(http.StatusNoContent)

}

func (h *Handler) Register_voter(c *gin.Context) {

	var req model.RegisterVoterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid body data",
		})
		return
	}

	//nin check
	_, err := h.queries.GetVoterByNIN(c.Request.Context(), req.NIN)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   err,
			"message": "Nin taken",
		})
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}

	//prse dob
	dob, _ := time.Parse("2006-01-02", req.DOB)

	voterReq := dbq.CreateVoterParams{
		ID:       helpers.GenerateVoterID(),
		FullName: req.FullName,
		Nin:      req.NIN,
		Dob:      pgtype.Date{Time: dob, Valid: true},
		State:    req.State,
		Lga:      req.Lga,
		Phone:    req.Phone,
	}

	voter, err := h.queries.CreateVoter(c.Request.Context(), voterReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "failed to create voter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    voter,
		"message": "voter created succesfully",
	})

}

// func delete_voter(w http.ResponseWriter, r *http.Request) {

// 	id := r.PathValue("id")
// 	//voters id check
// 	if !voterIDPat.MatchString(id) {
// 		writeError(w, http.StatusBadRequest, "invalid voter id format")
// 		return
// 	}
// 	mu.Lock()
// 	defer mu.Unlock()
// 	v, ok := store[id]

// 	if !ok {
// 		writeError(w, http.StatusNotFound, "voters not found")
// 		return
// 	}

// 	delete(store, id)
// 	delete(ninIndex, v.NIN)
// 	w.WriteHeader(http.StatusNoContent)

// }

func (h *Handler) Update_voter_status(c *gin.Context) {

	id := c.Param("id")

	if !voterIDPat.MatchString(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid voter id",
		})
		return
	}

	var req model.UpdateVoterStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid body data",
		})
		return
	}

	req.Status = strings.TrimSpace(strings.ToLower(req.Status))
	if req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "status is required",
			"message": "invalid body data",
		})
		return
	}
	voter, err := h.queries.UpdateVoterStatus(c.Request.Context(), dbq.UpdateVoterStatusParams{
		ID:     id,
		Status: req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "failed to update voter status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    voter,
		"message": "voter status updated succesfully",
	})

}

// healthcheck
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "service is running fine",
	})
}
