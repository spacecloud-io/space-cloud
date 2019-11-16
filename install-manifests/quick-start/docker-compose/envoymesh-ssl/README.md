Step 1: Install Docker<br>
<br>
  Ensure that you have a recent versions of docker and docker-compose installed.<br>
<br><br>

Step 2: Download space-cloud binary<br>
<br>
  Download the latest binary from here:<br>
  https://spaceuptech.com/downloads/linux/space-cloud.zip<br>
<br>
  extract the space-cloud binary file to service/bin folder.<br>
<br><br>

Step 3: Modify docker-compose.yaml<br>
  look for sc01.example.com-0001 and replace it to your valid letsencrypt certificate path.<br>
<br><br>

Step 4: Modify front/front-envoy.yaml<br>
  look for sc01.example.com and replace it with your valid domain that's acquired letsencrypt certificates<br>
<br><br>

Step 5: Modify admin and password of Space Cloud Admin portal from service-start.sh<br>
  look for service-start.sh in service folder and modify the values<br>
<br><br>

Step 6: Start all of our containers<br>
<br>
  docker-compose up --build -d<br>
<br><br>

Step 7: Login Space Cloud admin portal<br>
<br>
  https://sc01.example.com/mission-control<br>
<br><br>

Step 8: to scale service1<br/>
<br/>
  docker-compose scale service1=3<br/>
  To learn more about scaling, https://www.envoyproxy.io/docs/envoy/latest/start/sandboxes/front_proxy<br/>