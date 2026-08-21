.PHONY: migration migrate-up dev

migration:
	goose -dir migration create $(name) sql

migrate-up:
	goose -dir migration postgres "$(DATABASE_URL)" up

migrate-status:
	goose -dir migration postgres "$(DATABASE_URL)" status

dev:
	go run ./cmd/server

debug:
	@echo "DATABASE_URL=$(DATABASE_URL)"
