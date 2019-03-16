# Database Module

The database module is the core of Space Cloud. It provides instant REST, gRPC APIs on any database out there directly from the frontend. The API loosely follows the Mongo DB DSL. In addition to that, it expects a [JWT token](https://jwt.io) in the `Authorization` header. This JWT token is used along with user defined security rules to enforce authentication and authorization.

By CRUD I mean Create, Read, Update and Delete operations. These are the most basic operations that one can perform on a database. In addition to that, we offer a flexible query language (based on the Mongo DB query DSL) to slice and dice data as needed.

Currently the database module supports the following databases:
- Mongo DB
- MySQL and MySQL compatible databases
- Postgres and Postgres compatible databases

## Prerequisites
- A running database (We'll be using MongoDB in this example)
- Space Cloud binary ([Linux](https://spaceuptech.com/downloads/linux/space-cloud.zip), [Mac](https://spaceuptech.com/downloads/darwin/space-cloud.zip), [Windows](https://spaceuptech.com/downloads/windows/space-cloud.zip))

## Configure the crud module

The config pertaining to user management module can be found inside the `crud` key under the `modules` object. Here's the snippet:

```yaml
modules:
  crud:
    mongo:
      conn: mongodb://localhost:27017
      isPrimary: true
      collections:
        todos:
          isRealtimeEnabled: false
          rules:
            create:
              rule: allow
            read:
              rule: allow
            update:
              rule: allow
            delete:
              rule: allow

  # Config for other modules go here
```

As you can see `crud`, in this case has the key `mongo` which stands for the MongoDB database. You can have multiple databases in a single project by simply adding the config of each database under `crud`. The keys for the databases we currently support are `mongo` (for MongoDB), `sql-postgres` (for Postgres) and `sql-mysql` (for MySQL).

For each database, you need to specify the following fields:
- **conn**: This is the connection string to connect to the database with.
- **isPrimary**: Specifies if the database is to be used as the primary database. Note, you **cannot have more than one primary database**.
- **collections**: These are the table / collections which need to be exposed via Space Cloud. They contain two sub fields `isRealtimeEnabled` and `rules`. `rules` are nothing but the [security rules](/docs/security), to control the database access.

The snippet shown above configures Space Cloud to use `MongoDB` as the primary database present at `mongodb://localhost:27017`. It exposes a single collection `todos`. All types of operations (create, read, update and delete) are allowed on the `todos` collection. This implies that, any anonymous user will be able to perform any operations on the database. To expose more tables / collections, simple add new objects under the `collections` key.

Here's an example that has two collections `todos` and `users`. Note, updating and deleting users is denied.

```yaml
modules:
  crud:
    mongo:
      conn: mongodb://localhost:27017
      isPrimary: true
      collections:
        todos:
          isRealtimeEnabled: false
          rules:
            create:
              rule: allow
            read:
              rule: allow
            update:
              rule: allow
            delete:
              rule: allow
        users:
          isRealtimeEnabled: false
          rules:
            create:
              rule: allow
            read:
              rule: allow
            update:
              rule: deny
            delete:
              rule: deny

  # Config for other modules go here
```

> Note: The `allow` rule must only be used cautiously. It must never be used for update and delete operations in production. You can read more about security rules [here](/docs/security). 

## Next steps

Now you know the basics of the database module. The next step would be diving deeper into the [security rules](/docs/security) and its structure. Let's make sure that the apps we build are secure!

You can also check out the [API docs](/docs/api) to start building your app right away.

<div class="btns-wrapper">
  <a href="/docs/user-management/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/realtime" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
