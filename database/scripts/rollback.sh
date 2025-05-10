connectionString="postgresql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable"
count=${1:-1}
echo "Rolling back to $count version lower on connectionString ${connectionString}"
migrate -path ./database/migrations -database "$connectionString" down $count
