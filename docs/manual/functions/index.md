# Functions

The functions module is a means to extend the functionality provided by Space Cloud. Using this module you can write your own custom logic on the backend in the form of functions. The clients can call this function directly from the frontend.

> **Note:** The functions you write on the backend run as microservices on the backend.

## How it works

As an user, you need to write functions on the backend. In a single project you can write multiple functions. A collection of functions is called an `engine`. So in other words, you can write engines in the language of your choice which will have multiple functions registered to it.

All engines connect to [nats](https://nats.io) and subscribe to a subject name which is a derivative of the engine name. Nats is a pub-sub network which load balances requests between the engines.

The `space-cloud` server acts as an api-gateway which connects to nats as well. The request from the front end will be received by Space Cloud. Space Cloud would then publish it on nats, receive the response, and send it to the client.

## Enable the functions module

The config pertaining to functions module can be found inside the `functions` key under the `modules` object. Here's the snippet:

```yaml
modules:
  functions:
    enabled: true
    nats: nats://localhost:4222
  # Config for other modules go here
```

All you need to do is set the `enabled` field to true and provide the connection string to nats.

## Next steps

You can now see how to write the engines on the backend and invoke them from the frontend.

<div class="btns-wrapper">
  <a href="/docs/file-storage/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/functions/engine" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
