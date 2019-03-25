# Quick start (Manual)

This guide will help you get started with Space Cloud and Mongo DB on your local machine. It exposes complete functionality of Space Cloud.

In this guide I will walk you through how to develop a realtime todo app using Space Cloud. We'll be deploying the `space-cloud` binary manually.

## Prerequisites
- [MongoDB Database](https://docs.mongodb.com/manual/installation/)

> Note: MongoDB is not a dependency of Space Cloud. The sample app in this quick start uses MongoDB as its database.

## Step 1: Download `space-cloud`
The first step is to download the `space-cloud` binary. This binary is the server creating the endpoints and connecting to your database.  You need to download binary for your operating system or you could build it directly from its source code. You will need go version 1.11.2 or later to build it from source.

Download the binary for your OS from here:
- [Mac](https://spaceuptech.com/downloads/darwin/space-cloud.zip)
- [Linux](https://spaceuptech.com/downloads/linux/space-cloud.zip)
- [Windows](https://spaceuptech.com/downloads/windows/space-cloud.zip)

You can unzip the compressed archive

**For Linux / Mac:** `unzip space-cloud.zip && chmod +x space-cloud`

**For Windows:**      Right click on the archive and select `extract here`.

## Step 2: Download the sample config file
Space Cloud needs a config file in order to function. The config file is used to load information like the database to be used, its connection string, security rules, etc. You can find the config used for our todo app [here](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/realtime-todo-app/config.yaml).

This is how the config file looks like.

```yaml
---
id: space-cloud
secret: some-secret
modules:
  crud:
    mongo:
      conn: mongodb://localhost:27017
      isPrimary: true
      collections:
        todos:
          isRealtimeEnabled: true
          rules:
            create:
              rule: allow
            read:
              rule: allow
            update:
              rule: allow
            delete:
              rule: allow
  auth:
    email:
      enabled: true
  faas:
    enabled: false
    nats: nats://localhost:4222
  realtime:
    enabled: true
    kafka: localhost
  fileStore:
    enabled: false
    storeType: local
    conn: ./
    rules:
      rule1:
        prefix: /
        rule:
          create:
            rule: allow
          read:
            rule: allow
          delete:
            rule: allow
```

Quickly going through it, `id` is the project name. `secret` is the secret key used for signing and parsing JWT tokens. All the configuration for individual modules goes under the `modules` key. Currently, `crud`, `auth` (user management), `faas` (functions), `realtime` and `fileStore` are supported.

**Note:** When you are starting a project from scratch, you can run `./space-cloud init` (on Linux / Mac) or `space-cloud.exe init` (on Windows) to create a bare minimum config file.

## Step 3: Start Space Cloud
You can start `space-cloud` with the following command. Make sure MongoDB is running before this step.

**For Linux / Mac:** `./space-cloud run --config config.yaml`

**For Windows:** `space-cloud.exe run --config config.yaml`

That's it. Your backend is up and running!

That was quick wasn't it?

## Step 4: Try it out
Our back end is up and running. Time to show off it's awesome powers. We have built a [realtime todo app](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/realtime-todo-app/index.html) using html and javascript which uses the backend you have just setup. 

Open it in two different windows by double clicking the html file twice, login into both and then try adding some todos to see the magic.

## Next Steps
Awesome! We just made a realtime app without writing a single line of back end code. The next step is to dive into the various Space Cloud modules or run some [sample apps](/docs/quick-start/sample-apps).
- [User Management](/docs/user-management)
- [Database](/docs/database) (For CRUD operations)
- [Realtime](/docs/realtime)
- [Functions](/docs/functions)

<div class="btns-wrapper">
  <a href="/docs/quick-start/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/quick-start/sample-apps" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>