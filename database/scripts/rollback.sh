echo "Rolling back migration on: ${DB_HOST}:${DB_PORT} to $1 version lower"
migrate -path ./database/migrations \
  -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"\
  down $1
