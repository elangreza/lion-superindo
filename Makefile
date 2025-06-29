run:
	go run ./cmd/server/.
	
gen:
	go generate ./...

test:
	go test ./...

up:
	cp ./env.example docker.env
	docker compose --env-file docker.env up -d --build
	cat ./db/seed/seed_1.sql | docker exec -i superindo-database psql -h localhost -U superindo -f-

down:
	docker compose down

swag:
	swag fmt
	swag init -g cmd/server/main.go

wire:
	wire gen ./cmd/server/

.PHONY: run gen test up down swag