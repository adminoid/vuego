include .env

.PHONY: up down
db := postgres://${DbUser}:${DbPwd}@localhost:5432/${DbName}?sslmode=disable
up:
	migrate -database ${db} -path db/migrations up
down:
	migrate -database ${db} -path db/migrations down
