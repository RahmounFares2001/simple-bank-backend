postgres:
	sudo docker run --name bankContainer -p 5432:5432 -e POSTGRES_USER=fares -e POSTGRES_PASSWORD=fares -d postgres:latest

createdb:
	sudo docker exec -it bankContainer createdb --username=fares --owner=fares bankDB

dropdb:
	sudo docker exec -it bankContainer dropdb bankDB

migrateup:
	migrate -path db/migration -database "postgres://fares:fares@localhost:5432/bankDB?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgres://fares:fares@localhost:5432/bankDB?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgres://fares:fares@localhost:5432/bankDB?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgres://fares:fares@localhost:5432/bankDB?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go
	
.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server