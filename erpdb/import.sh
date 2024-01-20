#!/bin/bash
set -euo pipefail

source .env

# Use MYSQL_PWD environment variable for password to avoid exposing in command line
export MYSQL_PWD="${MYSQL_ROOT_PASSWORD}"

echo "Creating user tables..."
TTY_DISABLED=true mysql -h "${DB_HOST}" -P "${DB_PORT}" -u "${MYSQL_ROOT_USER}" -p"${MYSQL_ROOT_PASSWORD}"< ./erpdb/scripts/user.sql
# echo "Creating product tables..."
# TTY_DISABLED=true mysql -h "${DB_HOST}" -P "${DB_PORT}" -u "${MYSQL_ROOT_USER}" -p"${MYSQL_ROOT_PASSWORD}"< ./erpdb/scripts/product.sql
# echo "Inserting seed data..."
# TTY_DISABLED=true mysql -h "${DB_HOST}" -P "${DB_PORT}" -u "${MYSQL_ROOT_USER}" -p"${MYSQL_ROOT_PASSWORD}"< ./erpdb/scripts/seed.sql