#!/bin/bash

value="auth"
# Loop over each value and run the command
echo "Running seed for $value"
for file in ./db/seeds/*.sql; do
    filename=$(basename "$file")
    file_path=/var/lib/postgresql/seeds/$filename
    echo "Running seed file $file_path"
    docker-compose exec db psql -U postgres -d "$value" -a -f "$file_path"
done

