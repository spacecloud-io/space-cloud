# Writing custom logic

You can easily extend Space Cloud by writing your custom logic on the backend in the form of simple functions. These functions run as a microservice on the backend. This is how you write a simple function using the `service-api` -

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#service-js">Javascript</a></li>
      <li class="tab col s2"><a href="#service-go">Go</a></li>
      <li class="tab col s2"><a href="#service-java">Java</a></li>
      <li class="tab col s2"><a href="#service-python">Python</a></li>
    </ul>
  </div>
  <div id="service-js" class="col s12" style="padding:0">
    <pre>
      <code>
const service = require('space-service-node');

const service = new service('my-service');

service.registerFunc('my-func', (params, auth, cb) => {
  console.log('Params:', params, 'Auth', auth)
  // Do something

  const res = { ack: true, message: 'Function as a Service is Awesome!' }
  cb('response', res)
})
      </code>
    </pre>
  </div>
  <div id="service-go" class="col s12" style="padding:0">
    <pre>
      <code>
import "spaceuptech.com/space-service-go/service"

// Function to be registered
func myFunc(params service.M, auth service.M, cb service.CallBack) {
    log.Println("Params", params, "Auth", auth)
    // Do something

    // Call the callback
    cb(service.TypeResponse, service.M{"ack": true})
}

// Create an instance of service
myservice, err := service.Init("my-service", "")
if err != nil {
    log.Println("Err", err)
    return
}

// Register function
myservice.RegisterFunc("my-func", myFunc)

// Start service
myservice.Start()

// Call function of some other service
res, err := myservice.Call("some-service", "some-func", service.M{"msg": "space-service-go is awesome!"}, 1000)
      </code>
    </pre>
  </div>
  <div id="service-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
 <div id="service-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
# Python client coming soon!
      </code>
    </pre>
  </div>
</div>

Use `service.Init` to initialize an instance of an `service`. An `service` can harbour multiple functions which can be invoked by frontend. The `service.Init` function takes two parameters **serviceName** and **url** which are as follows:
- **serviceName:** Name of the service. Uniquely identifies an service
- **url:** Connection string of nats. Pass "" to use the default nats connection string

You can register a function that you have written to an `service` by calling `RegisterFunc` on an `service`. `RegisterFunc` takes a **name** and a **func** which are as follows:
- **name:** Name of the function. Uniquely identifies a function within an service
- **func:** Function that comprises of the custom logic

`func` is the function that will comprise of the custome logic that you want. It can be invoked by the client as and when required. The function takes 3 parameters as decsribed below:
- **params:** The params sent by the client
- **auth:** Auth object
- **cb:** Callback function used to return the response back to the client

## Next steps

Great! So now you know how to write custom logic on backend . Let's checkout how to invoke it from the frontend.

<div class="btns-wrapper">
  <a href="/docs/functions/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/functions/client" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div> 