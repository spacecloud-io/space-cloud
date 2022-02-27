<p align="center"><a href="https://space-cloud.io"><img src="https://space-cloud.io/images/kit/logo_full.svg" alt="Space Cloud"></a></p>

<h3 align="center">
  Develop, Deploy and Secure Serverless Apps on Kubernetes.
</h3>

<p align="center">
  <strong>
    <a href="https://space-cloud.io/">Website</a>
    •
    <a href="https://docs.space-cloud.io/">Docs</a>
    •
    <a href="https://discord.gg/RkGjW93">Support</a>
  </strong>
</p>
<p align="center">
    <a href="https://discord.gg/RkGjW93"><img src="https://img.shields.io/badge/chat-discord-brightgreen.svg?logo=discord&%20style=flat"></a>
    <br/>
    <a href="https://goreportcard.com/report/github.com/spaceuptech/space-cloud">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/spaceuptech/space-cloud">
    </a>
    <a href="https://opensource.org/licenses/Apache-2.0"><img
    alt="Apache License"
    src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"></a>
    <a href="https://twitter.com/intent/follow?screen_name=spacecloudio"><img src="https://img.shields.io/badge/Follow-spacecloudio-blue.svg?style=flat&logo=twitter"></a>
</p>

Space Cloud is a Kubernetes based serverless platform that provides **instant, realtime APIs on any database**, with [event triggers](https://docs.space-cloud.io/microservices/eventing) and unified APIs for your [custom business logic](https://docs.space-cloud.io/microservices/graphql).

Space Cloud helps you build modern applications without having to write any backend code in most cases.

It provides **GraphQL** and **REST** APIs which can be consumed directly by your frontend in a [secure manner](https://docs.spaceuptech.com/storage/database/securing-apis).

## Features 

View complete feature set [here](https://docs.spaceuptech.com/introduction/features).

- **Powerful CRUD**: Flexible queries, transactions, aggregations and cross-database joins
- **Realtime**: Make live queries to your database
- **File storage**: Upload/download files to scalable file stores (e.g., Amazon S3, Google Cloud Storage)
- **Extensible**: Unified APIs for your custom HTTP services
- **Event-driven**: Trigger webhooks or serverless functions on database or file storage events
- **Fine-grained access control**: Dynamic access control that integrates with your auth system (e.g., auth0, firebase-auth)
- **Scalable**: Written in Golang, it follows cloud-native practices and scales horizontally
- **Service Mesh**: Get all the capabilities of a service mesh without having to learn any of that!
- **Scale down to zero**: Auto scale your http workloads including scaling down to zero

Supported databases :heart::

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

> **If you are new to Space Cloud, we strongly recommend following our [step-by-step guide](https://learn.spaceuptech.com/space-cloud/basics/setup/) to get started**

### Other guides

View the installation guides for [Docker](https://docs.spaceuptech.com/install/docker) and [Kubernetes](https://docs.spaceuptech.com/install/kubernetes).

## Client-side tooling
Space Cloud exposes GraphQL and REST APIs. See [setting up project](https://docs.spaceuptech.com/introduction/setting-up-project) guide to choose a client and set it up. 

### GraphQL APIs
GraphQL is the recommended way to use Space cloud, and it works with any GraphQL client. However, we recommend using [Apollo Client](https://github.com/apollographql/apollo-client). See [awesome-graphql](https://github.com/chentsulin/awesome-graphql) for a list of clients.

### REST APIs

You can use the [REST APIs of Space Cloud](https://app.swaggerhub.com/apis/YourTechBud/space-cloud/0.15.0) if you are more comfortable with REST. 

To make it easy to consume the REST APIs in web projects, we have created a [**Javascript SDK**](https://docs.spaceuptech.com/introduction/setting-up-project/javascript) for you.

## How it works

Space Cloud is meant to replace any backend php, nodejs, java code you may write to create your endpoints. Instead, it _exposes your database over an external API_ that can be consumed directly from the frontend. In other words, it **allows clients to fire database queries directly**.

However, it's important to note that **the client does not send database (SQL) queries** to Space Cloud. Instead, it sends an object describing the query to be executed. This object is first **validated** by Space Cloud (using security rules). Once the client is authorized to make the request, **a database query is dynamically generated and executed**. The results are sent directly to the concerned client.

We understand that not every app can be built using only CRUD operations. Sometimes it's necessary to write business logic. For such cases, Space Cloud allows you to access your **custom HTTP servers** via the same consistent APIs of Space Cloud.  In this scenario, Space Cloud acts merely as an API gateway between your `services` and the client. However, the cool part is that you can even perform **joins on your microservices and database** via the GraphQL API of Space Cloud.

<div style="text-align: center">
<img src="https://space-cloud.io/images/graphQL-diagram.svg"  style="max-width: 80%" alt="Detailed Space Cloud architecture" />
</div>

Space Cloud integrates with [Kubernetes](https://kubernetes.io) and [Istio](https://istio.io) natively to bring to you a highly scalable Serverless Platform. It encrypts all traffic by default and lets you describe communication policies to protect your microservices.

With that, it also provides **autoscaling functionality** out of the box including **scaling down to zero**.

## Support & Troubleshooting

The documentation and community should help you troubleshoot most issues. If you have encountered a bug or need to get in touch with us, you can contact us using one of the following channels:

- Support & feedback: [Discord](https://discord.gg/RkGjW93)
- Issue & bug tracking: [GitHub issues](https://github.com/spacecloud-io/space-cloud/issues)
- Follow product updates: [@spaceupcloudio](https://twitter.com/spacecloudio)

## Contributing

Space Cloud is a young project. We'd love to have you onboard if you wish to contribute. To help you get started, here are a few areas you can help us with:

- Writing the documentation
- Making sample apps in React, Angular, Android, and any other frontend tech you can think of
- Deciding the road map of the project
- Creating issues for any bugs you find
- And of course, with code for bug fixes and new enhancements

## License

Space Cloud is [Apache 2.0 licensed](https://github.com/spacecloud-io/space-cloud/blob/master/LICENSE).
