# User Management Module

> Note: It is recommended to use your own user management module for a production environment. The current user management module is not production ready.

User management is used for managing the various sign in methods which are generally used to develop applications. It's basically a means for user to Sign Up or Log In into your application. In addition to that it provides the user with a JWT token which is used in all th other modules for authentication and authorization. 

The various sign in methods supported are:
- OAuth (Coming Soon)
  - Facebook
  - Google
  - Github
  - Twitter
- Basic (email & password sign in)

## Prerequisites
- [Space CLI](https://spaceuptech.com/docs/getting-started/space-cli)
- [Docker](https://docs.docker.com/install)

## Enable the user management module

The user management module is disabled by default. You can enable each sign in method individually from Space Console. OAuth based signed in methods requires you to enable it OAuth on the concerned platform as well as an additional step.

First start Space Console by running `space-cli start` and create a new project with the database of your choice.

// Image with create project model

The click the `User Management` tab on the side nav. It will open up a screen like this.

// user management screen

You can selectively enable the sign in methods as mentioned earlier. The basic sign in method is pretty straight forward. It enables the REST endpoints for login in and sign up.

Make the required changes as necessary and hit save. That's it. That's how you enable and disable the user management module.

## Next steps
You can check the usage of OAuth based sign in methods [here](https://spaceuptech.com/docs/user-management/oauth).