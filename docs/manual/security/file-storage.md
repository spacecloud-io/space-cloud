# Securing file storage

The security rules for file storage access works to authorize client request. Authorization works on the operation level (read, create and delete) for each path prefix. This means that you can have different rules for different path prefixes. Here's a sample snippet which shows the rules on the `/images/:userId` prefix. Operations `create`  and `read` are allowed while `delete` is blocked.

```yaml
modules:
  fileStore:
    enabled: true
    storeType: local
    conn: /tmp/data
    rules:                        # `rules` is a map of mutiple rules
      - prefix: /images/:userId   # `prefix` is the path prefix on which the rule applies
        rule:
          create:                 # `create` is the rule object for file write operations
            rule: allow
          read:                   # `read` is the rule object for file real operations         
            rule: allow
          delete:                 # `delete` is the rule object for file delete operations
            rule: deny

  # Config for other modules go here
```

## Points to note

The following needs to be kept in mind for the security rules in the file storage module.
- The `rules` is a map of rules. The key (`imageRule` in this case) is just a unique key to indentify each rule.
- Using the `allow` rule will instruct Space Cloud to disable JWT token parsing for that function. This means the `auth` parameter in the function will always be a `null` value when the rule is set to allow.
- All rules are applied on a `prefix`. A `prefix` is nothing but the path prefix where the file / folder is present or is to be created 
- The prefix may contain path parameters (`/:userId` in this case). The value of the path parameter is available in the `args.params` object. The key would be `userId` and the value would be the actual value in the path provided.

## Features
With security rules for functions you can:
- Allow / deny access to a file/folder.
- Allow access to a particular file/folder only if the user is authenticated.
- Allow access to a particular file/folder only if certain conditions are met (via JSON rules or custom logic).

## Popular use cases
- Allow only signed in users to upload a file.
- Role based authentication (For example only allow admin to delete a particular folder)
- Allow user to upload a file at his path only.
- The Instagram problem - Allow user to view a profile pic only if the profile is public or they are following them 
- Call another function to authorize the file/folder access

All these problems can be solved by the security module of Space Cloud.

## Available variables
All requests for file/folder access contains 2 variables which are availabe to you for matching conditions:

- **auth:** The claims present in the JWT token. If you are using in-built user management service of Space Cloud, then the `auth` has `id`, `name` and `role` of the user. While making a custom service, you are free to choose the claims which go inside the JWT token and thus available in the `auth` variable.
- **params:** The variables in the path prefix.


## Allow anonymous access
 
You can disable authentication and authorization for a particular file/folder completely by using `allow`. The request is allowed to be made even if the JWT token is absent in the request. You might want to use this when you want some files to be publicly available to all users even without signin (For example, access images of products in an e-commerce app). Here's how to give access to a particular path using `allow`:

```yaml
rule:
  read:
    rule: allow
```

## Deny access

This rule is to deny all incoming requests irrespective of any thing. It is especially useful to deny certain dangerous operations like `delete` while selectively allowing the other ones. (For example, deny access to delete product's image). Here's how to deny access to a particular operation using `deny`:

```yaml
rules:
  delete:
    rule: deny
```

## Allow only authenticated users

You can allow access to a certain path only if the user is authenticated. This rule is used to allow the request only if a valid JWT token is found in the request. No checks are imposed beyond that. Basically it authorizes every request which has passed the authentication stage. Here's how to allow a operation for authenticated users:

```yaml
rules:
  create:
    rule: authenticated
```

## Allow operation on certain conditions

Many a times you might want a user to access a file path only when certain conditions are met. For example, a user can delete a picture only if it is uploaded by him. Another use case might be allowing a user to read a profile's image only if that profile is public or the user is following him (Instagram problem). Such conditions might require you to check the value of certain fields from the incoming request or from the database. Or it can be a custom validation altogether. The security rules in Space Cloud are made keeping this flexibility in mind.

### Match incoming requests
This rule is used to allow a certain request only when a certain condition has been met and all the variables required for matching is present in the request itself. Each file access request contains 2 variables (`auth` and `params`) present in the `args` object. Generally this rule is used to match the path parameters with the auth object. It can also be used for role based authentication.

The basic syntax looks like this:

```yaml
rule: match
eval: == | != | > | >= | < | <=   # Any one of them
type: string | number | bool      # Any one of them
f1: field1                        # A value or variable
f2: field2                        # A value or variable 
```

Example (Make sure user can upload a file at his path only):

```yaml
rules:
  profilePicRule: 
    prefix: /users/:userId   # `prefix` is the path prefix on which the rule applies
      rule:
        create:                 # `create` is the rule object for file write operations
          rule: match
          eval: ==
          type: string
          f1: args.params.userId # params contain the userId variable from the path prefix
          f2: args.auth.id       # Assuming role is the JWT claim containing the role of user
```

Example (Role based authentication - allow only admin to delete a folder):

```yaml
rules:
  profilePicRule: 
    prefix: /projects/:projectId   # `prefix` is the path prefix on which the rule applies
      rule:
        delete:                 # `create` is the rule object for file write operations
          rule: match
          eval: ==
          type: string
          f1: args.auth.role     # Assuming role is the JWT claim containing the role of user
          f2: admin      
```

### Database Query
This rule is used to allow a certain request only if a database request returns successfully. The query's find clause is generated dynamically using this rule. The query is considered to be successful if even a single row is successfully returned.

The basic syntax looks like this:
```yaml
rule: query
db:   mongo | sql-mysql | sql-postgres  # Any one of them
col:  collection                        # Name of the table / collection
find: mongo-find query                  # Find object following MongoDB query syntax
```

The `query` rule executes a database query with the user defined find object with operation type set to `one`. It is useful for policies which depend on the values stored in the database.

Example (make sure user can query images of only public `profiles`):

```yaml
rules:
  profilePicRule: 
    prefix: /profiles/:profileId
      read:
        rule:   query
        db:   mongo
        col:  profiles
        find:
          userId: args.params.profileId  # Assuming profiles has field `userId`
          isPublic: true                 # Assuming profiles has field `isPublic`  
```

### Combine multiple conditions

You can mix and match several `match` and `query` rules together to tackle complex authorization tasks (like the Instagram problem) using the `and` and `or` rule.

The basic syntax looks like this:
```yaml
rule: and | or
clauses: array_of_rules
```

Example (The Instagram problem - Make sure the user can view a profile picture if the profile is public or he is a follower)
```yaml
rules:
  profilePicRule: 
    prefix: /profiles/:profileId
      read:
        rule: or
        clauses:
          - rule: query
            db:   mongo
            col:  profiles
            find:
              userId: args.params.profileId  # Assuming profiles has field `userId`
              isPublic: true 
          - rule: query
            db:   mongo
            col:  profiles
            find:
              followers:
                $in: args.auth.userId   # Assuming followers is an array of user ids
```

### Custom validations

In case where the matching and db query for validating conditions are not enough, you can bring your own custom validation logic by writing a function with the functions module. You can configure Space Cloud to use your function to authorize a particular request. Here's an example showing how to do this by rule `func`:

```yaml
## Asuming you have written a service `my-service` with a function `my-service` using functions mesh
rules:
  profilePicRule: 
    prefix: /profiles/:profileId
      read:
        rule: func
        service: my-service
        func: my-func
```

In the above case, `my-func` will receive the path `auth` and `params` objects as the arguments of the function. See [functions mesh](/docs/functions) to understand how to write a function. The request will be considered authorized by the Space Cloud in this case only when the function returns an object with `ack` property set to true.

## Next steps
Great! You have learned how to secure file access. You may now checkout the [security rules for functions module](/docs/security/functions).

<div class="btns-wrapper">
  <a href="/docs/security/database" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/security/functions" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>