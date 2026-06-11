run:
	go run cmd/main.go

build:
	go build cmd/main.go

wire:
	cd internal/injector && wire

air:
	air -c air.toml

test:
	go test -v  ./tests/v1

docker-up:
	docker compose up -d --build

container-reset:
	docker-compose down && docker-compose up -d