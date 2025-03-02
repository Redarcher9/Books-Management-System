DBUrl=""

start:
	go run cmd/main.go

swagger:
	swag init -g ./cmd/main.go -o ./docs

migrateup:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/books-management-system?sslmode=disable" up

startredis:
	/opt/homebrew/opt/redis/bin/redis-server /opt/homebrew/etc/redis.conf