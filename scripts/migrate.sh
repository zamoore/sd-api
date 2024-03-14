#!/bin/bash

# Check if the first argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 [up|down]"
  exit 1
fi

DATABASE_URL="postgres://zackmoore@localhost:5432/snipdrop?sslmode=disable"
MIGRATIONS_PATH="database/migrations"

case "$1" in
  up)
    echo "Migrating up..."
    migrate -database "$DATABASE_URL" -path $MIGRATIONS_PATH up
    ;;
  down)
    echo "Migrating down..."
    migrate -database "$DATABASE_URL" -path $MIGRATIONS_PATH down
    ;;
  *)
    echo "Invalid argument: $1"
    echo "Usage: $0 [up|down]"
    exit 1
    ;;
esac
