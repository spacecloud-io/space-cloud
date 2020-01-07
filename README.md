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

Space Cloud is a cloud-native backend server that provides **instant, realtime APIs on any database**, with [event triggers](https://docs.spaceuptech.com/advanced/event-triggers) and unified APIs for your [custom business logic](https://docs.spaceuptech.com/essentials/remote-services).

Space Cloud helps you build modern applications without having to write any backend code in most cases.

It provides **GraphQL** and **REST** APIs which can be consumed directly by your frontend in a [secure manner](https://docs.spaceuptech.com/auth).


## Features 

View complete feature set [here](https://docs.spaceuptech.com/getting-started/introduction/features).

- **Powerful CRUD**: Flexible queries, transactions and cross-database joins
- **Realtime**: Make live queries to your database
- **File storage**: Upload/download files to scalable file stores (e.g., Amazon S3, Google Cloud Storage)
- **Extensible**: Unified APIs for your custom HTTP services
- **Event-driven**: Trigger webhooks or serverless functions on database or file storage events
- **Fine-grained access control**: Dynamic access control that integrates with your auth system (e.g., auth0, firebase-auth)
- **Scalable**: Written in Golang, it follows cloud-native practices and scales horizontally

Supported databases:heart::

- **MongoDB**
- **PostgreSQL** and PostgreSQL compatible databases (For eg. CockroachDB, Yugabyte, etc.)
- **MySQL** and MySQL compatible databases (For eg. TiDB, MariaDB, etc.)
- **SQL Server**

## Table of Contents

- [Quick Start](#quick-start)
- [Client-side tooling](#client-side-tooling)
- [How it works](#how-it-works)
- [Support & Troubleshooting](#support--troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Quick start

### Using Docker Compose

Docker Compose is the easiest way to get started with Space Cloud. Let's see how to quickly get started with Postgres and Space Cloud using Docker Compose.

> **Note:** Postgres is not a dependency of Space Cloud. Space Cloud can run with any of it's supported databases.

**Prerequisites:**

- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

1. Get the docker-compose file.

```bash
wget https://raw.githubusercontent.com/spaceuptech/space-cloud/master/install-manifests/quick-start/docker-compose/postgres/docker-compose.yaml
```

2. Run Space Cloud & Postgres.

```bash
docker-compose up -d
```

3. Head over to http://localhost:4122/mission-control to open Mission Control and configure Space Cloud.

4. Create a project.

5. Then add Postgres to your project with the following connection string: `postgres://postgres:mysecretpassword@postgres:5432/postgres?sslmode=disable` and hit `Save`.

6. Head over to the `Explorer` section and follow this guide to [insert and read data via Space Cloud using GraphQL](https://docs.spaceuptech.com/getting-started/quick-start/explore-graphql).

### Other guides

To get started with Space Cloud without Docker compose, check out the [manual](https://docs.spaceuptech.com/getting-started/quick-start/manual) guide. For production settings, check out the [deployments](https://docs.spaceuptech.com/getting-started/deployment) guide.

## Client-side tooling
Space Cloud exposes GraphQL and REST APIs. See [setting up project](https://docs.spaceuptech.com/getting-started/setting-up-project) guide to choose a client and set it up. 

### GraphQL APIs
GraphQL is the recommended way to use Space cloud, and it works with any GraphQL client. However, we recommend using [Apollo Client](https://github.com/apollographql/apollo-client). See [awesome-graphql](https://github.com/chentsulin/awesome-graphql) for a list of clients.

### REST APIs

You can use the [REST APIs of Space Cloud](https://app.swaggerhub.com/apis/YourTechBud/space-cloud/0.15.0) if you are more comfortable with REST. 

To make it easy to consume the REST APIs in web projects, we have created a [**Javascript SDK**](https://docs.spaceuptech.com/getting-started/setting-up-project/javascript) for you.

## How it works

Space Cloud is meant to replace any backend php, nodejs, java code you may write to create your endpoints. Instead, it _exposes your database over an external API_ that can be consumed directly from the frontend. In other words, it **allows clients to fire database queries directly**.

However, it's important to note that **the client does not send database (SQL) queries** to Space Cloud. Instead, it sends an object describing the query to be executed. This object is first **validated** by Space Cloud (using security rules). Once the client is authorized to make the request, **a database query is dynamically generated and executed**. The results are sent directly to the concerned client.

We understand that not every app can be built using only CRUD operations. Sometimes it's necessary to write business logic. For such cases, Space Cloud allows you to access your **custom HTTP servers** via the same consistent APIs of Space Cloud.  In this scenario, Space Cloud acts merely as an API gateway between your `remote-services` and the client. However, the cool part is that you can even perform **joins on your remote services and database** via the GraphQL API of Space Cloud

<div style="text-align: center">
<img src="https://spaceuptech.com/icons/space-cloud-detailed.png"  style="max-width: 80%" alt="Detailed Space Cloud architecture" />
</div>

Apart from these, Space Cloud also integrates with tons of cloud technologies to give you several other features like `realtime database` (changes in the database are synced with all concerned clients in realtime), `file storage`, etc.

## Support & Troubleshooting

The documentation and community should help you troubleshoot most issues. If you have encountered a bug or need to get in touch with us, you can contact us using one of the following channels:

- Support & feedback: [Discord](https://discord.gg/ypXEEBr)
- Issue & bug tracking: [GitHub issues](https://github.com/spaceuptech/space-cloud/issues)
- Follow product updates: [@spaceuptech](https://twitter.com/spaceuptech)

## Contributing

Space Cloud is a young project. We'd love to have you onboard if you wish to contribute. To help you get started, here are a few areas you can help us with:

- Writing the documentation
- Making sample apps in React, Angular, Android, and any other frontend tech you can think of
- Deciding the road map of the project
- Creating issues for any bugs you find
- And of course, with code for bug fixes and new enhancements

## License

Space Cloud is [Apache 2.0 licensed](https://github.com/spaceuptech/space-cloud/blob/master/LICENSE).
