echo "Rolling back migration on: ${DATABASE_HOST}:${DATABASE_PORT} to $1 version lower"
migrate -path ./database/migrations \
  -database "postgresql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable"\
  down $1
