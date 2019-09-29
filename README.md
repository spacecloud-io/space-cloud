<p align="center"><a href="https://www.spaceuptech.com"><img src="https://spaceuptech.com/icons/space-cloud-block-diagram1.png" alt="Space Cloud"></a></p>

<h3 align="center">
	Instant Realtime APIs on any database
</h3>

<p align="center">
	<strong>
		<a href="https://spaceuptech.com/">Website</a>
		•
		<a href="https://docs.spaceuptech.com/">Docs</a>
		•
		<a href="https://discord.gg/ypXEEBr">Support</a>
	</strong>
</p>
<p align="center">
    <a href="https://discord.gg/ypXEEBr"><img src="https://img.shields.io/badge/chat-discord-brightgreen.svg?logo=discord&%20style=flat"></a>
    <br/>
		<a href="https://goreportcard.com/report/github.com/spaceuptech/space-cloud">
		<img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/spaceuptech/space-cloud">
	  </a>
    <a href="https://opensource.org/licenses/Apache-2.0"><img
		alt="Apache License"
		src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"></a>
    <a href="https://twitter.com/intent/follow?screen_name=spaceuptech"><img src="https://img.shields.io/badge/                 Follow-spaceuptech-blue.svg?style=flat&logo=twitter"></a>
</p>

Space Cloud is a cloud native backend server that provides **instant, realtime APIs on any database**, with [event triggers](https://docs.spaceuptech.com/advanced/event-triggers) and [space functions](https://docs.spaceuptech.com/essentials/custom-logic) for custom business logic.

Space Cloud helps you build modern applications without having to write any backend code in most cases.

It provides **GraphQL** and **REST** APIs which can be consumed directly by your frontend in a [secure manner](https://docs.spaceuptech.com/auth).


## Features 

View complete feature set [here](https://docs.spaceuptech.com/getting-started/introduction/features).

- **Powerful CRUD**: Flexible queries, transactions and cross database joins
- **Realtime**: Make live queries to your database
- **File storage**: Upload/download files to scalable file stores (eg: Amazon S3, Google Cloud Storage)
- **Extensible**: Write custom business logic in form of simple functions
- **Event driven**: Trigger webhooks or serverless functions on database or file storage events
- **Fine-grained access control**: Dynamic access control that integrates with your auth system (eg: auth0, firebase-auth)
- **Scalable**: Written in golang, it follows cloud native practices and scales horizontally

Supported databases:heart::

- **MongoDB**
- **PostgreSQL** and PostgreSQL compatible databases (For eg. CockroachDB, Yugabyte, etc.)
- **MySQL** and MySQL compatible databases (For eg. TiDB, MariaDB, etc.)

## Table of Contents

- [Quick Start](#quick-start)
- [Client-side tooling](#client-side-tooling)
- [How it works](#how-it-works)
- [Support & Troubleshooting](#support--troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Quick start

### Using Docker Compose

Docker compose is the easiest way to get started with Space Cloud. Let's see how to quickly get started with MongoDB and Space Cloud using Docker Compose.

> **Note:** MongoDB is not a dependency of Space Cloud. Space Cloud can run with any of it's supported databases.

**Prerequisites:**

- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

1. Get the docker-compose file.

```bash
wget https://raw.githubusercontent.com/spaceuptech/space-cloud/master/install-manifests/quick-start/docker-compose/mongo/docker-compose.yaml
```

2. Run Space Cloud & MongoDB.

```bash
docker-compose up -d
```

3. Head over to http://localhost:4122/mission-control to open Mission Control and configure Space Cloud.

4. Create a project.

5. Head over to the Database section in Mission Control and copy paste `mongodb://mongo:27017` in the connection string.

6. Head over to the `Explorer` section and follow this guide to [insert and read data via Space Cloud using GraphQL](https://docs.spaceuptech.com/getting-started/quick-start/explore-graphql).

### Other guides

To get started with Space Cloud without Docker compose, check out the [manual](https://docs.spaceuptech.com/getting-started/quick-start/manual) guide. For production settings, checkout the [deployments](https://docs.spaceuptech.com/getting-started/deployment) guide.

## Client-side tooling
Space Cloud exposes GraphQL, HTTP, websockets and gRPC endpoints. See [setting up project](https://docs.spaceuptech.com/getting-started/setting-up-project) guide to choose a client and set it up. 

### GraphQL APIs
Space Cloud works with any GraphQL client. We recommend using [Apollo Client](https://github.com/apollographql/apollo-client). See [awesome-graphql](https://github.com/chentsulin/awesome-graphql) for a list of clients.

### Client SDKs

If you don't want to use graphql, we have made following client SDKs for you:

- [Javascript](https://docs.spaceuptech.com/getting-started/setting-up-project/javascript) for web and Nodejs projects
- [Golang](https://docs.spaceuptech.com/getting-started/setting-up-project/golang) for Golang projects


## How it works

Space Cloud is meant to replace any backend php, nodejs, java code you may write to create your endpoints. Instead, it _exposes your database over an external API_ that can be consumed directly from the frontend. In other words, it **allows clients to fire database queries directly**.

However, it's important to note that **the client does not send database (SQL) queries** to Space Cloud. Instead, it sends an object describing the query to be executed. This object is first **validated** by Space Cloud (using security rules). Once the client is authorized to make the request, **a database query is dynamically generated and executed**. The results are sent directly to the concerned client.

We understand that not every app can be built using only CRUD operations. Sometimes it's necessary to write business logic. For such cases, Space Cloud offers you APIs to write `functions` (which runs as microservices alongside Space Cloud). These `functions` can be invoked from the frontend or by other `functions`. In this scenario, Space Cloud acts merely as an api gateway between your `functions` and the client.

<div style="text-align: center">
<img src="https://spaceuptech.com/icons/space-cloud-detailed.png"  style="max-width: 80%" alt="Detailed Space Cloud architecture" />
</div>

Apart from these, Space Cloud also integrates with tons of cloud technologies to give you several other features like `realtime database` (changes in the database are synced with all concerned clients in realtime), `file storage`, etc.

## Support & Troubleshooting

The documentation and community will help you troubleshoot most issues. If you have encountered a bug or need to get in touch with us, you can contact us using one of the following channels:

- Support & feedback: [Discord](https://discord.gg/ypXEEBr)
- Issue & bug tracking: [GitHub issues](https://github.com/spaceuptech/space-cloud/issues)
- Follow product updates: [@spaceuptech](https://twitter.com/spaceuptech)

## Contributing

Space Cloud is a young project. We'd love to have you on board if you wish to contribute. To help you get started, here are a few areas you can help us with:

- Writing the documentation
- Making sample apps in React, Angular, Android, and any other frontend tech you can think of
- Deciding the road map of the project
- Creating issues for any bugs you find
- And of course, with code for bug fixes and new enhancements

## License

Space Cloud is [Apache 2.0 licensed](https://github.com/spaceuptech/space-cloud/blob/master/LICENSE).
