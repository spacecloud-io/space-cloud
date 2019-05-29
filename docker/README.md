
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
