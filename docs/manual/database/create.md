# Creating data

You can add data to your app by simply calling `db.insert` on the frontend. Here's a code snippet to add a single document/record to your app:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#insert-js">Javascript</a></li>
      <li class="tab col s2"><a href="#insert-java">Java</a></li>
      <li class="tab col s2"><a href="#insert-python">Python</a></li>
    </ul>
  </div>
  <div id="insert-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The todo to be created
const doc = { _id: 1, text: "Star Space Cloud on Github!", time: new Date() };

db.insert("todos").doc(doc).apply()
  .then(res => {
    // Insert the todo in 'todos' collection / table
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
  <div id="insert-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
 <div id="insert-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
# Python client coming soon!
      </code>
    </pre>
  </div>
</div>

As you would have noticed, the `insert` method is asynchronous in nature. It takes the name of the concerned collection/table. The `doc` method takes an object to be inserted. The `apply` method actually triggers the given request to `space-cloud` and returns a promise.

## Adding multiple documents simultaneously:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#insertmany-js">Javascript</a></li>
      <li class="tab col s2"><a href="#insertmany-java">Java</a></li>
      <li class="tab col s2"><a href="#insertmany-python">Python</a></li>
    </ul>
  </div>
  <div id="insertmany-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
const docs = [
  { _id: 1, text: "Star Space Cloud on Github!", time: new Date() },
  { _id: 1, text: "Fork Space Cloud on Github!", time: new Date() }
];
db.insert('todos').docs(docs).apply().then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="insertmany-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
 <div id="insertmany-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
# Python client coming soon!
      </code>
    </pre>
  </div>
</div>

The `docs` method takes an array of objects to be inserted.

## Response

On response from the server, the callback passed to the `then` method is called with the response object as described below:

```
{
  "status": "number" // Status of the operation
}
```

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
