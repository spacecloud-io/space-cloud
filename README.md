<p align="center"><a href="https://www.spaceuptech.com"><img src="https://spaceuptech.com/icons/space-cloud-block-diagram1.png" alt="Space Cloud"></a></p>

<h3 align="center">
	Instant Realtime APIs on any database
</h3>

<p align="center">
	<strong>
		<a href="https://spaceuptech.com/">Website</a>
		•
		<a href="https://spaceuptech.com/docs/">Docs</a>
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

Space Cloud replaces your traditional backend servers and simplifies app development:

- **_Instant_**: Various pre-built modules such as User Management, Realtime CRUD and File Storage
- **_Secure_**: Authentication and authorization baked in by default
- **_Extensible_**: Provision to write custom backend logic


Supported databases:heart::
- **MongoDB**
- **PostgreSQL** and PostgreSQL compatible databases (For eg. CockroachDB, Yugabyte etc.)
- **MySQL** and MySQL compatible databases (For eg. TiDB)

## Table of Contents

- [Motivation](#motivation)
- [About Space Cloud](#about-space-cloud)
- [How it works](#how-it-works)
- [Design Goals](#design-goals)
  - [Ease of use](#ease-of-use)
  - [Security](#security)
  - [Enterprise-ready](#enterprise-ready)
  - [Leverage the existing tools](#leverage-the-existing-tools)
- [Documentation](#documentation)
- [Quick start](#quick-start)
- [Support & Troubleshooting](#support--troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Motivation

Making enterprise scale apps at the speed of **prototyping** is still a distant dream for many of us. Even a simple chat app becomes complicated **at scale**. Following best practices when starting something from scratch is very time consuming. Securing your app is a different ball game altogether.

Well, there are some excellent tools out there which help simplify app development like [Google Firebase](https://firebase.google.com/). But these tools come with their own **vendor lock-ins**. They either force you to use their own cloud or only work with a particular database. The next-gen apps on the other hand always require multiple task-specific databases and should even run on a private cloud.

We believed that **technology should adapt to your needs** and not the other way around. You should be able to choose any **database**, **cloud vendor** or **technology** of your preference. Agility should not come at the cost of flexibility. Space Cloud was born to solve precisely these problems.

## About Space Cloud

> **Note:** Space Cloud is still in Beta.

Space Cloud is essentially a web server that automatically integrates with an existing or a new database to provide instant realtime APIs over REST, websockets, gRPC, etc. Written in Golang, it provides a high throughput data access layer which can be consumed directly from the frontend. It's completely unopinionated and works with the tech stack of your choice.

<div style="text-align: center">
<img src="https://spaceuptech.com/icons/space-cloud-basic.png"  style="max-width: 80%" alt="Basic Space Cloud architecture" />
</div>

In a nutshell, Space Cloud provides you with all of the following **without having to write a single line of backend code**:

- Ready to use functionalities like User Management, Realtime CRUD and File Storage.
- Baked-in security.
- Freedom from vendor lock-ins.
- Flexibility to work with the tech stack of your choice.

## How it works

Space Cloud is meant to replace any backend php, nodejs, java code you may write to create your endpoints. Instead, it _exposes your database over an external API_ that can be consumed directly from the frontend. In other words, it **allows clients to fire database queries directly**.

However, it's important to note that **the client does not send database (SQL) queries** to Space Cloud. Instead, it sends an object describing the query to be executed. This object is first **validated** by Space Cloud (using security rules). Once the client is authorized to make the request, **a database query is dynamically generated and executed**. The results are sent directly to the concerned client.

We understand that not every app can be built using only CRUD operations. Sometimes it's necessary to write business logic. For such cases, Space Cloud offers you APIs to write `functions` (which runs as microservices alongside Space Cloud). These `functions` can be invoked from the frontend or by other `functions`. In this scenario, Space Cloud acts merely as an api gateway between your `functions` and the client.

<div style="text-align: center">
<img src="https://spaceuptech.com/icons/space-cloud-detailed.png"  style="max-width: 80%" alt="Detailed Space Cloud architecture" />
</div>

Apart from these, Space Cloud also integrates with tons of cloud technologies to give you several other features like `realtime database` (changes in the database are synced with all concerned clients in realtime), `file storage`, etc.

## Design Goals

There are a lot of design decisions taken by us while creating Space Cloud. These form the guiding principles which heavily influence the roadmap ahead. Understanding them would also make our **objectives of creating Space Cloud** a lot more clear.

### Ease of use

The main reason Space Cloud was born was to simplify the app/web development process. Right from making simple CRUD operations to syncing data reliably in a distributed environment, **everything must be as simple as a function call**. This is the prime reason we chose to have a consistent API across all the databases/technologies we support.

This also means that Space Cloud needs to be as unopinionated as possible to reuse the existing skill sets and tech you might be well versed with.

### Security

We take security a bit too seriously. In fact, we are close to being paranoid about it. All products built with Space Cloud must be highly secure.

The idea of exposing your database over a public API doesn't sound like a good one. But to make sure we can do it in a secure manner, we have added a powerful yet flexible feature called `security rules`. These `security rules` (written in JSON or YAML) along with JWT tokens help you take care of a wide variety of authentication and authorization problems.

### Enterprise-ready

We believe that each app built with Space Cloud must be extremely robust and future proof. We shall never comprise on the robustness of the platform at any cost. This also implies that we need to maintain strict backward compatibility.

### Leverage the existing tools

The goal of this project is not to re-invent the wheel over and over again. In fact, integration with proven technologies is preferred over implementing them ourselves. For example, we are using [Nats](https://nats.io/) under the hood, to implement `realtime`  and `functions` modules for high throughput and scale.

## Documentation

We are working hard to document every aspect of Space Cloud to give you the best onboarding experience. Here are links to the various docs we have:

- [Space Cloud](https://spaceuptech.com/docs/)
- Client APIs:
  - [JavaScript client](https://github.com/spaceuptech/space-api-js)
  - [Python client](https://github.com/spaceuptech/space-api-python)
  - [Go client](https://github.com/spaceuptech/space-api-go)
  - [Java client](https://github.com/spaceuptech/space-api-java) (In Progress)

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

6. Head over to the `Explorer` section and follow this [guide](https://spaceuptech.com/docs/quick-start/docker-compose) to insert and read data with Space APIs.

### Other guides
To get started with Space Cloud without Docker compose, check out the [manual](https://spaceuptech.com/docs/quick-start/manual) guide. For production settings, checkout the [deployments](https://spaceuptech.com/docs/deployment/overview) guide.

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
