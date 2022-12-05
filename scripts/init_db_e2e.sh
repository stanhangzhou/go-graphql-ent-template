#!/bin/bash

docker compose cp ./scripts/sql postgres:/var/lib/postgresql
docker compose exec postgres bash -c "psql -U postgres  < /var/lib/postgresql/sql/reset_e2e_database.sql"