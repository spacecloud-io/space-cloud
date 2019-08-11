# User Management Module

> **Note:** It is recommended to use your own user management module for a production environment. The current user management module is not production ready.

User management is used for managing the various sign in methods which are generally used to develop applications. It's basically a way for the user to sign up or login into your application. In addition to that it provides the user with a `JWT token` which is used in all the other modules for authentication and authorization. 



The various sign in methods supported are:
- Basic (email & password sign in)
- OAuth (Coming Soon)

## Prerequisites
- A running database (We'll be using MongoDB in this example)
- Space Cloud binary ([Linux](https://spaceuptech.com/downloads/linux/space-cloud.zip), [Mac](https://spaceuptech.com/downloads/darwin/space-cloud.zip), [Windows](https://spaceuptech.com/downloads/windows/space-cloud.zip))

## Enable the user management module

The config pertaining to user management module can be found inside the `auth` key under the `modules` object. Here's the snippet showing how to enable **basic email authentication**. This 3 line snippet will create the necessary endpoints required by the user management module.

```yaml
modules:
  auth:
    email:
      enabled: true

  # Config for other modules go here 
```

## Configuring the database

The module assumes a `users` table / collection is present in the database. The table needs to be created in the SQL databases. MongoDB, being schemaless does not require any configuration at all.

A sample object in JSON will have the following format
```json
{
  "id":    "string", // this will be _id for MongoDB
  "email": "string",
  "name":  "string",
  "pass":  "string", // The password field is bcrypted
  "role":  "string"
}
```

The SQL query to create the `users` table would look like this:

```sql
CREATE TABLE users (
  id    VARCHAR(50) PRIMARY KEY,
  email VARCHAR(50) NOT NULL,
  name  VARCHAR(50) NOT NULL,
  pass  VARCHAR(75) NOT NULL,
  role  VARCHAR(50)
);
```

> **Note:** You can always add more fields / columns as needed.

## Next steps
That's all you need to do to enable the user management module. You can check it's usage by heading over to next page and see how to consume the endpoints on the client side.

<div class="btns-wrapper">
  <a href="/docs/quick-start/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/user-management/signin" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>