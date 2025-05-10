connectionString="postgresql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable"
migrate -database "$connectionString" create -ext sql -dir database/migrations $1
