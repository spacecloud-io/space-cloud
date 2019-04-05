# Securing database access

The security rules for database access works to authorize client request. Authorization works on the query level for each table / collection in the database. This means that you can have different rules for different queries. Here's a sample snippet which shows the rules on the `users` collection. Operations `create`  and `read` are allowed while `update` and `delete` are blocked.

```yaml
modules:
  crud:
    mongo:
      conn: mongodb://localhost:27017
      isPrimary: true
      collections:
        users:
          isRealtimeEnabled: false
          rules:                    # These are the rules applied on this collection
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

## Rules in Space Cloud

### Allow ( `allow` )
This rule is used to disable authentication and authorization entirely. The request is allowed to be made even if the JWT token is absent in the `Authorization` header.

Example (allow all create requests):
```yaml
rules:
  create:
    rule: allow
```

### Authorized ( `authorized` )
This rule is used to allow the request if a valid JWT token is found in the `Authorization`. No checks are imposed beyond that. Basically it authorizes every request which has passed the authentication stage.

Example (allow read request with a valid JWT token):
```yaml
rules:
  read:
    rule: authorized
```

### Deny ( `deny` )
This rule is to deny all incoming requests. It is especially useful to deny certain operations like `delete` while selectively allowing the other ones.

Example (deny all delete requests):
```yaml
rules:
  delete:
    rule: deny
```

### Match ( `match` )
This rule is used to allow a certain request only when a certain condition has been met. Generally it is used to match the input parameters (like the where clause or certain fields in the document to be inserted) with the auth object. It can also be used for role based authentication.

The basic syntax looks like this:

```yaml
rule: match
eval: == | != | > | >= | < | <=   # Any one of them
type: string | number | bool      # Any one of them
f1: field1                        # A value or variable
f2: field2                        # A value or variable 
```

Each CRUD request contains a set of object which are available to you to match against. All variables are stored in the `args` object.

**The object available during each request for match is:**

```json
{
  "args": {
    "auth": {},       // The JWT claims present in the `Authorization` header 
    "find": {},       // Present when a where clause is supplied. (Follows the MongoDB query syntax)
    "update": {},     // Present for update operation. (Follows the MongoDB update syntax)
    "doc": {},        // Present for insert operation. (The document(s) to be inserted)
    "op": "one | all" // Present in all requests.
  }
}
```

Example (make sure user can query only his `todos`):

```yaml
rules:
  query:
    rule: match
    eval: ==
    type: string
    f1: args.auth.id      # Assuming id is the JWT claim containing the userId
    f2: args.find.userId  # Assuming the `todos` table contains the field `userId`
```

### Database Query ( `query` )
This rule is used to allow a certain request only if a database request returns successfully. The query's find clause is generated dynamically using this rule. The query is considered to be successful if even a single row is successfully returned.

The basic syntax looks like this:
```yaml
rule: query
db:   mongo | sql-mysql | sql-postgres  # Any one of them
col:  collection                        # Name of the table / collection
find: mongo-find query                  # Find object following MongoDB query syntax
```

The `query` rule executes a database query with the user defined find object with operation type set to `one`. It is useful for policies which depend on the values stored in the database.

Example (make sure user can query only public `profiles`):

```yaml
rules:
  query:
    rule: query
    db:   mongo
    col:  profiles
    find:
      userId: args.find.userId  # Assuming profiles has field `userId`
      isPublic: true            # Assuming profiles has field `isPublic`
```

### And | Or ( `and` | `or` )
These rule helps you mix and match several `match` and `query` rules to tackle complex authorization tasks.

The basic syntax looks like this:
```yaml
rule: and | or
clauses: array_of_rules # Find object following MongoDB query syntax
```

Example (make sure the user can query a profile only if it's his or he is the admin)
```yaml
rule: or
clauses:
  - rule: match
    eval: ==
    type: string
    f1:   args.auth.role    # Assuming role is the JWT claim containing the role of the user
    f2:   admin
  - rule: match
    eval: ==
    type: string
    f1:   args.auth.id      # Assuming id is the JWT claim containing the userId
    f2:   args.find.userId  # Assuming the `profiles` table contains the field `userId`
```

## Next steps

Great! You can now start securing your app. You may now checkout the [security rules for file storage module](/docs/security/file-storage) or head over to the section to [deploy your app](/docs/deploy/overview).

<div class="btns-wrapper">
  <a href="/docs/security/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/security/file-storage" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>