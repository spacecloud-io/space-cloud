#!/bin/sh
docker rmi -f sc-gateway
docker rmi -f spaceuptech/gateway:0.21.4
docker rmi -f spaceuptech/gateway:latest

docker build --no-cache -t sc-gateway .

docker tag sc-gateway spaceuptech/gateway:0.21.4
docker tag sc-gateway spaceuptech/gateway:latest

docker push spaceuptech/gateway:0.21.4
docker push spaceuptech/gateway:latest
