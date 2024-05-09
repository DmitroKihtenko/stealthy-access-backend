#!/bin/bash

INPUT_FILE="/home/ory/kratos-template.yml"
OUTPUT_FILE="/home/ory/kratos.yml"

echo "Running Kratos init script..."

if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: $INPUT_FILE not found"
    exit 1
fi

if [ -f "$OUTPUT_FILE" ]; then
    echo "Deleting old $OUTPUT_FILE file"
    rm -f $OUTPUT_FILE
fi

envsubst < "$INPUT_FILE" > "$OUTPUT_FILE"

echo "$OUTPUT_FILE generated successfully"

kratos migrate sql postgres://${SERVICE_DB_USER:-"access-backend"}:${SERVICE_DB_PASSWORD:-"password"}@${DATABASE_HOST:-"postgres-db"}:${DATABASE_PORT:-"5432"}/${SERVICE_DB_NAME:-"kratos"}?sslmode=disable --yes

echo "Initilization successful"
