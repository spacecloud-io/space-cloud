# Database Module

The database module is the core of Space Cloud. It provides instant realtime APIs on any database out there without having to write backend code. These APIs can be consumed in a secure manner directly from the frontend or from any other microservices your write at the backend. The APIs provided by this module are consistent across all the databases, simplifying database access. Thus, it acts as an unified data access layer making it easier to develop apps.

## Features

- Basic CRUD operations (create, read, update and delete)
- Get notified of the changes to your data in realtime
- Slice and dice the data being read as per your needs (i.e. filtering, sorting, skipping, limiting)
- Batch multiple mutations in a single db operation
- Transactions (Coming soon)
- Joins (Coming soon)

## Supported databases

- Mongo DB
- MySQL and MySQL compatible databases (TiDB)
- Postgres and Postgres compatible databases (CockroachDB, YugabyteDB)

## How it works

The client sends a JSON request (loosely following the MongoDB DSL) to Space Cloud describing the CRUD operation(s) to be performed. The API controller in Space Cloud first validates the request via the [security module](/docs/security). If the request is validated, then the CRUD controller converts the JSON query to the native DB query. This native db query is then provided to the respective db driver for execution on behalf of the client. The response from the db is sent back to the client.

If any clients had subscribed to realtime changes in database, then the realtime module propagates the database changes to the concerned clients in realtime. The realtime module makes sure that whenever there is any change in the result set (results get added, removed or updated), the client will be updated. The realtime module uses a pub-sub broker under the hood to make sure the realtime piece works in a distributed fashion.

> **Note:** Space Cloud runs a Nats server by default in the same process so that you don't have to run a broker. However, you can run your own broker and configure Space Cloud to use that instead. As of now, only Nats is supported as a broker with RabbitMQ and Kafka coming in future. 

## Limitations

There are a few limitations in the queries which can be used while using the `realtime` module.

- The primary key must be _id in case of MongoDB and id for MySQL and Postgres and of type text, varchar or string.
- All documents to be inserted must have the _id (for MongoDB) or id (for MySQL and Postgres) fields set.
- The fields used in the where clause of liveQuery should not be updated by another request.
- The insert many operation is not allowed. You can only insert one document at a time.
- All updates and deletes can be made on a single document only using the _id or id field in the where clause.

> **Note:** This limitations are only applicable if you intend to use the realtime functionality.

## Configure the database module

Head over to the `Database` section in Mission Control to configure the database module.

> **Note:** Make sure you have selected the right database from the topbar.

### Enabling the database
Enable the switch in the upper right corner of `Database` section to enable the selected database.

### Connection string
The `Connection string` input takes the connection string of your database.

You can use `environment variables` in the connection string to take db credentials. For example: `$MONGO_URL`

### Configuring collections

Mission Control by default creates the configuration for the `default` collection/table for you. The configuration of `default` collection/table is used when the configuration of requested collection/table is not found.    

The configuration of table/collection looks like the following: 

```json
{
  "isRealtimeEnabled": true,
  "rules": {
    "create": {
      "rule": "allow"
    },
    "read": {
      "rule": "allow"
    },
    "update": {
      "rule": "allow"
    },
    "delete": {
      "rule": "allow"
    }
  }
}
```

All operations (create, read, update and delete) are allowed in the above example. You can learn more about securing the functions module in depth over [here](/docs/security/database).  

## Next steps

Now you know the basics of the database module. The next step would be diving deeper into the [configuration](/docs/database/config) and its structure. You could also continue to see how to use the database module on the frontend.

<div class="btns-wrapper">
  <a href="/docs/user-management/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/create" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
