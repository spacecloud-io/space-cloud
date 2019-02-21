# Space Cloud
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

<a href="https://discord.gg/ypXEEBr"><img src="https://img.shields.io/badge/chat-discord-brightgreen.svg?logo=discord&style=flat"></a>
<a href="https://twitter.com/intent/follow?screen_name=spaceuptech"><img src="https://img.shields.io/badge/Follow-spaceuptech-blue.svg?style=flat&logo=twitter"></a>

Space Cloud is an open source, high performance web service which provides instant Realtime APIs on the database of your choice. Build Internet Scale apps with the agility of a prototype!

## Table of Contents

- [About Space Cloud](#about-space-cloud)
- [What makes Space Cloud unique](#what-makes-space-cloud-unique)
- [Documentation](#documentation)
- [Getting started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Download Space Cloud](#download-space-cloud)
    - [Download config file](#download-the-config-file)
    - [Start Space Cloud](#start-space-cloud)
    - [Download TODO App](#download-the-todo-app)
- [Support & Troubleshooting](#support--troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## About Space Cloud

> Note: Space Cloud is in still in Beta.

Space Cloud is an open-source, high performance web engine which lets you create instant Realtime APIs on most of the databases out there. Written in Golang, it provides a high throughput layer for your backend services. It's completely unopinionated and works with the tech stack of your choice.

Space Cloud is purpose-built to power fast-growing, realtime online services on public, private and hybrid clouds requiring global scale at prototyping agility. Built with extensibility in mind, it provides APIs for you to extend the platform in the form of functions in any language.

## What makes Space Cloud unique?
Space Cloud is a single web engine which integrates with all you back end technologies and exposes them as easy-to-use REST APIs. You can leverage the power of the tools you already love without changing a single line of code. No migrations are necessary. Add new platforms or technologies as needed without having to worry about how to use them.

In a nutshell, Space Cloud provides you:
- Instant Realtime APIs to consume directly from your frontend.
- Authentication and authorization baked in by default.
- Freedom from vendor lock ins.
- Flexibility to work with the tech stack of your choice.
- Various pre-built modules such as User Management, Realtime CRUD and File Storage.

## Documentation
We are working hard to document every aspect of Space Cloud to give you the best onboarding experience. Here are links to the various docs we have:
- Space Cloud (Coming soon!)
- Client APIs:
    - [Javascript client](https://github.com/spaceuptech/space-api-js/wiki)
    - Java client (Coming soon!)

## Getting started
Let's see how to build an end-to-end todo app using Space Cloud

### Prerequisites
- [MongoDB database](https://docs.mongodb.com/manual/installation/)

### Download Space Cloud
You need to download the `space-cloud` binary for your operating system or you could build it directly from its source code. You need go version 1.11.2 or later to build it from source.

Download the binary for your OS from here:
- [Mac](https://spaceuptech.com/downloads/darwin/space-cloud.zip)
- [Linux](https://spaceuptech.com/downloads/linux/space-cloud.zip)
- [Windows](https://spaceuptech.com/downloads/windows/space-cloud.zip)

Make the `space-cloud` binary executable and add it to you `PATH`.

### Download the config file
Space Cloud needs a config file in order to function properly. It relies on the config file to load information like the database connection string, security rules, etc. 

You can find a sample config for the todo app [here](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/basic-todo-app/config.yaml). Feel free to explore the file.

### Start Space Cloud
You can start `space-cloud` with the following command. Make sure mongo db is running before this step
```
space-cloud run --config config.yml
```

That's it. Your backend is up and running!

### Download the TODO App
We have built a simple todo app using HTML and vanilla Javascript which works with the backend you have just created. You can find it [here](https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/basic-todo-app/index.html).

Feel free to play around with it to explore all the capabilities of Space Cloud.

## Support & Troubleshooting

The documentation and community will help you troubleshoot most issues. If you have encountered a bug or need to get in touch with us, you can contact us using one of the following channels:

* Support & feedback: [Discord](https://discord.gg/ypXEEBr)
* Issue & bug tracking: [GitHub issues](https://github.com/spaceuptech/space-cloud/issues)
* Follow product updates: [@spaceuptech](https://twitter.com/spaceuptech)

## Contributing
Space Cloud is a young project. We'd love to have you on board if you wish to contribute. To help you get started, here are a few areas you can help us with:
- Writing the documentation
- Making sample apps in React, Angular, Android, and any other frontend tech you can think of
- Deciding the road map of the project
- Creating issues for any bugs you find
- And of course, with code for bug fixes and new enhancements

## License
Space Cloud is [Apache 2.0 licensed](https://github.com/spaceuptech/space-cloud/blob/master/LICENSE).
