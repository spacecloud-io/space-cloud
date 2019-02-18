# Quick start (Manual)

This guide will help you get started with Space Cloud and Mongo DB on your local machine. It exposes complete functionality of Space Cloud.

In this guide I will walk you through how to develop a todo app using Space Cloud. We'll be deploying the `space-cloud` binary manually. The recommended way to get started is using the [`space-cli`](/docs/getting-started/space-cli).

## Prerequisites
- [MongoDB Database](https://docs.mongodb.com/manual/installation/)

## Step 1: Download `space-cloud`
You need to download the `space-cloud` binary for your operating system or you could build it directly from its source code. You need go version 1.11.2 or later to build it from source.

Download the binary for your OS from here:
- [Mac](https://spaceuptech.com/downloads/darwin/space-cloud.zip)
- [Linux](https://spaceuptech.com/downloads/linux/space-cloud.zip)
- [Windows](https://spaceuptech.com/downloads/windows/space-cloud.zip)

Make the `space-cloud` binary executable and add it to you `PATH`.


## Step 2: Download the sample config file
Space Cloud needs a config file in order to function properly. It relies on the config file to load information like the database to be used, it's connection string, security rules, etc. 

You can find a sample config for the todo app [here](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/basic-todo-app/config.yaml). Feel free to explore the file.

## Step 3: Start Space Cloud
You can start `space-cloud` with the following command. Make sure mongo db is running before this step
```
space-cloud run --config config.yml
```

That's it. Your backend is up and running!

That was quick wasn't it?

## Step 4: Try it out
Our back end is up and running. We have built a [basic todo app](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/basic-todo-app/index.html) using html and javascript which uses the backend you have just setup. Try signing up and adding some todos to test it out.

## Next Steps
Awesome! We just made an end-to-end app without writing a single line of back end code. The next step is to dive into the various Space Cloud modules or run some [sample apps](https://spaceuptech.com/docs/getting-started/sample-apps).
- [User Management](https://spaceuptech.com/docs/user-management)
- [Database](https://spaceuptech.com/docs/database) (For CRUD operations)
- [Realtime](https://spaceuptech.com/docs/real-time)
- [Functions](https://spaceuptech.com/docs/functions)

<div class="btns-wrapper">
  <a href="/docs/getting-started/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/getting-started/sample-apps" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>