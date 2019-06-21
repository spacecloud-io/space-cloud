# Securing your project

Space Cloud offers robust authentication and authorization mechanism with the philosophy that applications and platforms built with Space Cloud should be secure by default. Treating security as an afterthought is a recipe for disaster.

The security module in Space Cloud secures the requests to database, functions and file storage module via security rules written by user. Any operation to a resource (database, function or a file path) not specified in the configuration file is denied. This ensures that all operations are secure by default in Space Cloud.

## What can I do with security module?
- Allow / deny a particular operation irrespective of any conditions.
- Allow a particular operation only if a client is logged in.
- Allow a particular operation only if certain conditions are met (via JSON rules or custom logic).

## How it works

As an user, you have to write security rules for the various operations on all the resources (database, file storage and custom functions) exposed by the Space Cloud. These security rules have to be written in the config file provided to Space Cloud. All incoming requests to Space Cloud are first validated by the API controller via the security module based on the corresponding rule in the config file. Validation happens in two stages: Authentication and Authorization.

### JWT based authentication
Space Cloud uses JWT-based authentication. It expects a JWT token in every incoming request. The security module verifies if the signature of the token is valid or not based on a `secret` provided to it. This makes sure that the user is authenticated and hasn't tried to change or create his / her own false token. You can check out the [official website](https://jwt.io) of the JWT project to learn more about it.


Each JWT token provided to user on signin/signup is signed with a `secret` by an authentication service (in-built or your custom service). This is the same `secret` which is provided to the config file of Space Cloud. JWT tokens also contain a JSON object (known as claims) as a payload. The user is free to decide what claims should go into the JWT token while writing a custom service. When using the in-built user authentication module, the claims consist of the `id`, `name` and `role` of the user. This JSON object is parsed and is made available as the `args.auth` variable in security rules. The integrity of the auth variable is maintained due to the nature of JWT tokens.

### Authorization
This stage decides whether an authenticated user is authorized or not to make a request. The request is validated only if the security rule is resolved based on it's type. Various ways to resolve rule:
- Matching based on the fields in the incoming request and the auth object (JWT claims).
- Making a database query.
- Calling a custom function to return `true/false` to validate the request.

## Security rule types

Security rules are a mechanism used to enforce authorization. The request is allowed to be made only if the conditions specified by the security rules are met. Currently the following rules are supported:

- **allow:** This rule is used to disable authentication and authorization entirely. The request is allowed to be made even if the JWT token is absent in the request.
- **authenticated:** This rule is used to allow the request if a valid JWT token is found in the request. No checks are imposed beyond that. Basically it authorizes every request which has passed the authentication stage.
- **deny:** This rule is to deny all incoming requests. It is especially useful to deny certain operations like `delete` while selectively allowing the other ones.
- **match:** This rule is used to allow a certain request only when a certain condition has been met. Generally it is used to match the input parameters (like the where clause or certain fields in the document to be inserted) with the auth object. It can also be used for role based authentication (match the role of user to a particular value).
- **query:** This rule is used to allow a certain request only if a database request returns successfully. The query's find clause is generated dynamically using this rule. The query is considered to be successful if even a single row is successfully returned.
- **and, or:** These rules helps you mix and match the `match` and `query` rules to tackle complex authorization tasks.
- **func:** These rule allows you to write a custom function in any language to authorize the incoming request.

## Next steps

Each module has their own way of using security rules. You can head over to the module specific security rule page to know how they are used and some examples.

- [Database](/docs/security/database)
- [File storage](/docs/security/file-storage)
- [Functions](/docs/security/functions)

<div class="btns-wrapper">
  <a href="/docs/functions/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/security/database" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
