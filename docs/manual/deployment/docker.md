# Run Space Cloud using Docker

This guide will help you set up Space Cloud via Docker.

## Prerequisites

- [Docker](https://docs.docker.com/install/)


## Step 1: Run Space Cloud

The following command runs `space-cloud` in a docker container and exposes the HTTP and grpc endpoints on ports `4122` and `4124` respectively:  

```bash
docker run -d -p 4122:4122 -p 4124:4124 --name space-cloud \
  -e ADMIN_USER=some-admin \
  -e ADMIN_PASS=some-pass \
  -e ADMIN_SECRET=some-secret \
  spaceuptech/space-cloud:latest
```

Here, `ADMIN_USER` and `ADMIN_PASS` are the credentials to login into Mission Control (Admin UI), whereas `ADMIN_SECRET` is the JWT secret used to authenticate login requests for Mission Control. 

> **Note:** The HTTP and grpc endpoints are available in a secure fashion over SSL on ports `4126` and `4128` respectively.

To expose the HTTP and grpc endpoints of Space Cloud in a secure way via SSL run the following command:
```bash
docker run -d -p 4126:4128 -p 4128:4128 --name space-cloud \
  -v /path/to/ssl-folder:/ssl
  -e ADMIN_USER=some-admin \
  -e ADMIN_PASS=some-pass \
  -e ADMIN_SECRET=some-secret \
  -e SSL_CERT=/ssl/some-ssl.crt \
  -e SSL_KEY=/ssl/some-ssl.key \
  spaceuptech/space-cloud:latest
```

Check if the `space-cloud` container is up and running:
```bash
docker ps
```

## Step 2: Configure Space Cloud

If you exec into docker container of Space Cloud, you can see a `config.yaml` file and a `raft-store` folder would have been generated in the home directory.

Space Cloud needs this config file in order to function. The config file is used to load information like the database to be used, its connection string, security rules, etc. 

Space Cloud has it's own Mission Control (admin UI) to configure all of this in an easy way. 

> **Note:** All changes to the config of `space-cloud` has to be done through the Mission Control only. Changes made manually to the config file will get overwritten. 


### Open Mission Control

Head over to `http://localhost:4122/mission-control` or `https://localhost:4126/mission-control` to open Mission Control depending on how you started `space-cloud`.

> **Note:** Replace `localhost` with the address of your Space Cloud if you are not running it locally. 


## Next Steps

Awesome! We just started Space Cloud using Docker. Next step would be to [set up a frontend/backend project](/docs/setting-up-project/) to use Space Cloud in your preffered langauage. 

Feel free to dive into various modules of Space Cloud:

- Perform CRUD operations using [Database](/docs/database/) module
- Manage files with ease using [File Management](/docs/file-storage) module
- Allow users to sign-in into your app using [User management](/docs/user-management) module
- Write custom logic at backend using [Functions](/docs/functions/) module
- [Secure](/docs/security) your apps

<div class="btns-wrapper">
  <a href="/docs/quick-start/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/setting-up-project/" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
