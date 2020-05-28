#!/bin/sh

db_path=/tmp/fv-db

docker run --rm -d \
    --name fv-db \
    --network=host \
    -e POSTGRES_PASSWORD=fvpgsecret \
    -e PGDATA=/var/lib/postgresql/data \
    -v ${db_path}:/var/lib/postgresql/data \
    -p 5432:5432 \
    postgres

echo "initializing db table"
sleep 5

# init db
docker run --rm -i \
    --network=host \
    -e POSTGRES_PASSWORD=fvpgsecret \
    postgres \
    psql -h localhost -U postgres < scripts/init.sql
