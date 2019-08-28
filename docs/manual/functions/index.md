# Functions Mesh

The functions module is a means to extend the functionality provided by Space Cloud. Using this module you can write your own custom logic on the backend in the form of functions in any language. These functions can be called in a secure manner either directly from the frontend or from any other backend function you have written.

> **Note:** The functions you write run as long lived processes on the backend.

## Service mesh redefined!
The functions module redefines the traditional service mesh by asbtracting away service discovery, load balancing and all networking under the hood. It allows you to architect your microservices in the form of simple functions rather than end to end service. Thus, the functions module brings all the advantages of a microservice architecture at the ease of writing simple functions!

## What can I do with functions?

Since the functions module can technically allow you to write your own custom logic, the possibilities with it are endless. Here are a few ways in which you can use the functions module: 

- Integrate with third party services and APIs.
- Trigger custom logic on changes in database (For example, send a welcome email when a user completes signup).
- Extend security module of Space Cloud by providing custom validations (For example, allow a particular crud operation only if the current day is Sunday).
- Add a database support to Space Cloud.
- Run some ML algorithm periodically on data to update the model.

## How it works

As an user, you need to write functions on the backend. You can group related functions together as a `service`. Each service runs as a long lived process on the backend. In a single project, you can have multiple services in the language of your choice. Each service connects and subscribes to Space Cloud with its service name.

> **Note:** A service acts as a client to Space Cloud and not an upstream server.

The frontend or some other function can request to trigger a specific function by providing a service name, function name and params for that function. On receiving the request, Space Cloud first validates the request via the [security module](/docs/security). If the request can be made, it then makes an RPC (remote procedural call) on behalf of the client to the requested function. The function is executed with the params and returns a response to Space Cloud which is then returned to the client.

Space Cloud uses a broker with RPC semantics under the hood to heavy lift the scaling, networking and load balancing of RPC calls. Space Cloud runs a Nats server by default in the same process so that you don't have to run a broker. However you can run your own broker and configure Space Cloud to use that broker instead. As of now, only Nats is supported as a broker with the support of RabbitMQ and Kafka coming soon.

## Configure the functions module

Head over to the `Functions Mesh` section in Mission Control to configure the functions module.

### Enabling the functions module
Enable the switch in the upper right corner of `Functions Mesh` section to enable the functions module.

### Configuring services

Mission Control by default creates the configuration for a `default` service for you with a `default` function. The configuration of the `default` service/function is used when the configuration of requested service/function is not found.    

The configuration of a service looks like the following: 

```json
{
  "functions": {
    "func1": {
      "rule": {
        "rule": "allow"
      }
    }
  }
}
```

The above example exposes all function calls to `func1`. You can learn more about securing the functions module in depth over [here](/docs/security/functions).  

## Next steps

You can now see how to write the services on the backend and invoke them from the frontend.

<div class="btns-wrapper">
  <a href="/docs/file-storage/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/functions/service" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
