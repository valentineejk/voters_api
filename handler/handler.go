package handler

import (
	"errors"
	"math"
	"net/http"
	"regexp"
	"strconv"
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
	validStatus = map[string]bool{
		"pending": true,
		"verified": true,
		"rejected": true,
	}
	validStates = map[string]bool{
		"abia": true, "adamawa": true, "akwa ibom": true, "anambra": true,
		"bauchi": true, "bayelsa": true, "benue": true, "borno": true,
		"cross river": true, "delta": true, "ebonyi": true, "edo": true,
		"ekiti": true, "enugu": true, "gombe": true, "imo": true,
		"jigawa": true, "kaduna": true, "kano": true, "katsina": true,
		"kebbi": true, "kogi": true, "kwara": true, "lagos": true,
		"nasarawa": true, "niger": true, "ogun": true, "ondo": true,
		"osun": true, "oyo": true, "plateau": true, "rivers": true,
		"sokoto": true, "taraba": true, "yobe": true, "zamfara": true,
	}
)

type Handler struct {
	queries *dbq.Queries
}

func New(queries *dbq.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

// handlers
func (h *Handler) GetAllVoters(c *gin.Context) {
	status := c.Param("status")
	state := c.Param("state")

	if status != "" && !validStatus[status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status must be pending, verified or rejected",
		})
		return
	}

	if state != "" && !validStates[state] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "state must be a valid nigerian state",
		})
		return
	}
	var statusAddress *string
	if status != "" {
		statusAddress = &status
	}

	var stateAddress *string
	if state != "" {
		stateAddress = &state
	}

	page, limit := 1, 50
	if p := c.Param("page"); p != "" {
		pageInt, err := strconv.Atoi(p)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "page must be a positive number",
			})
			return
		}
		page = pageInt
	}

	if l := c.Param("limit"); l != "" {
		limitInt, err := strconv.Atoi(l)
		if err != nil || limitInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "limit must be a positive number",
			})
			return
		}
		limit = limitInt
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


	params := dbq.ListVotersParams{
		Limit: int32(limit),
		Offset: int32((page - 1) * limit),
		State: stateAddress,
		Status: statusAddress,
	}
	
	voters, err := h.queries.ListVoters(c, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "voter not found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	total := len(voters)
	totalPages := int(math.Ceil(float64(total/limit)))
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, map[string]any{
		"data": voters,
		"meta": model.PaginatedMeta{
			Page: page,
			Limit: limit,
			Total: total,
			TotalPages:	totalPages,
			HasNext: page < totalPages,
			HasPrev: page > 1,
		},
	})
	
}

func (h *Handler) Get_one_voter(c *gin.Context) {

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

	//TODO: add another check

	c.JSON(http.StatusOK, voter)

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

	//state validation
	if !helpers.ValidateState(req.State) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid state",
			"message": "Well, the state must be a valid Nigerian state",
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
