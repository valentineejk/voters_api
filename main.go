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

	v1.POST("/polling-stations", h.AddPollingStation)
	v1.GET("/polling-stations", h.GetAllPollingStations)
	v1.GET("/polling-stations/:id", h.GetPollingStation)

	v1.GET("/health", h.HealthCheck)

	r.Run(PORT)

}

//update_voter_status -
//delete_voter -
//get_all_voters -
//add state vildation function to the create voter, valid nin to the struct, valid phone
