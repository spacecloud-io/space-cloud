# Configuring the database module

> **Note:** All changes to the config of `space-cloud` has to be done through the Mission Control only. Changes made manually to the config file will get overwritten. 

The config pertaining to crud module can be found inside the `crud` key under the `modules` object. Here's the snippet:

```yaml
modules:
  crud:
    mongo:
      enabled: true
      conn: mongodb://localhost:27017
      collections:
        default:
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
      enabled: true
      conn: user:my-secret-pwd@/project
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
  realtime:
    enabled: true
    broker: nats
    conn: nats://localhost:4222

  # Config for other modules go here
```

The `crud` module in config specifies which databases to connect to and which collections to expose in that database along with the respective security rules. As you can see, you can have multiple databases in a single project by simply adding the config of each database under `crud`. The keys for the databases we currently support are `mongo` (for MongoDB), `sql-postgres` (for Postgres and Postgres compatible databases) and `sql-mysql` (for MySQL and MySQL compatible databases).

# Enable the realtime module

Here's a snippet showing how to **enable the realtime module**. This 4 line snippet will set up all the necessary routines required by the realtime module.

```yaml
modules:
  realtime:
    enabled: true    # Enable the realtime module globally
    broker: nats     # Broker to be used as pub sub for realtime module
    conn: nats://localhost:4222  #  Connection string of broker 
```

The realtime feature also needs to be enabled on a collection level for the collections that you want to sync in realtime. Here's a snippet configuring space cloud to use MongoDB and MySQL. MongoDB will hold the `todos` collection which will be synced in realtime while MySQL has the `users` table (not synced in realtime).

```yaml
modules:
  crud:
    mongo:
      conn: mongodb://localhost:27017
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
    sql-mysql:
      conn: user:my-secret-pwd@/project
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
  realtime:
    enabled: true
    broker: nats
    conn: nats://localhost:4222
```

For each database, you need to specify the following fields:
- **conn:** This is the connection string to connect to the database with.
- **collections:** These are the table / collections which need to be exposed via Space Cloud. They contain two sub fields `isRealtimeEnabled` and `rules`. `rules` are nothing but the [security rules](/docs/security/database), to control the database access.

The snippet shown above configures Space Cloud to use `MongoDB` present at `mongodb://localhost:27017`. It exposes a single collection `todos`. All types of operations (create, read, update and delete) are allowed on the `todos` collection. This implies that, any anonymous user will be able to perform any operations on the database. To expose more tables / collections, simply add new objects under the `collections` key.

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

Now you know the basics of the database module. The next step would be checking out the file storage module to bring files to your app!

<div class="btns-wrapper">
  <a href="/docs/database/transactions" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/file-storage/overview" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>