postgres:
	sudo docker run --name bankContainer -p 5432:5432 -e POSTGRES_USER=fares -e POSTGRES_PASSWORD=fares -d postgres:latest

createdb:
	sudo docker exec -it bankContainer createdb --username=fares --owner=fares bankDB

dropdb:
	sudo docker exec -it bankContainer dropdb bankDB

migrateup:
	migrate -path db/migration -database "postgres://fares:fares@localhost:5432/bankDB?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://fares:fares@localhost:5432/bankDB?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...
	
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test