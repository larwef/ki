#!/bin/bash

cp -p "../target/$(ls -t ../target/app | grep -v /orig | head -1)" toRoot/app

docker build -t go-docker-test .
docker run -it --rm -p 8080:8080 -p 8081:8081  --env-file ki.properties go-docker-test ./app -disable-tls=true