export GOOSE_DRIVER = postgres
export GOOSE_DBSTRING = postgres://postgres:your_secure_password@localhost:5432/mydatabase?sslmode=disable

run:
	go run main.go

migrate:
	cd db/migrations && goose up
