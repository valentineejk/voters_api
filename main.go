package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	database "github.com/valentineejk/voters_api/database/postgres"
	"github.com/valentineejk/voters_api/handler"
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

	v1.POST("/voters", h.Register_voter)
	v1.GET("/voters/:id", h.Get_one_voter)
	v1.PUT("/voters/:id/status", h.Update_voter_status)
	v1.GET("/health", h.HealthCheck)

	auth := r.Group("/api/v1/auth")
	auth.POST("/register", h.RegisterHandler)
	auth.POST("/login", h.Login)

	r.Run(PORT)

}

//update_voter_status -
//delete_voter -
//get_all_voters -
//add state vildation function to the create voter, valid nin to the struct, valid phone
