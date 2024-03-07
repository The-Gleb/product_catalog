postgres:
	docker run --name catalog_db -e POSTGRES_USER=catalog_db -e POSTGRES_PASSWORD=catalog_db -p 5434:5432 -d postgres:alpine

postgresrm:
	docker stop catalog_db
	docker rm catalog_db

migrateup:
	migrate -path internal/adapter/db/migration -database "postgres://catalog_db:catalog_db@localhost:5434/catalog_db?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/adapter/db/migration -database "postgres://catalog_db:catalog_db@localhost:5434/catalog_db?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown