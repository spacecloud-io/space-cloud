# Static Hosting

The static module of Space Cloud provides an easy way to host your static files. It can also act as an reverse proxy to other upstream servers in your system. However, it's not meant to replace full blown web servers like Nginx or Apache. In fact, for more advanced use cases, Space Cloud would come behind a reverse proxy server like Nginx or Apache. Space Cloud can come in handy in simple use cases like hosting a website alongside your backend (Space Cloud), where something like Nginx or Apache will be an overkill or difficult to start with.

## Configuring static module

The configurations for static module goes inside the `static` section in `modules`. It contains two options:
- **enabled:** - To enable the static module
- **routes:** An array of `route` to host

A `route` consists of the following fields:
- **host:** Name of the host
- **prefix:** Prefix of the incoming request
- **path:** Path of the folder containing static resources (used for static hosting)
- **proxy:** Address of the upstream server to proxy incoming request (used for reverse proxy)


## Serve static resources
Here's an sample config of how to host static resources via Space Cloud:

```yaml
  static:
    routes:
      - host: console.spaceuptech.com
        prefix: /home
        path: /public/console
```
The above configuration will host the folder the static resources at `/public/console` for any requests at console.spaceupetch.com starting with `/home`

## Reverse proxy

Here's how Space Cloud can be configured to act as a reverse proxy server to other servers in your system:

```yaml
  static:
    routes:
      - host: spaceuptech.com
        prefix: /v1/
        proxy: http://localhost:8090/v1/
```


## Serve SPAs

Serving an SPA like ReactJS or AngularJS is different from serving a static website. Space Cloud can't natively host an SPA as of now. However, you can serve them via something like [serve](https://www.npmjs.com/package/serve) and use Space Cloud as a reverse proxy to them.

## Next steps

Great, you have learned how to host static resources via Space Cloud. You can continue to see how to [deploy Space Cloud](/docs/deploy/overview).

<div class="btns-wrapper">
  <a href="/docs/security/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/deploy/overview" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
