FROM golang:alpine AS base
RUN mkdir -p /app
RUN apk update && apk add ca-certificates tzdata
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
WORKDIR /app

FROM base AS dev

FROM base AS builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app/application

FROM alpine:3.20
WORKDIR /app
COPY --from=builder $GOROOT/go/bin/migrate /usr/local/bin/

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/application /app/application
COPY --from=builder /app/database/migrations /app/database/migrations
COPY --from=builder /app/database/scripts/migrate.sh /app/database/scripts/migrate.sh
COPY docker-entrypoint.sh .

RUN chmod +x ./docker-entrypoint.sh
CMD ["/bin/sh", "/app/docker-entrypoint.sh"]
