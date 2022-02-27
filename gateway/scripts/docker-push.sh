#!/bin/sh
docker rmi -f sc-gateway
docker rmi -f spacecloudio/gateway:0.21.5
docker rmi -f spacecloudio/gateway:latest

docker build --no-cache -t sc-gateway .

docker tag sc-gateway spacecloudio/gateway:0.21.5
docker tag sc-gateway spacecloudio/gateway:latest

docker push spacecloudio/gateway:0.21.5
docker push spacecloudio/gateway:latest
