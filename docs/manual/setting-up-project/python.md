# Add Space Cloud to your Python project

Follow this guide to use Space Cloud APIs in your python app.

## Step 1: Install Space Cloud API
**Install via pip:**
```bash
$ pip install space-api-py
```

## Step 2: Create a Client instance

A client instance of Space Cloud will help you talk to `space-cloud` binary and perform backend operations directly from the frontend.

The `API` constructor takes two parameters: 
- **PROJECT_ID:** Unqiue identifier of a project. It's derived by converting your project name to lowercase and replacing all spaces to hiphens. For example `Todo App` becomes `todo-app`.
- **SPACE_CLOUD_URL:** This is the url of your `space-cloud` binary. It's `localhost:4124` or `localhost:4128` for grpc and grpc secure endpoints respectively.

> **Note:** Replace `localhost` with the address of your Space Cloud if you are not running it locally. 

```python
from space_api import API

api = API('books-app', 'localhost:4124')
```


## Step 3: Create a DB instance

The `api` instance created above will help you to directly use `fileStorage` and `functions` modules. However, to use `crud`, `realTime` and `auth` modules you will also need to create a `db` instance.

> **Note:** You can use multiple databases in the same project. (For eg. MongoDB and MySQL)

**For MongoDB:**
```python
db = api.mongo()
```

**For PostgreSQL:**

> **Note:** This can also be used for any other database that is PostgreSQL compatible (For eg. CockroachDB, Yugabyte etc.)
```python
db = api.postgres()
```



**For MySQL:**

> **Note:** This can also be used for any other database that is MySQL compatible (For eg. TiDB)
```python
db = api.my_sql()
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