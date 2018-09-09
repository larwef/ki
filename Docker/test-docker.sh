#!/bin/bash

cp -p "../target/$(ls -t ../target/app | grep -v /orig | head -1)" toRoot/app

docker build -t go-docker-test .
docker run -it --rm -p 8080:8080 go-docker-test