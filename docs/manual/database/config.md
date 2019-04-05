# Configuring the database module

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

As you can see `crud`, in this case has the key `mongo` which stands for the MongoDB database. You can have multiple databases in a single project by simply adding the config of each database under `crud`. The keys for the databases we currently support are `mongo` (for MongoDB), `sql-postgres` (for Postgres) and `sql-mysql` (for MySQL).

Here's a snippet configuring space cloud to use MongoDB and MySQL. MongoDB will hold the `todos` collection while MySQL has the `users` table.

```yaml
modules:
  crud:
    mongo:
      conn: mongodb://localhost:27017
      isPrimary: false
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
    sql-mysql:
      conn: root:my-secret-pw@/project
      isPrimary: true
      collections:
        users:
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

For each database, you need to specify the following fields:
- **conn:** This is the connection string to connect to the database with.
- **isPrimary:** Specifies if the database is to be used as the primary database. Note, you **cannot have more than one primary database**.
- **collections:** These are the table / collections which need to be exposed via Space Cloud. They contain two sub fields `isRealtimeEnabled` and `rules`. `rules` are nothing but the [security rules](/docs/security/database), to control the database access.

The snippet shown above configures Space Cloud to use `MongoDB` as the primary database present at `mongodb://localhost:27017`. It exposes a single collection `todos`. All types of operations (create, read, update and delete) are allowed on the `todos` collection. This implies that, any anonymous user will be able to perform any operations on the database. To expose more tables / collections, simply add new objects under the `collections` key.

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

> **Note:** The `allow` rule must only be used cautiously. It must never be used for update and delete operations in production. You can read more about security rules [here](/docs/security/database). 

## Next steps

Now you know the basics of the database module. The next step would be checking out the realtime module to bring realtime updates to your app!

<div class="btns-wrapper">
  <a href="/docs/database/delete" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/realtime/overview" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>