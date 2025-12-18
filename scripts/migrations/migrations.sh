#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."
MAX_RETRIES=30
RETRY_COUNT=0

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' 2>/dev/null; do
  RETRY_COUNT=$((RETRY_COUNT + 1))
  if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
    echo "ERROR: Failed to connect to PostgreSQL after $MAX_RETRIES attempts"
    echo "Check your POSTGRES_PASSWORD in .env file matches the database"
    exit 1
  fi
  echo "Attempt $RETRY_COUNT/$MAX_RETRIES: PostgreSQL not ready yet, waiting..."
  sleep 1
done

echo "PostgreSQL is ready. Running migrations..."

# Run all migrations in order (sorted by filename)
for migration in $(ls -1 /migrations/*.sql | sort -V); do
    echo "Running migration: $(basename $migration)"
    PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$migration"
done

echo "Database initialized successfully!"

