
## Build Docker Images

Run the `build.sh` file to create and push Docker images.

```bash
$ ./build.sh
```

**Note:** To publish docker images in your DockerHub account. Change the username (rawsanj) to your DockerHub username when building and publish.
Also update the docker image names in `docker-compose.yml` file.

## Run Space Cloud with MongoDb and Sample app

Run below command under `space-cloud/docker` directory.

```bash
$ docker-compose up
``` 

## Deploy to Play-with-Docker


Ctrl + Click below button to deploy space cloud stack with sample app to [Play with Docker](https://labs.play-with-docker.com).

[![Run](https://img.shields.io/badge/RUN-Deploy%20to%20Play%20with%20Docker-red.svg?style=for-the-badge&logo=appveyor)](https://labs.play-with-docker.com/?stack=https://raw.githubusercontent.com/RawSanj/space-cloud/master/docker/docker-compose.yml#)
