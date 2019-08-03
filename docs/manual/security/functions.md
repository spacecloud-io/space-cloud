# Securing function calls

The security rules for functions works to authorize client request for function calls. Authorization works on the function level of each service. This means that you can have different rules for each function in a service. Here's a sample snippet which shows the security rules to access function `func1` of service `service1`.

```yaml
modules:
  functions:
    enabled: true
    broker: nats
    conn: nats://localhost:4222
    services:
      service1:
        functions:
          func1:
            rule: 
              rule: allow
          func2: 
            rule: 
              rule: allow

# Config for other modules go here 
```

You can add write rules for multiple functions under each service in the `functions` section. A function call is denied if there is no corresponding rule for it in the config file. This ensures that all function calls are secure by default.

## Features
With security rules for functions you can:
- Allow / deny access to a function.
- Allow a particular function call only if the user is authenticated.
- Allow a particular function call only if certain conditions are met (via JSON rules or custom logic).

## Popular use cases
- Allow only signed in users to call a function (For example only allow signed in users to make a payment).
- Role based authentication (For example only allow admin to access a particular function)
- Check if the params sent by user contains a certain field.
- Call another function to authorize the function call (For example you might have an authorization service which validates all types of request).

All these problems can be solved by the security module of Space Cloud.

## Available variables
All requests for function calls contains 2 variables which are availabe to you for matching conditions:

- **auth:** The claims present in the JWT token. If you are using in-built user management service of Space Cloud, then the `auth` has `id`, `name` and `role` of the user. While making a custom service, you are free to choose the claims which go inside the JWT token and thus available in the `auth` variable.
- **params:** The params object sent by the user to call the function.

## Allow anonymous access
 
You can disable authentication and authorization for a particular function of a service completely by using `allow`. The request is allowed to be made even if the JWT token is absent in the request. You might want to use this when you want your users to perform certain operation without signin. Here's how to give access to a particular operation using `allow`:

```yaml
services:
  service1:
    functions:
      func1:
        rule:
          rule: allow
```

## Deny access

This rule is to deny all calls to a particular function irrespective of any thing. It might be useful to temporarily deny access to a function (For example in testing). Here's how to deny access to a particular function using `deny`:

```yaml
services:
  service1:
    functions:
      func1:
        rule:
          rule: deny
```

## Allow only authenticated users

You can allow a certain function to be called by a user only if the user is authenticated. (For example, allow only logged in users to make a payment). This rule is used to allow the request only if a valid JWT token is found in the request. No checks are imposed beyond that. Basically it authorizes every request which has passed the authentication stage. Here's how to allow a function call for authenticated users:

```yaml
services:
  service1:
    functions:
      func1:
        rule:
          rule: authenticated
```

## Allow function call on certain conditions

Many a times you might want a user to call a particular function only when certain conditions are met. Such conditions might require you to check the value of certain fields from the incoming request or from the database. Or it can be a custom validation altogether. The security rules in Space Cloud are made keeping this flexibility in mind.

### Match incoming requests
This rule is used to allow a certain request only when a certain condition has been met and all the variables required for matching are present in the request itself. Every request for a function call contains 2 variables - `auth` and `params` present in the `args` object. Generally this rule is used to match the parameters sent by user with the auth object. It can also be used for role based authentication.

The basic syntax looks like this:

```yaml
rule: match
eval: == | != | > | >= | < | <=   # Any one of them
type: string | number | bool      # Any one of them
f1: field1                        # A value or variable
f2: field2                        # A value or variable 
```

Example (Match the value of a field in `params` sent by the user):

```yaml
rule: match
eval: ==
type: string
f1: args.auth.id        # Assuming id is the JWT claim containing the userId
f2: args.params.userId  # Assuming user sent userId in params
```

Example (Role based authentication - allow only admin to call a certain function):

```yaml
rule: match
eval: ==
type: string
f1: args.auth.role     # Assuming role is the JWT claim containing the role of user
f2: admin
```

Example (Check if a field is present in the `params`):

```yaml
rule: match
eval: ==
type: bool
f1: utils.exists(args.params.postId)
f2: true
```
`utils.exists` is a utility function by the security rules which checks if a given field exists or not and returns true or false.

### Database Query
This rule is used to allow a certain function call only if a database request returns successfully. The query's find clause is generated dynamically using this rule. The query is considered to be successful if even a single row is successfully returned.

The basic syntax looks like this:
```yaml
rule: query
db:   mongo | sql-mysql | sql-postgres  # Any one of them
col:  collection                        # Name of the table / collection
find: mongo-find query                  # Find object following MongoDB query syntax
```

The `query` rule executes a database query with the user defined find object with operation type set to `one`. It is useful for policies which depend on the values stored in the database.

Example (Make sure a user can call a function only if he is author of some book):

```yaml
rule: query
db:   mongo
col:  books
find:
  authorId: args.params.bookId  # Assuming books table has `authorId` field
```

### Combine multiple conditions

You can mix and match several `match` and `query` rules together to tackle complex authorization tasks (like the Instagram problem) using the `and` and `or` rule.

The basic syntax looks like this:
```yaml
rule: and | or
clauses: array_of_rules
```

Example (Make sure an user can call a function only if he has the role `admin` or `super-user`)
```yaml
rule: or
clauses:
  - rule: match
    eval: ==
    type: string
    f1: args.auth.role # Assuming role is the JWT claim containing the role of user
    f2: admin
  - rule: match
    eval: ==
    type: string
    f1: args.auth.role
    f2: super-user    
```

### Custom validations

In case where the matching and db query for validating conditions are not enough, you can authorize the request within the function itself or configure Space Cloud to call another function for authorization. Here's an example showing how to configure Space Cloud to use another function to authorize a particular function call:

```yaml
modules:
  functions:
    services:
      service1:
        functions:
          func1:
            rule: 
              rule: func
              service: auth-service
              func: auth-func
```

In the above case, to authorize a request to call `func1`, the Space Cloud will make a call to the `auth-func` function of the `auth-service`. The request to `func1` will be considered authorized by the Space Cloud only when the `auth-func` returns an object with `ack` property set to true.

## Next steps

Great! You have learned how to secure function calls. You may now head over to the [static module](/docs/static) to learn how to host static resources.

<div class="btns-wrapper">
  <a href="/docs/security/file-storage" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/static/overview" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>