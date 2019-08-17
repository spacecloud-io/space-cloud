# Quick Start

This guide will help you quickly get started with _Space Functions_. As mentioned [earlier](/docs/functions/overview), the functions module is a means to extend the functionality provided by Space Cloud and write custom business logic.

In this guide we will create a function to carry out complex arithmetic operation. To be more precise, we will be computing the addition of two numbers and returning the result. 

## Prerequisites

- [Nodejs](https://nodejs.org/en/download/)
- [Docker](https://docs.docker.com/install/)

> **Note:** There is no restriction on which language you choose to write your functions. We are using Nodejs in this guide just because it's relatively easy to get started with.

## Step 1: Start Space Cloud

You can start Space Cloud in a docker container by simply runnning the following command:

```bash
docker run -d -p 4122:4122 -p 4124:4124 --name space-cloud \
  -e DEV=true \
  spaceuptech/space-cloud:latest
```

> **Note:** You can also start space-cloud without docker by following [this guide](/docs/quick-start/manual).

## Step 2: Configure Space Cloud

On running `space-cloud`, a `config.yaml` file and a `raft-store` folder gets generated in the directory from where you run `space-cloud`.

Space Cloud needs this config file in order to function. The config file is used to load information like the database to be used, its connection string, security rules, etc. 

Space Cloud has it's own Mission Control (admin UI) to configure all of this in an easy way. 

> **Note:** All changes to the config of `space-cloud` has to be done through the Mission Control only. Changes made manually to the config file will get overwritten. 


### Open Mission Control

Head over to `http://localhost:4122/mission-control` to open Mission Control.

> **Note:** Replace `localhost` with the address of your Space Cloud if you are not running it locally. 

### Creating a project

Click on `Create a Project` button. 

Give a `name` to your project. For this guide we'll use the name `my-adder`. 

It doesn't matter which database you choose here since we'll be disabling it anyways.

Hit `Next` to create the project. On creation of project you will be directed to the project overview screen. The functions module is enabled for us by default.

### Disable the Database Module

Head over to the `Database` section.

In the upper right corner, you'll find a switch to disable the database. Since we don't need to use the database in this guide, we can disable it and hit `Save`. Otherwise, you can leave it the way it is.

## Step 3: Setup our NPM Project

Our function is going to be a npm project. Then we will have to create an `index.js` file which will contain our super complex function.

- First create a folder which will be our working directory.
- Run `npm init -y` to create an npm project
- Install `space-api` which we'll use to register the function with Space Cloud by running `npm install --save space-api`

Our working environment is set up.

> **Note:** You can install any packages when creating a Space Cloud function. Since functions run as longed lived microservices, we impose no restrictions on your dependencies and execution environment.

## Step 4: Write our function

Create a file named `index.js`. Our entire code will reside in this function itself.

In Space Functions, we have a concept of `services` which is nothing but a collection of functions. All instances of the same service is expected to have identical functions. Space Cloud will take care of balancing load between all services.

We'll be naming our service `arithmetic` and our function `sum`.

Copy the following content in our file:

```js
const { API } = require('space-api');           // Import the space-api library

const projectId = 'my-adder';                   // The name of our project
const spaceCloudURL = 'http://localhost:4122';  // The url of Space Cloud

// Create an space api object
const api = new API(projectId, spaceCloudURL);

// Make a service object
const service = api.Service('arithmetic');

// Register function to a service
service.registerFunc('sum', (params, auth, cb) => {
  
  // params - contains the payload send by the client.
  // auth   - contains the JWT claims of the client making the request.
  // cb     - is the functions used to give back a response to the client
  
  // In our use case:
  // params = { num1: SOME_NUMBER, num2: ANOTHER_NUMBER }

  // Lets add the two numbers
  const sum = params.num1 + params.num2;

  // Return the response to the clinet
  cb('response', { sum: sum });
});

// Start the service
service.start() 
```

Going through the code will make it clear whats happening. 

First we created an api object which essentially provides the name of the project and location of Space Cloud. This is an important point. **The service is a client to Space Cloud**. It does not run an HTTP server or anything of that sort. 

This also means that the service can run anywhere (behind a firewall, in another cloud, etc.) as long as it can connect to a Space Cloud cluster. Space Cloud pushes requests down to the service via a **bi-directional link**. No need to worry about service discovery or anything of that sort.

Next we create a service and register a function on that service.

Communication with Space Cloud begin once we execute `service.start()`.

## Step 5: Start the service

Simply run `node index.js`. Pretty straight forward.

## Step 6: Try it out

Our service is up and running. Time to explore it's awesome powers. 

### Space Cloud Explorer 

Open mission control and head over to the `Explorer` section. 

The Explorer is a tool in Mission Control that lets you try Space Cloud APIs without actually setting up any frontend or backend project. It directly lets you run javascript APIs of `space-cloud` against itself.  

> **Note:** The `api` object and certain helpers like `and`, `or` and `cond` for generating where clauses are available to all code snippets you run through the Explorer.

### Call the function

Now simply call the function. Copy the following snippet in the text area of `Explorer` and hit `Apply`.

```js
const serviceName = 'arithmetic'; // Our service name
const funcName = 'sum';           // Our function name

// The numbers we want to add
const params = { num1: 10, num2: 12 };

// Simply call the function
api.call(serviceName, funcName, params);
```

You should be able to see the `status` as `200` and the `result` containing the object returned to us by the function.

## Next Steps

Awesome! We just experienced the power of Space Functions. Super proud to be solving a complex mathematical operation without having to worry about networking and stuff.

It might not be possible to come up with something as complex as additon in the real world, but you can definately try something more simple like querying your tensorflow models using Space Functions.

The next step is to dive into the various Space Cloud modules.

- Dive deeper into the [Functions](https://spaceuptech.com/docs/functions/service) module
- [Secure](https://spaceuptech.com/docs/security/functions) your functions

<div class="btns-wrapper">
  <a href="/docs/functions/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/functions/service" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>