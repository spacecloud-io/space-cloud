# Securing your project

Space Cloud leverages multiple security systems with the philosophy that applications and platforms should be secure by default. Treating security as an afterthought is a recipe for disaster.

Every module in Space Cloud supports JWT-based authentication. Access to any database not specified in the configuration file is denied by default. Whenever possible, use JWT-based authentication in production. You can check out the [official website](https://jwt.io) of the JWT project to learn more about it.

## How JWT is used

`space-cloud` looks for the JWT token in the `Authorization` header in each HTTP request made to it.

The first stage where JWT tokens are used is to verify if the signature of the token is valid or not. This makes sure that the user is authenticate and has tried to change or create his / her own token. This stage is what we refer to as authentication. Every module (except user management) performs authentication.

JWT tokens also contain a JSON object (known as claims) as a payload. The user is free to decide what claims must go into the JWT token. This JSON object is parsed and is made available as the `args.auth` variable in security rules. The integrity of the auth variable is maintained due to the nature of JWT tokens.

## Security rules

Security rules are a mechanism used to enforce authorization. The request is allowed to be made only if the conditions specified by the security rules are met. Currently the following rules are supported:

- **allow**: This rule is used to disable authentication and authorization entirely. The request is allowed to be made even if the JWT token is absent in the `Authorization` header.
- **authenticated**: This rule is used to allow the request if a valid JWT token is found in the `Authorization`. No checks are imposed beyond that. Basically it authorizes every request which has passed the authentication stage.
- **deny**: This rule is to deny all incoming requests. It is especially useful to deny certain operations like `delete` while selectively allowing the other ones.
- **match**: This rule is used to allow a certain request only when a certain condition has been met. Generally it is used to match the input parameters (like the where clause or certain fields in the document to be inserted) with the auth object. It can also be used for role based authentication.
- **query**: This rule is used to allow a certain request only if a database request returns successfully. The query's find clause is generated dynamically using this rule. The query is considered to be successful if even a single row is successfully returned.
- **and, or**: These rules helps you mix and match the `match` and `query` rules to tackle complex authorization tasks.

## Next steps

Each module has their own way of using security rules. You can head over to the module specific security rule page to know more.

- [Database](/docs/security/database)
- [File storage](/docs/security/file-storage)

<div class="btns-wrapper">
  <a href="/docs/functions" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/security/database" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
