VERSION := 1.0

include .env

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o playbypost cmd/server/main.go

tidy:
	go get -u -v ./...
	go mod tidy

migrate_up:
	migrate -path=core/sys/db/migrate/migrations -database "postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=core/sys/db/migrate/migrations -database "postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=disable" -verbose down 1

psql:
	psql -h localhost -U postgres -p 5432

drop_db:
	psql -h localhost -U postgres -p 5432 -c "DROP DATABASE playbypost;"

create_db:
	psql -h localhost -U postgres -p 5432 -c "CREATE DATABASE playbypost;"

stop_postgres:
	docker stop postgres 

start_postgres:
	docker start postgres