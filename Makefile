-include .env

DATABASE=mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=true&loc=Local
BINARY_NAMES="golang-migrate|migrate"

dbu: ## Migrates the database to the latest version.
	migrate -database="$(DATABASE)" -path=db/migrations -lock-timeout=20  -verbose up

dbdrop: ## Drops the database.
	migrate -database="$(DATABASE)" -path=db/migrations -lock-timeout=20 -verbose drop -f
