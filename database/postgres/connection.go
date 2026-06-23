package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	dbq "github.com/valentineejk/voters_api/database/sqlc"
)

func Connection() *dbq.Queries {

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// pgxpool manages a pool of connections
	// never open a new connection per request
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer pool.Close()

	// ping to verify connection is alive at startup
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}
	fmt.Println("connected to postgres")

	// db.New wraps the pool and gives us typed query methods
	queries := dbq.New(pool)

	return queries
}
