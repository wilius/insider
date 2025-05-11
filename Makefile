SHELL := /bin/bash
.SILENT:
clean-env:
	@ find . -type f -name ".env" -exec sh -c 'rm -rf $${1%}' _ {} \;

env:
	@ find . -type f -name ".env.example" -exec sh -c '[ -e "$${1%.example}" ] || cp -v "$${1%}" "$${1%.example}"' _ {} \;

create-migration: env
	docker compose run --rm --user "$(shell id -u):$(shell id -g)" app "database/scripts/create_migration.sh $(name)"

migrate: env
	docker compose run --rm app "database/scripts/migrate.sh"

rollback: env
	docker compose run --rm app "database/scripts/rollback.sh $(count)"

build: env
	docker compose build app webhook

server: build
	docker compose run --service-ports --rm app './docker-entrypoint.sh'

develop:
	docker compose up postgres redis rabbitmq webhook -d

stop:
	docker compose down

clean:
	docker compose down --remove-orphans --volumes

clean-all: clean clean-env