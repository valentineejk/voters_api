package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbq "github.com/valentineejk/voters_api/database/sqlc"
	"github.com/valentineejk/voters_api/helpers"
	model "github.com/valentineejk/voters_api/modal"
)

// register
func (h *Handler) RegisterHandler(c *gin.Context) {

	var req model.RegisterUserRequest

	//check req from from frontend
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//hash before storage
	hash, err := helpers.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to hash password",
		})
		return
	}

	user, err := h.queries.CreateUser(c.Request.Context(), dbq.CreateUserParams{
		ID:           helpers.GenerateVoterID(),
		Email:        req.Email,
		PasswordHash: hash,
	})

	//takehome - check if email exists
	//generate token for user

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error registrering user",
		})
		return
	}

	userResponse := model.RegisterUserResponse{
		Email: user.Email,
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":   user.ID,
		"data": userResponse,
	})
}

// login
func (h *Handler) Login(c *gin.Context) {

	var req model.LoginUserRequest

	//check req from from frontend
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	FoundUser, err := h.queries.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   err.Error(),
				"message": "invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "login failed",
		})
		return
	}

	//compare password
	correct := helpers.CheckPasswordHash(req.Password, FoundUser.PasswordHash)
	if !correct {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "login failed, invaid password",
		})
		return
	}

	//TODO: lock out user for multiple wrong password

	//generate access | create a 3rd func in uth to generate both access and refresh token

	accessTkn, err := GenerateAccess(FoundUser.ID, FoundUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "login failed",
			"error":   err.Error(),
		})
		return
	}

	refreshTkn, err := GenerateRefresh(FoundUser.ID, FoundUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "login failed",
			"error":   err.Error(),
		})
		return
	}

	//store the refresh toekn
	_, err = h.queries.CreateRefreshToken(c.Request.Context(), dbq.CreateRefreshTokenParams{
		ID:        uuid.NewString(),
		UserID:    FoundUser.ID,
		TokenHash: helpers.HashToken(refreshTkn),
		// ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save session",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            FoundUser.ID,
		"access_token":  accessTkn,
		"refresh_token": refreshTkn,
	})

}

// refresh
func (h *Handler) RefreshToken() {

}

// logout - takehome
func (h *Handler) Logout() {

}
