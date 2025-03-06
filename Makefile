db_login:
	psql ${DATABASE_URL}

migrateCreate:
	migrate create -ext sql -dir migrations -seq $(name)

migrateUp:
	migrate -database ${DATABASE_URL} -path migrations up

migrateDown:
	migrate -database ${DATABASE_URL} -path migrations down 1

migrateDrop:
	migrate -database ${DATABASE_URL} -path migrations drop

dockerUp:
	docker compose up -d

dockerDown:
	docker compose down

http:
	go run . http


