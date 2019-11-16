Step 1: Install Docker<br/>
<br/>
  Ensure that you have a recent versions of docker and docker-compose installed.<br/>
<br/><br/>


Step 2: Download space-cloud binary<br/>
<br/>
  Download the latest binary from here:<br/>
  https://spaceuptech.com/downloads/linux/space-cloud.zip<br/>
<br/>
  extract the space-cloud binary file to service/bin folder.<br/>
<br/><br/>


Step 3: Start all of our containers<br/>
<br/>
  docker-compose up --build -d<br/>
<br/><br/>


Step 4: Login Space Cloud admin portal<br/>
<br/>
  http://localhost:8000/mission-control<br/>
  login as username: admin, password: admin<br/>
<br/><br/>


Step 5: to scale service1<br/>
<br/>
  docker-compose scale service1=3<br/>
  To learn more about scaling, https://www.envoyproxy.io/docs/envoy/latest/start/sandboxes/front_proxy<br/>