#!/bin/sh
set -e
sh /app/database/scripts/migrate.sh
exec "/app/application"