# Space Cloud Documentation

## What is Space Cloud?

> Note: Space Cloud is in still in Beta.

Space Cloud is an open-source, high performance web engine which lets you create instant REST APIs on most of the databases out there. Written in Golang, it provides a high throughput layer for your backend services. It's completely unopiniated and works with the tech stack of your choice.

Space Cloud is purpose-built to power fast-growing, realtime online services on public, private and hybrid clouds requiring global scale at prototyping agility. Built with extensibility in mind, it provides APIs for you to extend the platform in the form of functions in any language.

## What makes Space Cloud unique?
Space Cloud is a single web engine which integrates with all you back end technologies and exposes them as easy-to-use REST APIs. You can leverage the power of the tools you already love without changing a single line of code. No migrations are necessary. Add new platforms or technologies as needed without having to worry about how to use them.

In a nutshell, Space Cloud provides you:
- Instant REST APIs to consume directly from your frontend.
- Authentication and authorization backed in by default.
- Freedom from vendor lock ins.
- Flexibility to work with the tech choice of your choice.
- Various pre-built modules such as [User Management](https://spaceuptech.com/docs/user-mangement/), [Realtime CRUD](https://spaceuptech.com/docs/realtime) and [File Storage](https://spaceuptech.com/docs/file-storage).

## Commonly used terms
There are a couple of terms we use to the refer the various pieces of Space Cloud. Let's clear them up right away.

### Space Cloud - [(Product Page)](https://spaceuptech.com)
Space Cloud is the name for the umbrella project. Together with some other projects it serves to provide instant REST APIs on any database. Space Cloud has the following components
- **Space Console** - [(Github Repo)](https://github.com/spaceuptech/space-console): Space Console is the visual tool the create and configure Space Cloud projects. All project config is stored in a simple json file.
- **A CLI** - [(Github Repo)](https://github.com/spaceuptech/space-cli): Space CLI (also referred to as `space-cli`) is a command line utility to simplify the process of creating and deploying Space Cloud projects. It's written in Node.js and is installed via npm.
- **A Binary** - [(Github Repo)](https://github.com/spaceuptech/space-cloud): This is the binary which is physically deployed to run Space Cloud. Referred to as `space-cloud`, it creates the REST endpoints and makes connections to the appropriate backend technologies depending on the project configuration.

## What's next?
- Head straight to our [getting started guide](https://spaceuptech.com/docs/getting-started).
- New to Space Cloud? Checkout the [tutorial](https://spaceuptech.com/tutorials) instead.

[next](https://spaceuptech.com/docs/getting-started) >>

