# Space Cloud Documentation

## About Space Cloud

> **Note:** Space Cloud is still in Beta.

Space Cloud is essentially a web server that automatically integrates with an existing or a new database to provide instant realtime APIs over REST, websockets, gRPC, etc. Written in Golang, it provides a high throughput data access layer which can be consumed directly from the frontend. It's completely unopinionated and works with the tech stack of your choice.


<img src="https://spaceuptech.com/icons/space-cloud-basic.png"  alt="Basic Space Cloud architecture" />

In a nutshell, Space Cloud provides you with all of the following **without having to write a single line of backend code**:

- Ready to use functionalities like User Management, Realtime CRUD and File Storage.
- Baked-in security.
- Freedom from vendor lock-ins.
- Flexibility to work with the tech stack of your choice.

Supported databases❤️:
- **MongoDB**
- **PostgreSQL** and PostgreSQL compatible databases (For eg. CockroachDB, Yugabyte etc.)
- **MySQL** and MySQL compatible databases (For eg. TiDB)

## How it works

Space Cloud is meant to replace any backend php, nodejs, java code you may write to create your endpoints. Instead, it _exposes your database over an external API_ that can be consumed directly from the frontend. In other words, it **allows clients to fire database queries directly**.

However, it's important to note that **the client does not send database (SQL) queries** to Space Cloud. Instead, it sends an object describing the query to be executed. This object is first **validated** by Space Cloud (using security rules). Once the client is authorized to make the request, **a database query is dynamically generated and executed**. The results are sent directly to the concerned client.

<img src="https://spaceuptech.com/icons/space-cloud-detailed.png" alt="Detailed Space Cloud architecture" />

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

## What's next?

- Take Space Cloud up for a spin with our [quick started guide](/docs/quick-start).  
- [Deploy](/docs/deployment/) Space Cloud and [set up a project](/docs/setting-up-project/) from scratch.

<div class="btns-wrapper">
  <a href="/docs/quick-start" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
