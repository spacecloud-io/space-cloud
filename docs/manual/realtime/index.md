# Realtime database

The realtime module is responsible to propagate database changes to the concerned clients in realtime. The clients execute a `liveQuery` which is similar to a normal query with a where clause. The realtime module makes sure that whenever there is any change in the result set (results get added, removed or updated), the client will be updated.

## How it works?

The realtime module guarantees that all changes will propagate to the client irrespective of failures on the backend. To achieve these guarantees, Space Cloud uses [Apache Kafka](https://apache.kafka.org) under the hood.

However, Space Cloud doesn't use Kafka by default. It uses Kafka ony when the `PROD` environment variable is set to true or when `space-cloud` is started with the `--prod` flag.

> **Note:** The distributed mode of realtime module is not implemented yet.

## Limitations

There are a few limitations in the queries which can be used while using the realtime module.

- The primary key must be `_id` in case of MongoDB and `id` for MySQL and Postgres and of type `text`, `varchar` or `string`.
- All documents to be inserted must have the `_id` (for MongoDB) or `id` (for MySQL and Postgres) fields set.
- The fields used in the where clause of `liveQuery` should not be updated by another request.
- All updates and deletes can be made on a single document only using the `_id` or `id` field in the where clause.

## Enable the realtime module

The configuration pertaining to the realtime module can be found under the `realtime` key under the `modules` object. Here's a snippet showing how to **enable the realtime module**. This 3 line snippet will set up all the necessary routines required by the realtime module.

```yaml
modules:
  realtime:
    enabled: true
    kafka: localhost

# Config for other modules go here 
```

You also need to to enable the feature on the collections you want to sync. Here's a snippet which indicates that the collection `todos` is realtime.

```yaml
modules:
  crud:
    mongo:
      conn: mongodb://localhost:27017
      isPrimary: true
      collections:
        todos:
          isRealtimeEnabled: true   # This makes the todos collection realtime
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

## Next steps
That's all you need to do to enable the realtime module. It's time to go ahead and how to use the realtime module on the frontend. Here are the client side APIs.
- [Javascript Client](/docs/api/javascript)
- Java (Coming soon!)

<div class="btns-wrapper">
  <a href="/docs/database" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/file-storage" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>