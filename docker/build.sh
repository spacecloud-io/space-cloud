#!/usr/bin/env bash

echo "##################### Building space-cloud Docker Image! #####################"
docker build -t rawsanj/space-cloud space-cloud-docker/

echo ""
echo "##################### Building realtime-react-chat App Docker Image! #####################"
docker build -t rawsanj/realtime-chat-react realtime-chat-react/

echo ""
echo "##################### Publishing Docker Images to Docker Hub. #####################"
echo "NOTE: change the username (rawsanj) to your Dockerhub username when building and publish, if you want to publish those in your account."

docker push rawsanj/space-cloud

docker push rawsanj/realtime-chat-react
