---
title: "Install Space Cloud"
description: The first step to using Space Cloud is to install it
date: 2019-07-16T20:11:26+05:30
draft: true
index: 2
---

Hey guys, its Noorain Your tech bud, and in this tutorial I’m gonna teach you everything you’ll need to build a react app with space cloud and mongodb.

But before we get started, let’s talk about what space cloud really is. Space cloud is an open source web server which provides a realtime data access layer and a full fledged functions mesh for your microservices.

Like Google firebase, it lets you query your database from your frontend (like a react app) or from your backend microservices. And it’s realtime! So all changes get synced between all concerned clients immediately. But unlike firebase, SC can run with mongodb, mysql and postgres with many more dbs yet to come. Basically you could replace the mongodb we’ll be using in this tutorial with mysql or postgres seamlessly. This is the realtime data access layer i was talking about.

This sounds very similar to another super cool project named prisma. But unlike prisma you can use SC directly from the frontend securely. SC has got a super powerful security module which lets you secure database access.

The next one is the functions mesh. Basically it lets you write microservices but instead of exposing functionality as HTTP endpoints, you expose them directly as functions. What this means is that you can invoke these functions, (which are running on your backend) directly from another service or from your frontend. All networking, service discovery and load balancing is completely taken care of. So if you are running two instances of the same microservice, SC will automatically load balance between them... which is pretty cool
