#!/bin/bash


echo "Running migrations for auth"
migrations_dir=./db/migrations
echo "Running migrations from $migrations_dir"
dbmate -d $migrations_dir -u postgres://postgres:postgres@localhost:5432/auth?sslmode=disable --no-dump-schema up


