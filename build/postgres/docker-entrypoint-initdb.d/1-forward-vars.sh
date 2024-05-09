#!/bin/bash

SERVICE_DB_NAME=${SERVICE_DB_NAME:-"kratos"}
SERVICE_DB_USER=${SERVICE_DB_USER:-"access-backend"}
SERVICE_DB_PASSWORD=${SERVICE_DB_PASSWORD:-"password"}
SQL_VARS_FILE="/docker-entrypoint-initdb.d/init/variables.sql"

forward_variable() {
  local var=$1
  local value=$2
  local sql_vars_file=${3:-$SQL_VARS_FILE}
  echo "\set $var $value" >> "$sql_vars_file"
}

main() {
  echo -n > $SQL_VARS_FILE
  forward_variable database "$SERVICE_DB_NAME"
  forward_variable user_name "$SERVICE_DB_USER"
  forward_variable user_password "$SERVICE_DB_PASSWORD"
}

main
