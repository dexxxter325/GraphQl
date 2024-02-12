migrate-up:
	migrate -path ./postgres/migrations -database "postgres://postgres:qwerty@localhost:5431/postgres?sslmode=disable" up
migrate-down:
	migrate -path ./postgres/migrations -database "postgres://postgres:qwerty@localhost:5431/postgres?sslmode=disable" down
