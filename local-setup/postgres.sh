#!/bin/bash

echo "Pulling the latest PostgreSQL image..."
docker pull postgres:latest

echo "Creating a new PostgreSQL container..."
docker run --name my-postgres-container -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres

echo "Waiting for PostgreSQL to start and be ready to accept connections..."
while ! docker logs my-postgres-container | grep -q "database system is ready to accept connections"; do
  sleep 3
done

echo "Checking PostgreSQL logs..."
docker logs my-postgres-container

echo "Verifying the PostgreSQL container..."
docker ps | grep my-postgres-container

echo "PostgreSQL container setup is complete."

export LOCAL_POSTGRES_CONNECTION_STRING="postgresql://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
echo "Postgres local connection string [$LOCAL_POSTGRES_CONNECTION_STRING]"