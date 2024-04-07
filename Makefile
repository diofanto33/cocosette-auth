.PHONY: all build run dkpsql clean

DB_USER=pepe 
DB_PASSWORD=pepe123
DB_PORT=5432
CONTAINER_NAME=cocosette
DB_DRIVER=postgres 

all: build run connect clean

make run:
	go run ./cmd/main.go 

build:
	go build -o bin/zulo ./cmd/api

dkpsql:
	docker run --name $(CONTAINER_NAME) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-p $(DB_PORT):$(DB_PORT) \
		-d $(DB_DRIVER)
	sleep 2
	docker exec -it $(CONTAINER_NAME) bash -c "\
		PGPASSWORD=$(DB_PASSWORD) psql -U $(DB_USER) -c 'CREATE DATABASE \"auth_svc\";' && \
		PGPASSWORD=$(DB_PASSWORD) psql -U $(DB_USER) -d \"auth_svc\";" 

host:
	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(CONTAINER_NAME)


clean:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true
	kill -9 $(shell lsof -t -i:50051) || true
