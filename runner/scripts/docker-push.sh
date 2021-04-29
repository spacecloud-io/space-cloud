#!/bin/sh
docker rmi -f sc-runner
docker rmi -f spaceuptech/runner:0.21.4
docker rmi -f spaceuptech/runner:latest

docker build --no-cache -t sc-runner .

docker tag sc-runner spaceuptech/runner:0.21.4
docker tag sc-runner spaceuptech/runner:latest

docker push spaceuptech/runner:0.21.4
docker push spaceuptech/runner:latest

