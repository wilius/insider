# Project Setup

This document provides the necessary steps to set up the environment and run the application. Please follow the instructions carefully to ensure a smooth setup.
<br>
To skip all the details and run the project, run the command below if you make sure the prerequisites met by your environment:

swaggerhub deployment also accessible [here](https://app.swaggerhub.com/apis/insider-435/messaging/11). The swagger definition is also attached to the project called `swagger-openapi.yaml`
```sh
  make server
``` 

## **Running the Application**

```sh
make server       # Run the server
```

Access the application at [http://localhost:8080](http://localhost:8080). The postman collection is also shared [here](https://tinyurl.com/2p9zaemd). Please import into your Postman and start consuming what app does

---

## **Stopping the Application**

```sh
make stop
```

To fully clean up Docker volumes and environment files:

```sh
make clean-all
```

---

## **Prerequisites**

Ensure the following tools are installed on your machine:

1. **Docker**

2. **Docker Compose**

3. **Make Utility**

    * For Ubuntu/Debian:

   ```sh
   sudo apt-get update
   sudo apt-get install -y make
   ```

    * For MacOS:

   ```sh
   brew install make
   ```

---

## **Environment Variables**

Below are the environment variables required for the application:

```env
DATABASE_HOST=postgres
DATABASE_PORT=5432
DATABASE_NAME=insider
DATABASE_PASSWORD=postgres
DATABASE_USERNAME=postgres

REDIS_HOST=redis
REDIS_PORT=6379

RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_VHOST=insider
RABBITMQ_PASSWORD=guest
RABBITMQ_USERNAME=guest

WEBHOOK_PROVIDER_WEBHOOKSITE_URL=http://webhook:5000/webhook
```

Make sure these variables are present in your `.env` file or set in your environment before running the application. if you're trying to compose it up from scratch, make is going to create one for you.

---

## **Makefile Commands**

### **Setup Commands**

* `make env` - Creates `.env` files from `.env.example` if they don't exist.
* `make clean-env` - Removes all `.env` files.

### **Database Migrations**

* `make create-migration name=<migration_name>` - Creates a new database migration script.
* `make migrate` - Runs the database migrations.
* `make rollback count=<number>` - Rolls back the last `<number>` migrations.

### **Build & Server Commands**

* `make build` - Builds the Docker containers for `app` and `webhook`.
* `make debug-build` - Builds with detailed output.
* `make server` - Runs the application server with live code reloading.

### **Environment Management**

* `make develop` - Starts `postgres`, `redis`, `rabbitmq`, and `webhook` services in detached mode.
* `make stop` - Stops all running containers.
* `make clean` - Removes all Docker containers and orphans.
* `make clean-all` - Cleans all containers and removes `.env` files.

---

## **Configuration (config.yaml)**

The application configuration is stored in `config.yaml`. Below is the structure:

```
server:
  port: {{ default .SERVER_PORT 8080 }}

database:
  host: {{ default .DATABASE_HOST "localhost" }}
  port: {{ default .DATABASE_PORT 5432 }}
  name: {{ default .DATABASE_NAME "insider" }}
  username: {{ default .DATABASE_USERNAME "postgres" }}
  password: {{ default .DATABASE_PASSWORD "postgres" }}
  logLevel: {{ default .DATABASE_LOG_LEVEL 1 }}
  pool:
    maxIdleConnections: {{ default .DATABASE_MAX_IDLE_CONNECTIONS 5 }}
    maxOpenConnections: {{ default .DATABASE_MAX_OPEN_CONNECTIONS 100 }}
    connectionMaxLifetime: {{ default .DATABASE_CONNECTION_MAX_LIFETIME 60000 }}
    connectionMaxIdleTime: {{ default .DATABASE_CONNECTION_MAX_IDLE_TIME 1000 }}

rabbitmq:
  host: {{ default .RABBITMQ_HOST "localhost" }}
  port: {{ default .RABBITMQ_PORT 5672 }}
  vhost: {{ default .RABBITMQ_VHOST "insider" }}
  username: {{ default .RABBITMQ_USERNAME "guest" }}
  password: {{ default .RABBITMQ_PASSWORD "guest" }}

redis:
  host: {{ default .REDIS_HOST "localhost" }}
  port: {{ default .REDIS_PORT 6379 }}

scheduler:
  interval: {{ default .SCHEDULER_INTERVAL "10s" }}
  itemCountPerCycle: {{ default .SCHEDULER_ITEM_COUNT_PER_CYCLE 2 }}

provider:
  type: {{ default .WEBHOOK_PROVIDER "webhook.site" }}
  url: {{ default .WEBHOOK_PROVIDER_WEBHOOKSITE_URL "http://localhost:5000/webhook" }}
  requestTimeout: {{ default .WEBHOOK_PROVIDER_WEBHOOKSITE_REQUEST_TIMEOUT "5s" }}
  circuitBreaker:
    maxRequests: 5
    maxFailure: 5
    interval: 60s
    timeout: 5s
```

You can modify the default values as needed.

---

## **Docker Compose Services**

Below is the structure of your `docker-compose.yml`:

```yaml
services:
  app:
    build:
      context: .
      target: dev
    entrypoint: sh -c
    init: true
    pull_policy: "build"
    environment:
      DOCKER_BUILDKIT: 1
      COMPOSE_DOCKER_CLI_BUILD: 1
    env_file: .env
    volumes:
      - .:/app
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis

  webhook:
    build:
      context: test_files/docker/python
    init: true
    pull_policy: "build"
    environment:
      DOCKER_BUILDKIT: 1
      COMPOSE_DOCKER_CLI_BUILD: 1
    ports:
      - "5000:5000"

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: insider
    hostname: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - '5432:5432'
    healthcheck:
      test: [ "CMD-SHELL", "pg_isready -U postgres -d insider" ]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis/redis-stack:7.4.0-v3
    hostname: redis
    ports:
      - '6379:6379'
    volumes:
      - redis_data:/data

  rabbitmq:
    image: rabbitmq:3.12-management
    hostname: rabbitmq
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    ports:
      - 15672:15672
      - 5672:5672
    environment:
      - RABBITMQ_DEFAULT_USER=guest
      - RABBITMQ_DEFAULT_PASS=guest
      - RABBITMQ_DEFAULT_VHOST=insider

volumes:
  postgres_data:
  redis_data:
  rabbitmq_data:
```

This defines all services, networking, and data volumes for persistence.

---

## **Troubleshooting**

If you encounter any issues:

* Verify Docker and Docker Compose are installed and running.
* Ensure ports `8080`, `5432`, `6379`, and `5672` are available.
* Use `docker compose logs <service>` to inspect service logs.
