create-migrate:
	migrate create -ext sql -dir ./migrations -seq regauth

migrate-up:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/regauth?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/regauth?sslmode=disable' down

migrate-force:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/regauth?sslmode=disable' force 1