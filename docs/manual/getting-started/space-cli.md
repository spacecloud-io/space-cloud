# Quick start (Space CLI)

This guide will help you get started with Space Cloud and Mongo DB on your local machine. It exposes complete functionality of Space Cloud.

In this guide I will walk you through installing the `space-cli` via npm. We will need this tool to start the `space-cloud` binary and Mongo DB. `space-cli` will also host the Space Cloud Console which is a visual tool to create and configure projects built with Space Cloud. I must also add that `space-cli` will be deploying our back end for us. However, it is recommended to deploy Space Cloud or any other tool for that matter with a proven deployment solution such as Kubernetes.

## Prerequisites
- [Docker](https://docs.docker.com/install/)
- [Node.js / NPM](https://nodejs.org/en/)

## Step 1: Install Space CLI
The first step is to install `space-cli`. `space-cli` is a command line utility which makes is very convenient to create and configure projects made using Space Cloud.
```
npm install -g space-cli
```

## Step 2: Start the Space Cloud Console
Now we need to start the Space Cloud Console. Console is the User Interface to interact with `space-cloud`. `space-cli` is responsible to automatically fetch the latest version of the Console to make sure you always are up-to-date.
```
space-cli start
```
This command will download the latest version of Space Cloud Console and make it available on [http://localhost:8080/console](http://localhost:8080/console)

// Space Cloud Console image goes here

## Step 3: Create a project
Let's create out first project. You can create your own project fro scratch using the `Add Project` button on the home screen. For now, let's use the Todo App template to create a project with the default configuration.

// Create project model image goes here

You can change the name of the project if required. Let the database be Mongo DB and hit create. 

## Step 4: Deploy the back end
We have created our first Space Cloud project. You must be welcomed by the `projects` page.

// Space console projects page goes here

 On the left you have all the modules which Space Cloud has to offer. The top bar contains a few useful tips and resources which could come in handy. The upper right `deploy` button deploys the configuration to `space-cloud`.

 Hit the `deploy` button and choose the `on laptop` option. This will instruct `space-cli` to start `space-cloud` binary along with the database.

> Note: It is recommended to deploy Space Cloud projects using the `cloud` option for production use cases. This command is used to push the project config to an already running cluster. For more notes on how to deployments check out the [deploy page](https://spaceuptech.com/docs/deploy).

// Space console image with deploy popup goes here.

## Step 5: Try it out
Our back end is up and running. We have built a [realtime to-do app](https://spaceuptech.com/downloads/todo-app.html) which used the back end you have just setup. Try signing up and adding some to-dos from two browsers simultaneously.

// Screen with two todo apps running simultaneously goes here

## Next Steps
Awesome! We just made a realtime app without writing a single line of back end code. The next step is to dive into the various Space Cloud modules or run some [sample apps](https://spaceuptech.com/docs/getting-started/sample-apps).
- [User Management](https://spaceuptech.com/docs/user-management)
- [Database](https://spaceuptech.com/docs/database) (For CRUD operations)
- [Realtime](https://spaceuptech.com/docs/real-time)
- [Functions](https://spaceuptech.com/docs/functions)

<< [previous](https://spaceuptech.com/docs/getting-started) | [next](https://spaceuptech.com/docs/getting-started/sample-apps) >>
