# Writing custom logic

You can easily extend Space Cloud by writing your custom logic on the backend in the form of simple functions. These functions run as a microservice on the backend. This is how you write a simple function using the `engine-api` -

```go
import "spaceuptech.com/space-engine-go/engine"

// Function to be registered
func myFunc(params engine.M, auth engine.M, cb engine.CallBack) {
    log.Println("Params", params, "Auth", auth)
    // Do something

    // Call the callback
    cb(engine.TypeResponse, engine.M{"ack": true})
}

// Create an instance of engine
myEngine, err := engine.Init("my-engine", "")
if err != nil {
    log.Println("Err", err)
    return
}

// Register function
myEngine.RegisterFunc("my-func", myFunc)

// Start engine
myEngine.Start()

// Call function of some other engine
res, err := myEngine.Call("some-engine", "some-func", engine.M{"msg": "space-engine-go is awesome!"}, 1000)
```

Use `engine.Init` to initialize an instance of an `engine`. An `engine` can harbour multiple functions which can be invoked by frontend. The `engine.Init` function takes two parameters **engineName** and **url** which are as follows:
- **engineName:** Name of the engine. Uniquely identifies an engine
- **url:** Connection string of nats. Pass "" to use the default nats connection string

You can register a function that you have written to an `engine` by calling `RegisterFunc` on an `engine`. `RegisterFunc` takes a **name** and a **func** which are as follows:
- **name:** Name of the function. Uniquely identifies a function within an engine
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