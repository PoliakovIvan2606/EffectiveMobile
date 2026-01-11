run:
	go run cmd/app/main.go

build:
	mkdir -p volumes/pgdata/
	docker-compose up --build