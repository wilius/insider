SHELL := /bin/bash
.SILENT:
clean-env:
	@ find . -type f -name ".env" -exec sh -c 'rm -rf $${1%}' _ {} \;

env:
	@ find . -type f -name ".env.example" -exec sh -c '[ -e "$${1%.example}" ] || cp -v "$${1%}" "$${1%.example}"' _ {} \;

create-migration: env
	docker compose run --rm --user "$(shell id -u):$(shell id -g)" app "database/scripts/create_migration.sh $(name)"

migrate: env
	docker compose run --rm app "sleep 3 && database/scripts/migrate.sh"

rollback: env
	docker compose run --rm app "sleep 3 && database/scripts/rollback.sh $(count)"


build: env
	docker compose build ...

debug-build: env
	DOCKER_BUILDKIT=1 ... --progress=plain

server: build
	docker compose up -d

develop:
	docker compose up postgres redis webhook -d

stop:
	docker compose down

clean:
	docker compose down --remove-orphans --volumes

clean-all: clean clean-env