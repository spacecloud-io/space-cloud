# Add Space Cloud to your JavaScript project

Follow this guide to use Space Cloud APIs in your web app or any Javascript/Node.js project.

## Step 1: Install Space Cloud API
**Install via npm:**
```bash
$ npm install space-api --save
```

**Or import as a stand alone library:**
```js
<script src="https://spaceuptech.com/downloads/libraries/space-api.js"></script>
```

## Step 2: Create an API instance

An `api` instance of Space Cloud on the frontend will help you talk to `space-cloud` binary and perform backend operations directly from the frontend. 

The `API` constructor takes two parameters: 
- **PROJECT_ID:** Unqiue identifier of a project. It's derived by converting your project name to lowercase and replacing all spaces to hiphens. For example `Todo App` becomes `todo-app`.
- **SPACE_CLOUD_URL:** This is the url of your `space-cloud` binary. It's `http://localhost:4122` or `https://localhost:4126` for HTTP and HTTPS endpoints respectively.

> **Note:** Replace `localhost` with the address of your Space Cloud if you are not running it locally. 

**For ES6:**
```js
import { API } from 'space-api';

const api = new API('todo-app', 'http://localhost:4122');
```

**For ES5/CommonJS:**
```js
const { API } = require('space-api');

const api = new API('todo-app', 'http://localhost:4122');
```

**For stand alone:**
```js
var api = new Space.API("todo-app", "http://localhost:4122");
```


## Step 3: Create a DB instance

The `api` instance created above will help you to directly use `fileStorage` and `functions` modules. However, to use `crud`, `realTime` and `auth` modules you will also need to create a `db` instance.

> **Note:** You can use multiple databases in the same project. (For eg. MongoDB and MySQL)

**For MongoDB:**
```js
const db = api.Mongo();
```

**For PostgreSQL:**

> **Note:** This can also be used for any other database that is PostgreSQL compatible (For eg. CockroachDB, Yugabyte etc.)
```js
const db = api.Postgres();
```

**For MySQL:**

> **Note:** This can also be used for any other database that is MySQL compatible (For eg. TiDB)
```js
const db = api.MySQL();
```

## Next steps
Great! Since you have initialized the `api` and `db` instance you can start building apps with `space-cloud`. Check out these modules to explore all that you can do with `space-cloud`:
- Perform CRUD operations using [Database](/docs/database/) module
- [Realtime](/docs/realtime/) data sync across all devices
- Manage files with ease using [File Management](/docs/file-storage) module
- Allow users to sign-in into your app using [User management](/docs/user-management) module
- Write custom logic at backend using [Functions](/docs/functions/) module
- [Secure](/docs/security) your apps

<div class="btns-wrapper">
  <a href="/docs/setting-up-project/" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>