GOOSE_DBSTRING ?= "postgresql://vinhtsts99:YmCbT3klRAZ9ZvmXNRNm3LQfwfSWrhow@dpg-d4c5eeali9vc73bpqsr0-a.oregon-postgres.render.com/taskly_us9h"
GOOSE_MIGRATION_DIR ?= sql/schema
GOOSE_DRIVER ?= postgres
APP_NAME = server
docker_build:
	docker-compose up -d --build 
	docker-compose ps
docker_stop:
	docker-compose down
dev: 
	go run ./cmd/${APP_NAME}
docker_up:
	docker-compose up -d

up_by_one:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) up-by-one
create_migration:
	@goose -dir=$(GOOSE_MIGRATION_DIR) create $(name) sql
upse:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) up 
downse: 
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) down 

resetse:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) reset

dev1: 
	echo "GOOSE_DBSTRING=$(GOOSE_DBSTRING)"
	echo "GOOSE_DRIVER=$(GOOSE_DRIVER)"
	echo "GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR)"
sqlgen:
	sqlc generate

swag:
	swag init -g ./cmd/server/main.go -o ./cmd/swag/docs
.PHONY: dev downse upse resetse dev docker_build docker_stop docker_up swag
.PHONY: air
