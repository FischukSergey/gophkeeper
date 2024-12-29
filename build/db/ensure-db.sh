#!/bin/bash
set -e

until pg_isready -U "$POSTGRES_USER"; do
    echo "Waiting for PostgreSQL to start..."
    sleep 1
done

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    DO
    \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gophkeeper') THEN
            CREATE DATABASE gophkeeper;
        END IF;
    END
    \$\$;
EOSQL 