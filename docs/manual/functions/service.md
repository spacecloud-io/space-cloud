# Writing custom logic

You can easily extend Space Cloud by writing your custom logic on the backend in the form of simple functions. These functions run as a microservice on the backend. This is how you write a simple service -

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#service-js">Javascript</a></li>
      <li class="tab col s2"><a href="#service-java">Java</a></li>
      <li class="tab col s2"><a href="#service-python">Python</a></li>
      <li class="tab col s2"><a href="#service-golang">Golang</a></li>
    </ul>
  </div>
  <div id="service-js" class="col s12" style="padding:0">
    <pre>
      <code>
const { API, cond } = require('space-api');

const api = new API('my-app', 'http://localhost:4122');

// Make a service
const service = api.Service('service-name');

// Register function to a service
service.registerFunc('function-name', (params, auth, cb) => {
  // Your custom logic goes here

  // Response to be returned to client
  const res = { ack: true, message: 'Functions mesh is amazing!' }
  cb('response', res)
})

// Start the service
service.start()      
      </code>
    </pre>
  </div>
  <div id="service-golang" class="col s12" style="padding:0">
    <pre>
      <code>
import (
	"github.com/spaceuptech/space-api-go/api"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/service"
	"fmt"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	service := api.Service("service")
	service.RegisterFunc("echo_func", Echo)
	service.Start()
	
}

func Echo(params, auth *model.Message, fn service.CallBackFunction) {
	var i interface{}
	params.Unmarshal(&i)
	fn("response", i)
}
      </code>
    </pre>
  </div>
  <div id="service-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
API api = new API("books-app", "localhost", 4124);
Service service = api.service("service");
service.registerFunc("echo_func", new ServiceFunction() {
    @Override
    public void onInvocation(Message params, Message auth, ReturnCallback cb) {
        cb.send("response", params.getValue(Object.class));
    }
});
service.start(); // Is blocking
      </code>
    </pre>
  </div>
 <div id="service-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

api = API('books-app', 'localhost:4124')


def my_func(params, auth, cb):  # Function to be registered
    print("Params", params, "Auth", auth)

    # Do Something
    cb('response', {"ack": True})


service = api.service('service')  # Create an instance of service
service.register_func('my_func', my_func)  # Register function
service.start()  # Start service (This is a blocking call)

api.close()
      </code>
    </pre>
  </div>
</div>

Use `api.Service` to initialize an instance of an `service`. The `api.Service` function takes only one parameter - **serviceName** which uniquely identifies the service. 

A `service` can harbour multiple functions which can be invoked by client. `service.registerFunc` is used to register a function to a service. The `registerFunc` method takes two parameters:
- **funcName:** Name of the function which uniquely identifies a function within a service
- **func:** The function to be executed

## Writing a function

Any registered function gets three arguments during execution when triggered by client as follows:   

- **params:** The params object sent by the client.
- **auth:** Auth object (consists the claims of JWT Token)
- **cb:** Callback function used to return the response back to the client

### Send JSON object back to client
To send JSON object as a response back to client, call the `cb` function with type as `response` and the second parameter being the response object: 
```js
// Any object that you want to send as response
const response = { ack: true, message: 'I love functions mesh!' } 
cb('response', response)
```

### Render HTML page back to client
Coming soon!

## Next steps

Great! So now you know how to write custom logic on backend. Let's checkout how to invoke it from the frontend.

<div class="btns-wrapper">
  <a href="/docs/functions/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/functions/client" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div> 