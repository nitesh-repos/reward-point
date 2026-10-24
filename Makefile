.PHONY: run test cover docker-up docker-down

run:
	go run ./cmd/reward-point

test:
	go test -v ./...

cover:
	go test -cover ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down
