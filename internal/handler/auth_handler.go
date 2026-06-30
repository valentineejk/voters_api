package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbq "github.com/valentineejk/voters_api/database/sqlc"
	"github.com/valentineejk/voters_api/internal/helpers"
	"github.com/valentineejk/voters_api/internal/model"
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
func (h *Handler) RefreshToken(c *gin.Context) {

	var req model.RefreshToken

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. validate JWT signature and expiry
	claims, err := Validate(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired refresh token",
		})
		return
	}

	// look up the hash in the DB
	// if not found: already used (rotation) or revoked or never existed

	tokenRecord, err := h.queries.GetRefreshToken(
		c.Request.Context(), helpers.HashToken(req.RefreshToken))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token not found or already used",
		})
		return
	}

	// revoke the OLD refresh token (rotation — one use only)
	if err := h.queries.RevokeRefreshToken(
		c.Request.Context(), tokenRecord.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke token"})
		return
	}

	// fetch the user to get current role (may have changed since last login)
	user, err := h.queries.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// issue new access token
	newAccessToken, err := GenerateAccess(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 6. issue new refresh token and store its hash
	newRefreshToken, err := GenerateRefresh(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	_, err = h.queries.CreateRefreshToken(c.Request.Context(), dbq.CreateRefreshTokenParams{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: helpers.HashToken(newRefreshToken),
		// ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// logout - takehome
func (h *Handler) Logout(c *gin.Context) {

	userID, _ := c.Get("user_id")

	// revoke ALL refresh tokens for this user
	// this logs them out of every device simultaneously

	err := h.queries.RevokeAllUserTokens(
		c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
				"code":  "MISSING_TOKEN",
			})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := Validate(tokenStr)

		if err != nil {

			code := "INVALID_TOKEN"
			if errors.Is(err, jwt.ErrTokenExpired) {
				code = "TOKEN_EXPIRED"
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
				"code":  code,
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// role-based access — only admins
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
				"code":  "FORBIDDEN",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// reading user from context inside any protected handler:
func (h *Handler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	user, err := h.queries.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  role,
	})
}
