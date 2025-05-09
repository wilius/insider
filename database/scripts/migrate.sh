connectionString="postgresql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable"
echo "Running migration on: ${connectionString}"
migrate -path ./database/migrations -database "$connectionString" up
