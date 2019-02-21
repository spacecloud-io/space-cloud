# Database Module

The database module is the core of Space Cloud. It provides instant REST APIs on any database out there. it allows you to perfrom CRUD operations on the database directly from the frontend. The API loosely follows the Mongo DB DSL. In addition to that, it expects a JWT token in the `Authorization` header. This JWT token is used along with user defined security rules to enforce authentication and authorization.

By CRUD I mean Create, Update and Delete operations. These are the most basic operations that one can perform on a database. In addition to that, we offer a flexible query language (based on the Mongo DB query DSL) to slice and dice data as needed.

Currently the database module supports the following databases:
- Mongo DB
- MySQL and MySQL compatible databases
- Postgres and Postgres compatible databases

## Prerequisites
- [Space CLI](https://spaceuptech.com/docs/getting-started/space-cli)
- [Docker](https://docs.docker.com/install)

## Exploring the database module

The database module requires a `primary database` to begin with. This is the database selected while creating a new project. All the other modules default to using this database to store their metadata (if any).

// Create project model goes here

Once you have created a project, you can head over to the `Database` tab in the side nav. You'll be greeted with a screen which looks something like this.

// Image of db rules

This is where you configure the security rules I spoke about earlier. The upper left dropdown is to select the database you are currently configuring. Each database you add in your project (yes you can use multiple databases in the same project) needs to be configured separately.

You can hit on the add table / collection button to add another table to Space Cloud. Note, this does not actually create a table / collection in the database. All tables need to be created by the user. This can be done either by logging into the database directly or using the Data Explorer feature in Space Cloud (Coming Soon!). Space Cloud only exposes those tables / collections via REST which have been added to it. This prevents accidentally exposing certain hidden or internal tables to the world.

The connection string in the usual connection string used to connect to the database. It's database specific, i.e. each database may have it's own format for writing connection string. Showing below the connection string formats for the supported databases:
- Mongo DB: mongo://[username:password@]host1[:port1][,...hostN[:portN]]][/[database][?options]]
- MySQL: mysql://[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
- Postgres: postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]

At this point I'll urge you to go ahead and add a table. Let's call it `posts` for now.

// Image with the security rules

The box beside the tables list is something we like to call the security rules. The icon right above it briefly describes it's usage. To know more about how security rules works, you can head over to the [security page](https://spaceuptech.com/docs/security).

## Adding or removing databases

Remember I said you can use multiple databases in a single Space Cloud project. You do that by adding secondary databases. Head over to the `Config` tab in the database module. You'll find all the supported databases right here. Feel free to add or remove databases. Just don't forget to configure them in the `Rules` tab.

## Next steps

Now you know the basics of the database module. The next step would be diving deeper into the security rules and it's structure. Let's make sure that the apps we build are secure!

You can also check out the [API docs](https://spaceuptech.com/docs/api) to start building your app right away.
