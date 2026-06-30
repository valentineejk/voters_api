package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	database "github.com/valentineejk/voters_api/database/postgres"
	"github.com/valentineejk/voters_api/internal/handler"
)

func main() {

	// Load .env into the process environment. Not fatal if it's missing —
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on existing environment")
	}

	PORT := ":8000"

	q, p := database.Connection()
	defer p.Close()

	h := handler.New(q)

	r := gin.Default()

	v1 := r.Group("/api/v1")

	//protected routes
	protected := v1.Group("/", handler.AuthMiddleware())
	{
		protected.GET("/voters/:id", h.Get_one_voter)
		protected.DELETE("/voters/:id", h.Delete_voter)
		protected.GET("/voters", h.GetAllVoters)
	}

	v1.POST("/voters", h.Register_voter)
	v1.PUT("/voters/:id/status", h.Update_voter_status)
	v1.GET("/health", h.HealthCheck)

	auth := v1.Group("/auth")
	auth.POST("/register", h.RegisterHandler)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)

	r.Run(PORT)

}

//revoke - take home group
