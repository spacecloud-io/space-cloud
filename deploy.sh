#!/bin/bash

mkdir linux && mkdir windows && mkdir darwin
GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w -extldflags "-static"' .
zip space-cloud.zip space-cloud
mv ./space-cloud.zip ./darwin/
rm space-cloud
GOOS=windows GOARCH=amd64 go build -ldflags '-s -w -extldflags "-static"' .
zip space-cloud.zip space-cloud.exe
mv ./space-cloud.zip ./windows/
rm space-cloud.exe
GOOS=linux GOARCH=amd64 go build -ldflags '-s -w -extldflags "-static"' .
zip space-cloud.zip space-cloud
mv ./space-cloud.zip ./linux/
rm space-cloud

# echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
# docker push spaceuptech/space-cloud:latest

curl -H "Authorization: Bearer $JWT_TOKEN" -F 'file=@./darwin/space-cloud.zip' -F 'fileType=file' -F 'makeAll=false' -F 'path=/darwin' https://spaceuptech.com/v1/api/downloads/files
curl -H "Authorization: Bearer $JWT_TOKEN" -F 'file=@./windows/space-cloud.zip' -F 'fileType=file' -F 'makeAll=false' -F 'path=/windows' https://spaceuptech.com/v1/api/downloads/files
curl -H "Authorization: Bearer $JWT_TOKEN" -F 'file=@./linux/space-cloud.zip' -F 'fileType=file' -F 'makeAll=false' -F 'path=/linux' https://spaceuptech.com/v1/api/downloads/files
