DB=postgresql://root:pasd@localhost:5432/allegro_trade?sslmode=disable

pgsql:
	docker run --name postgreSQL -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pasd -d postgres:14-alpine

createdb:
	docker exec -it postgreSQL createdb --username=root --owner=root allegro_trade

dropdb:
	docker exec -it postgreSQL dropdb allegro_trade

migrateup:
	migrate -path db/migration -database "$(DB)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -count=1 ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/YuanData/allegro-trade/db/sqlc Store

.PHONY: pgsql
