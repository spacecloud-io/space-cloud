# Access custom logic
You can call a function running on the backend (written via functions module of Space Cloud) on frontend by simply calling `api.call` on frontend. Here's a code snippet showing how to do it:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#client-js">Javascript</a></li>
      <li class="tab col s2"><a href="#client-java">Java</a></li>
      <li class="tab col s2"><a href="#client-python">Python</a></li>
    </ul>
  </div>
  <div id="client-js" class="col s12" style="padding:0">
    <pre>
      <code>
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Call a function running 'my-func' of 'my-service' running on backend
api.call('my-service', 'my-func', { msg: 'Space Cloud is awesome!' }, 1000)
  .then(res => {
    if (res.status === 200) {
      console.log('Response: ', res.data)
    }
  }).catch(ex => {
    // Exception occured while processing request
  })
      </code>
    </pre>
  </div>
  <div id="client-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
 <div id="client-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

api = API("books-app", "localhost:8081")
# Call a function, 'my-func' of 'my-engine' running on backend
response = api.call('my-engine', 'my-func', {"msg": 'Space Cloud is awesome!'}, 1000)
if response.status == 200:
    print(response.result)
else:
    print(response.error)
      </code>
    </pre>
  </div>
</div>

The `call` function takes four arguments which are as follows:
- **serviceName** - Name of the service
- **funcName** - Name of the function
- **params** - An object that can contain any data that you want to pass to the function on backend
- **timeOut** - Timeout in seconds

As you would have noticed, the above function is asynchronous in nature. 

## Response

A response object sent by the server contains the **status** and **data** fields explained below:

**status:** Number describing the status of the operation. Following values are possible:
- 200 - Operation was successful
- 500 - Internal server error

**data:** Object returned by the function.

## Next steps

Now you know the basics of all the modules. So let's take a deep dive at securing your app! 

<div class="btns-wrapper">
  <a href="/docs/functions/service" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/security/overview" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div> 