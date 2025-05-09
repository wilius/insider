echo "Running migration on: ${DB_HOST}:${DB_PORT}"
migrate -path ./database/migrations \
  -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"\
  up
