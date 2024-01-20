#!/bin/bash
set -euo pipefail

source .env

# Use MYSQL_PWD environment variable for password to avoid exposing in command line
export MYSQL_PWD="${MYSQL_ROOT_PASSWORD}"

echo "Creating user schema..."
TTY_DISABLED=true mysql -h "${DB_HOST}" -P "${DB_PORT}" -u "${MYSQL_ROOT_USER}" -p"${MYSQL_ROOT_PASSWORD}"< ./erpdb/scripts/drop.sql