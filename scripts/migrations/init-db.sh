#!/bin/bash
set -e

echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q'; do
  sleep 1
done

echo "PostgreSQL is ready. Running migrations..."

# Run all migrations in order (sorted by filename)
for migration in $(ls -1 /migrations/*.sql | sort -V); do
    echo "Running migration: $(basename $migration)"
    PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$migration"
done

echo "Database initialized successfully!"

