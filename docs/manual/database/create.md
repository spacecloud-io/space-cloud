# Creating data

You can add data to your app by simply calling `db.insert` on the frontend. Here's a code snippet to do so:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js1">Javascript</a></li>
      <li class="tab col s2"><a href="#java1">Java</a></li>
      <li class="tab col s2"><a href="#python1">Python</a></li>
    </ul>
  </div>
  <div id="js1" class="col s12" style="padding:0">
    <pre>
      <code>
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The todo to be created
const doc = { _id: 1, text: "Star Space Cloud on Github!", time: new Date() };

// Insert the todo in 'todos' collection / table
db.insert("todos").one(doc).then(res => {
    if (res.status === 200) {
      // Todo was created successfully
      return;
    }
  })
  .catch(ex => {
    // Exception occured while processing request
  });
      </code>
    </pre>
  </div>
  <div id="java1" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python1" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

As you would have noticed, the above function is asynchronous in nature. `doc` is the document that you want to insert while `todos` is the name of the collection / table in which you want to insert your `doc`.

You can insert a single document using the `one` method or multiple documents using the `all` method.

## Response

A response object sent by the server contains the **status** fields explained below:

**status:** Number describing the status of the operation. Following values are possible:

- 200 - Operation was successful
- 401 - Request was unauthenticated
- 403 - Request was unauthorized
- 500 - Internal server error

## Adding a single document:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js2">Javascript</a></li>
      <li class="tab col s2"><a href="#java2">Java</a></li>
      <li class="tab col s2"><a href="#python2">Python</a></li>
    </ul>
  </div>
  <div id="js2" class="col s12" style="padding:0">
    <pre>
      <code>
// Create a single todo
const doc = { _id: 1, text: 'Star Space Cloud on Github!', time: new Date()};
db.insert('todos').one(doc).then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="java2" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python2" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

## Adding multiple documents simultaneously:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js3">Javascript</a></li>
      <li class="tab col s2"><a href="#java3">Java</a></li>
      <li class="tab col s2"><a href="#python3">Python</a></li>
    </ul>
  </div>
  <div id="js3" class="col s12" style="padding:0">
    <pre>
      <code>
const docs = [
  { _id: 1, text: 'Star Space Cloud on Github!', time: new Date()},
  { _id: 1, text: 'Fork Space Cloud on Github!', time: new Date()}
];
db.insert('todos').all(docs).then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="java3" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python3" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

## Next steps

Now you know how to insert data into a database using Space Cloud. So let's check how to query it back.

<div class="btns-wrapper">
  <a href="/docs/database/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/read" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
