---
title: "Learn React"
description: Learn how to bind your react app with Space Cloud
date: 2019-07-16T20:01:12+05:30
draft: true
cover:
- "point 1"
- "point 2"
---


Hey guys, its Noorain Your tech bud, and in this tutorial I’m gonna teach you everything you’ll need to build a react app with space cloud and mongodb.

But before we get started, let’s talk about what space cloud really is. Space cloud is an open source web server which provides a realtime data access layer and a full fledged functions mesh for your microservices.

Like Google firebase, it lets you query your database from your frontend (like a react app) or from your backend microservices. And it’s realtime! So all changes get synced between all concerned clients immediately. But unlike firebase, SC can run with mongodb, mysql and postgres with many more dbs yet to come. Basically you could replace the mongodb we’ll be using in this tutorial with mysql or postgres seamlessly. This is the realtime data access layer i was talking about.

This sounds very similar to another super cool project named prisma. But unlike prisma you can use SC directly from the frontend securely. SC has got a super powerful security module which lets you secure database access.

The next one is the functions mesh. Basically it lets you write microservices but instead of exposing functionality as HTTP endpoints, you expose them directly as functions. What this means is that you can invoke these functions, (which are running on your backend) directly from another service or from your frontend. All networking, service discovery and load balancing is completely taken care of. So if you are running two instances of the same microservice, SC will automatically load balance between them... which is pretty cool

All this functionality can be accessed over HTTP, Websockets and gRPC. We have also got client libraries in python, java, node and go to make your life much easier.

Coming back to the tutorial, we’ll be making a realtime todo app. We’ll be able to add todos, mark them as checked, delete them and ofcourse login and sign up. Everything will be realtime. So updates on one client will be reflected on the other ones.

We’ll be starting with downloading and running the space cloud binary. So SC is distributed as a single file which you simply download and run. There is no need to get a bunch of services up or anything.

Next would be configuring SC. It needs to be told stuff like which db to connect to, what’s the connection string… the security rules and stuff like that. You can do this by providing a config file. In this tutorial, however, I’ll be using mission control which is SC’s admin ui.

I’m leaving it upto you to get mongodb up and running as a prerequisite.

Well that’s gonna be our entire backend setup. We don’t need to write any code in node java or whatever. Just running SC should be fine.

The frontent will be a react app. I won’t be covering writing the react app in this series. I already have a todo app hosted on github. What i will cover... is how to use the space api to talk to space cloud and bind the app with mongodb. So its not the UI part, but the backend integration part.

To give you a rough overview. Your react app will intialise the space api and point it to the SC url, which will be localhost. All requests will be sent to SC. SC will authorise our requests, generate an appropriate mongodb query and run it against our db and return the response.

So in this tutorial we’ll be covering most of the concepts you would need to know while making an app with react SC. The goal here is to leave you with enough confidence so you can start using SC in your own projects.

Drop in a comment if you want me to cover a particular topic or you want to give some feedback. I would really appreciate feedback. We gonna be releasing a ton of content so subscribe to stay tuned.

See you in the next video