# Getting Started with Space Cloud

Follow this guide to get started using `space-cloud` in a project from scratch.

> Note: If you instead want to play around with `space-cloud` to see what it can do, then check out the [Quick start](/docs/quick-start).

Since `space-cloud` bridges the gap between your frontend and database, it involves a 3-tier architecture as shown below:

<img src="https://spaceuptech.com/icons/space-cloud-basic.png"  alt="Basic Space Cloud architecture" />

Thus, using `space-cloud` in a project requires you to configure three parts:

- Database (Prerequisite)
- Space Cloud
- Frontend

## Step 1: Start your Database

`space-cloud` exposes realtime CRUD functionalities over any database of your choice. In order for that to work, you first need to have your database up and running. `space-cloud` supports the following databases as of now:

- Mongo DB
- MySQL and MySQL compatible databases (For eg. TiDB)
- Postgres and Postgres compatible databases (For eg. CockroachDB, Yugabyte etc.)

## Step 2: Download Space Cloud

Once you have your database up and running, you need to make sure that you have the latest version of the `space-cloud` binary on your machine. This binary is the server connecting to your database and creating the endpoints for it. You need to download a binary for your operating system or you could build it directly from its source code. You will need go version 1.11.2 or later to build it from the source.

Download the binary for your OS from here:

- [Mac](https://spaceuptech.com/downloads/darwin/space-cloud.zip)
- [Linux](https://spaceuptech.com/downloads/linux/space-cloud.zip)
- [Windows](https://spaceuptech.com/downloads/windows/space-cloud.zip)

You can unzip the compressed archive

**For Linux / Mac:** `unzip space-cloud.zip`

**For Windows:** Right click on the archive and select `extract here`.

Make the `space-cloud` binary executable and add it to your path.

**For Linux / Mac:** `chmod +x space-cloud`

## Step 3: Configure Space Cloud

Space Cloud needs a config file in order to function. The config file is used to load information like the database to be used, its connection string, security rules, etc.

You can use a sample config file from here [here](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/realtime-todo-app/config.yaml).

Or else you can create a bare minimum config file for a project using the following command from the folder containing `space-cloud`:

**For Linux / Mac:** `./space-cloud init`

**For Windows:** `space-cloud.exe init`

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

## Step 4: Start Space Cloud

You can start `space-cloud` with the following command.

**For Linux / Mac:** `./space-cloud run --config config.yaml`

**For Windows:** `space-cloud.exe run --config config.yaml`

## Next steps

As you have both the database and the `space-cloud` up and running, the next task is to set up your frontend app to use `space-cloud` and start building an app! Check out the language specific guides below to help you do this:

- [Javascript](/docs/getting-started/javascript) for web and Nodejs projects.
- [Java](/docs/getting-started/java) for Android and Java projects. (coming soon)
- [Python](/docs/getting-started/python) for Python projects.

<div class="btns-wrapper">
  <a href="/docs/quick-start/" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
