# Migrations
#
# Examples:
#   make migrate-create name=create_users_table   # new migration pair
#   make migrate-up                                # apply all pending migrations
#   make migrate-down                              # roll back the last migration
#   make sqlc                                      # regenerate typed query code
#   make start                                     # run the API

# usage: make migrate-create name=create_users_table
migrate-create:
	migrate create -ext sql -dir ./database/migrations -seq $(name)

migrate-up:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" up

migrate-down:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" down 1

sqlc:
	#generate sql commands
	sqlc generate

start:
	go run .
