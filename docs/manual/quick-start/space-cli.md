# Quick start (Space CLI)

This guide will help you get started with Space Cloud and Mongo DB on your local machine. It exposes complete functionality of Space Cloud.

In this guide I will walk you through installing the `space-cli` via npm. We will need this tool to start the `space-cloud` binary and Mongo DB. `space-cli` also hosts Space Console which is a visual tool to create and configure projects built with Space Cloud. I must also add that `space-cli` will be deploying our back end for us. However, it is recommended to deploy Space Cloud or any other tool for that matter with a proven deployment solution such as Kubernetes.

## Prerequisites
- [Docker](https://docs.docker.com/install/)
- [Node.js / NPM](https://nodejs.org/en/)

## Step 1: Install Space CLI
The first step is to install `space-cli`. `space-cli` is a command line utility which makes is very convenient to create and configure projects made using Space Cloud.
```
npm install -g space-cli
```

## Step 2: Create a new project
The first step is to create a project. A Space Cloud project is nothing but a YAML file which contains the configuration to start the `space-cloud` binary. The YAML file includes config such as the database to be used, it's connection string and security rules.

For this example we will create a simple todo app with MongoDB as ou database. We will start by using a sample config file. To do that run:
```
space-cli new --sample basic-todo-app
```
It will prompt you to put a project id. Let's keep is `todo-app` for now.

Hit enter.

And your done. `space-cli` has automatically created the config file for the todo-app. Feel free to explore the config file. You can read more about it to create your own apps right [here](/docs/config/overview).

## Step 3: Deploy the backend
The only task remaining for us to do is deploy our backend. `space-cli` comes with a neat API to deploy the entire backend along with the database using docker. You can deploy the backend using the following command:
```
space-cli deploy --local --config todo-app.yml
```

This function will load the YAML file and deploy MongoDB and `space-cloud` using docker. Note, you could use this command to deploy the database of your choice. `space-cli` goes through the config file and deploys all the dependencies for you.

> Note: It is recommended to deploy Space Cloud projects using the `cloud` option for production use cases. This command is used to push the project config to an already running cluster. For more notes on how to deployments check out the [deploy page](/docs/deploy).

Awesome!

Our backend is up and running.

That was quick wasn't it?

## Step 4: Try it out
Our back end is up and running. We have built a [basic todo app](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/basic-todo-app/index.html) using html and javascript which uses the backend you have just setup. Try signing up and adding some todos to test it out.

## Next Steps
Awesome! We just made an end-to-end app without writing a single line of back end code. The next step is to dive into the various Space Cloud modules or run some [sample apps](/docs/quick-start/sample-apps).
- [User Management](/docs/user-management)
- [Database](/docs/database) (For CRUD operations)
- [Realtime](/docs/real-time)
- [Functions](/docs/functions)

<div class="btns-wrapper">
  <a href="/docs/quick-start/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/quick-start/sample-apps" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>