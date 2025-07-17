postgres:
	docker run --name mypostgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=qwe123 -d postgres:latest

createdb:
	docker exec -it mypostgres createdb --username=root --owner=root simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:qwe123@localhost:5432/simple_bank?sslmode=disable" -verbose up


migrateup1:
	migrate -path db/migration -database "postgres://root:qwe123@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgres://root:qwe123@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgres://root:qwe123@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
	

dropdb:
	docker exec -it mypostgres dropdb --username=root --owner=root simple_bank 

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/JaidenShall/simplebank/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migratedown sqlc test server mock
