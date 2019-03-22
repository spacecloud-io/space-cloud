# Database Module

The database module is the core of Space Cloud. It provides instant REST, gRPC APIs on any database out there directly from the frontend. The API loosely follows the Mongo DB DSL. In addition to that, it expects a [JWT token](https://jwt.io) in the `Authorization` header. This JWT token is used along with user defined security rules to enforce authentication and authorization.

By CRUD I mean Create, Read, Update and Delete operations. These are the most basic operations that one can perform on a database. In addition to that, we offer a flexible query language (based on the Mongo DB query DSL) to slice and dice data as needed.

Currently the database module supports the following databases:
- Mongo DB
- MySQL and MySQL compatible databases
- Postgres and Postgres compatible databases

## Enable the crud (database) module

The config pertaining to crud module can be found inside the `crud` key under the `modules` object. Here's the snippet:

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

The above snippet instructs `space-cloud` to connect to MongoDB at `mongodb://localhost:27017`. All operations (create, read, update and delete) are allowed on the `todos` collection.

You can learn more about on the various parameters available for configuring the database module [here](/docs/database/config).

## API

Comming soon

## Next steps

Now you know the basics of the database module. The next step would be diving deeper into the [configuration](/docs/database/config) and its structure. Let's make sure that the apps we build are secure!

<div class="btns-wrapper">
  <a href="/docs/user-management/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/create" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
