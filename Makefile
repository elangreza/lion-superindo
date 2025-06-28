run:
	go run cmd/server/main.go
	
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

.PHONY: migrate migrate-create run gen up down test