#!/bin/bash
set -e

echo "Starting database initialization script..."
echo "POSTGRES_USER: $POSTGRES_USER"
echo "POSTGRES_DB: $POSTGRES_DB"

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

echo "Database initialization completed." 