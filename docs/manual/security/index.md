# Securing your project

There are several security provisions taken in Space Cloud. Our philosophy is to build applications and platforms which are secure by default. In today's world where security is given maximum priority, thinking of security as an after thought isn't possible.

Every module in Space Cloud uses JWT tokens for authentication. It is switched off by default. You can check out this [website](https://jwt.io) to learn more about JWT tokens. `space-cloud` looks for the JWT token in the `Authorization` header in each HTTP request made to it.

## How are JWT tokens used

The first stage where JWT tokens are used is to verify if the signature of the token is valid or not. This makes sure that the user is authenticate and has tried to change or create his / her own token. This stage is what we refer to as authentication. Every module (except user management) performs authentication.

JWT tokens also contain a json object known as claims as a payload. The user is free to decide what claims must go into the JWT token. This json object is parsed and is made available as the `args.auth` variable in security rules. The integrity of the auth variable is maintained due to the nature of JWT tokens.

## Security rules

Security rules is a mechanism used to enforce authorization. The request is allowed to be made only if the conditions specified by the security rules are met. Currently the following rules are supported:
- **allow:** This rule is used to disable authentication and authorization entirely. The request is allowed to be made even if the JWT token is absent in the `Authorization` header.
- **authenticated**: This rule is used to allow the request if a valid JWT token is found. No checks are imposed beyond that. Basically it authorizes every request which has passed the authentication stage.
- **deny**: This rule is to deny all incoming requests. It is especially useful to deny certain operations like `delete` which selectively allowing the other ones.
- **match**: This rule is used to allow a certain request only when a certain condition has been met. Generally is is used to match the input parameters (like the where clause or certain fields in the document to be inserted) with the auth object. It can also be used for role based authentication.
- **query**: This rules is used to allow a certain request only if a database request returns successfully. The query's where clause is generated dynamically using this rule. The query is considered to be successful if even a single row is successfully returned.
- **and, or**: These rules helps you mix and match the `match` and `query` rules to tackle complex authorization tasks.

## Next steps

Each module has their own way of using security rules. You can head over to the module specific security rule page to know more.
- [Database](/docs/security/database)
- [File storage](/docs/security/file-storage)