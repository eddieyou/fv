#!/bin/sh


docker build -t fv-app -f ./build/Dockerfile .

docker run --rm -d \
    --network=host \
    -e FV_DB="host=localhost port=5432 user=postgres password=fvpgsecret sslmode=disable" \
    --name fv-app \
    fv-app \
    fv
