Step 1: Install Docker

  Ensure that you have a recent versions of docker and docker-compose installed.

Step 2: Download space-cloud binary

  Download the latest binary from here:
  https://spaceuptech.com/downloads/linux/space-cloud.zip

  extract the space-cloud binary file to service/bin folder.

Step 3: Modify docker-compose.yaml
  look for sc01.example.com-0001 and replace it to your valid letsencrypt certificate path.

Step 4: Modify front/front-envoy.yaml
  look for sc01.example.com and replace it with your valid domain that's acquired letsencrypt certificates

Step 5: Modify admin and password of Space Cloud Admin portal from service-start.sh
  look for service-start.sh in service folder and modify the values

Step 6: Start all of our containers

  docker-compose up --build -d

Step 7: Login Space Cloud admin portal

  http://sc01.example.com/mission-control

    