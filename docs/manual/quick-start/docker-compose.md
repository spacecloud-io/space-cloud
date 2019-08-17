# Quick start (Using Docker Compose)

This guide will help you run a local development setup that sets up both the `space-cloud` and MongoDB. It will guide you through exploring the Space Cloud APIs on MongoDB without having to set up any project.

If you instead want to start a project from scratch using `space-cloud`, then check out the [getting started](/docs/getting-started) guide.

> **Note:** MongoDB is not a dependency of Space Cloud. We are using MongoDB in this guide for ease of use because of it's schemaless nature.

## Prerequisites

- [Docker](https://docs.docker.com/install/) >= 18.02.0
- [Docker Compose](https://docs.docker.com/compose/install/) >= 1.20.0


## Step 1: Get the docker-compose file

The [spaceuptech/space-cloud/install-manifests](https://github.com/spaceuptech/space-cloud/tree/master/install-manifests) repo contains all installation manifests required to deploy Space Cloud anywhere. Get the docker compose file from there:

```bash
wget https://raw.githubusercontent.com/spaceuptech/space-cloud/master/install-manifests/quick-start/docker-compose/mongo/docker-compose.yaml
```

## Step 2: Run Space Cloud & MongoDB

```bash
docker-compose up -d
```

Check if the containers are running:
```bash
docker ps
```

## Step 3: Configure Space Cloud

If you exec into docker container of Space Cloud, you can see a `config.yaml` file and a `raft-store` folder would have been generated in the home directory.

Space Cloud needs this config file in order to function. The config file is used to load information like the database to be used, its connection string, security rules, etc. 

Space Cloud has it's own Mission Control (admin UI) to configure all of this in an easy way. 

> **Note:** All changes to the config of `space-cloud` has to be done through the Mission Control only. Changes made manually to the config file will get overwritten. 


### Open Mission Control

Head over to `http://localhost:4122/mission-control` to open Mission Control.

> **Note:** Replace `localhost` with the address of your Space Cloud if you are not running it locally. 

### Creating a project
Click on `Create a Project` button. 

Give a `name` to your project. MongoDB will be selected as your database by default. Keep it as it is for this guide.

Hit `Next` to create the project. On creation of project you will be directed to the project overview screen. 

### Configuring DB config

Head over to the `Database` section. 

Copy paste `mongodb://mongo:27017` in the `connection string` input.

Hit `Save` button. That's all what is required to configure Space Cloud for this guide!

## Step 4: Try it out

Our backend is up and running, configured to expose APIs on MongoDB. Time to explore it's awesome powers. 

## Space Cloud Explorer 

Head over to the `Explorer` section. 

The Explorer is a tool in Mission Control that lets you try Space Cloud APIs without actually setting up any frontend or backend project. It directly lets you run javascript APIs of `space-cloud` against itself.  

> **Note:** The `api` object and certain helpers like `and`, `or` and `cond` for generating where clauses are available to all code snippets you run through the Explorer.

### Inserting data

Copy paste the following code in the Explorer and hit apply to insert multiple todos:

```js
const db = api.Mongo()
const docs = [
  {_id: "1", text: "Star Space Cloud on Github", date: new Date()},
  {_id: "2", text: "Follow us on Twitter", date: new Date()},
  {_id: "3", text: "Spread the love!", date: new Date()}
]

db.insert("todos").docs(docs).apply()
```

On successful insert, you should be able to see the `status` as `200` which means the documents were inserted in the database.

### Querying all documents back
To retrieve the documents from MongoDB using Space cloud, copy paste the following code in the Explorer and hit apply:

```js
const db = api.Mongo()

db.get("todos").apply()
```

You should be able to see the `status` as `200` and the `result` with the documents you inserted in the previous step:
```json
{
  "result": [
    {
      "_id": "1",
      "date": "2019-08-03T03:24:43.641Z",
      "text": "Star Space Cloud on Github"
    },
    {
      "_id": "2",
      "date": "2019-08-03T03:24:43.641Z",
      "text": "Follow us on Twitter"
    },
    {
      "_id": "3",
      "date": "2019-08-03T03:24:43.641Z",
      "text": "Spread the love!"
    }
  ]
}
```

### Querying single document
To retrieve the document with `_id` equals to `2`, copy paste the following code in the Explorer and hit apply:

```js
const db = api.Mongo()

db.getOne("todos").where(cond("_id", "==", "2")).apply()
```

You should be able to see the `status` as `200` and the following `result`:
```json
{
  "result": {
      "_id": "2",
      "date": "2019-08-03T03:24:43.641Z",
      "text": "Follow us on Twitter"
    }
}
```


## Next Steps

Awesome! We just performed few CRUD operations on MongoDB without writing a single line of backend code. The next step is to dive into the various Space Cloud modules.

- Perform CRUD operations using [Database](/docs/database/) module
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
