.PHONY: export-db import-db

DB_USER=postgres
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=postgres

export PGUSER=$(DB_USER)
export PGPASSWORD=$(DB_PASSWORD)
export PGHOST=$(DB_HOST)
export PGPORT=$(DB_PORT)
export PGDATABASE=$(DB_NAME)

export-db:
	- mkdir -p out
	- docker-compose exec -T database pg_dump -U $(DB_USER) $(DB_NAME) > out/db_dump.sql

import-db:
	cat out/db_dump.sql | docker-compose exec -T database psql -U $(DB_USER) $(DB_NAME)
