.PHONY: run test lint migrate-up migrate-down

run:
	go run cmd/api/main.go

test:
	go test ./... -v

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "$$DATABASE_DSN" up

migrate-down:
	migrate -path migrations -database "$$DATABASE_DSN" down
