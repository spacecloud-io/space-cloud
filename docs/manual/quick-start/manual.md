# Quick start (Manual)

This guide will help you get started with Space Cloud quickly on your local machine. You will go through how to develop a realtime todo app using Space Cloud and MongoDB. We'll be deploying the `space-cloud` binary manually.

> Note: If you instead want to start a project from scratch using `space-cloud`, then check out the [getting started](/docs/getting-started) guide.

## Prerequisites

- [MongoDB Database](https://docs.mongodb.com/manual/installation/)

> Note: MongoDB is not a dependency of Space Cloud. The sample app in this quick start uses MongoDB as its database.

## Step 1: Download `space-cloud`

The first step is to download the `space-cloud` binary. This binary is the server creating the endpoints and connecting to your database. You need to download binary for your operating system or you could build it directly from its source code. You will need go version 1.11.2 or later to build it from source.

Download the binary for your OS from here:

- [Mac](https://spaceuptech.com/downloads/darwin/space-cloud.zip)
- [Linux](https://spaceuptech.com/downloads/linux/space-cloud.zip)
- [Windows](https://spaceuptech.com/downloads/windows/space-cloud.zip)

You can unzip the compressed archive

**For Linux / Mac:** `unzip space-cloud.zip && chmod +x space-cloud`

**For Windows:** Right click on the archive and select `extract here`.

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
  functions:
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

Quickly going through it, `id` is the project name. `secret` is the secret key used for signing and parsing JWT tokens. All the configuration for individual modules goes under the `modules` key. Currently, `crud`, `auth` (user management), `functions` (functions), `realtime` and `fileStore` are supported.

> Note: The in-depth configurations of various modules are explained in their corresponding sections.

## Step 3: Start Space Cloud

You can start `space-cloud` with the following command. Make sure MongoDB is running before this step.

**For Linux / Mac:** `./space-cloud run --config config.yaml`

**For Windows:** `space-cloud.exe run --config config.yaml`

That's it. Your backend is up and running!

That was quick wasn't it?

## Step 4: Try it out

Our backend is up and running. Time to show off it's awesome powers. We have built a [realtime todo app](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/realtime-todo-app/index.html) using html and javascript which uses the backend you have just setup.

Open it in two different windows by double clicking the html file twice, login into both and then try adding some todos to see the magic.

## Next Steps

Awesome! We just made a realtime app without writing a single line of backend code. The next step is to dive into the various Space Cloud modules or run some [sample apps](/docs/quick-start/sample-apps).

- Perform CRUD operations using [Database](/docs/database/) module
- [Realtime](/docs/realtime/) data sync across all devices
- Manage files with ease using [File Management](/docs/file-storage) module
- Allow users to sign-in into your app using [User management](/docs/user-management) module
- Write custom logic at backend using [Functions](/docs/functions/) module
- [Secure](/docs/security) your apps

<div class="btns-wrapper">
  <a href="/docs/quick-start/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/quick-start/sample-apps" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
