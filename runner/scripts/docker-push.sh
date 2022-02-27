#!/bin/sh
docker rmi -f sc-runner
docker rmi -f spacecloudio/runner:0.21.5
docker rmi -f spacecloudio/runner:latest

docker build --no-cache -t sc-runner .

docker tag sc-runner spacecloudio/runner:0.21.5
docker tag sc-runner spacecloudio/runner:latest

docker push spacecloudio/runner:0.21.5
docker push spacecloudio/runner:latest

