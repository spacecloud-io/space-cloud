# Deleting Data

You can delete data in your app by simply calling `db.delete` on the frontend. Here's a code snippet to delete all documents matching a specific condition:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#delete-js">Javascript</a></li>
      <li class="tab col s2"><a href="#delete-java">Java</a></li>
      <li class="tab col s2"><a href="#delete-python">Python</a></li>
    </ul>
  </div>
  <div id="delete-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API, and, or, cond } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The condition to be matched
const condition = cond("userId", "==", "user1");

// Delete all todos of a particular user which match a particular condition
db.delete("todos").where(condition).apply()
  .then(res => {
    if (res.status === 200) {
      // Documents were deleted successfully
      return;
    }
  })
  .catch(ex => {
    // Exception occured while processing request
  });
      </code>
    </pre>
  </div>
   <div id="delete-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
  <div id="delete-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
# Python client coming soon!
      </code>
    </pre>
  </div>
</div>

As you would have noticed, the `delete` method is asynchronous in nature. It takes the name of the concerned collection/table and deletes all the matching documents. The `apply` method actually triggers the given request to `space-cloud` and returns a promise.

## Deleting a single document:

> Note: `deleteOne` method is not available for SQL databases.

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#delete-one-js">Javascript</a></li>
      <li class="tab col s2"><a href="#delete-one-java">Java</a></li>
      <li class="tab col s2"><a href="#delete-one-python">Python</a></li>
    </ul>
  </div>
  <div id="delete-one-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
db.deleteOne('todos').where(cond('_id', '==', 1)).apply().then(res => ...).catch(ex => ...);      
      </code>
    </pre>
  </div>
   <div id="delete-one-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
  <div id="delete-one-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
# Python client coming soon!      
      </code>
    </pre>
  </div>
</div>

The `deleteOne` method takes the name of the concerned table/collection. It deletes a single document matching the condition. If no matching document is found, it returns an error (`res.status` - 400).
## Delete documents selectively

You can selectively delete only a few documents which you desire and leave the rest by using `where` clause. The `where` method accepts a `condition` object. After validation, `space-cloud` generates a database specific query. The documents or rows which match this query are deleted by the update operations described above.

### Specifying a single condition

The `cond` function is used to specify a single condition as shown below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#cond-js">Javascript</a></li>
      <li class="tab col s2"><a href="#cond-java">Java</a></li>
      <li class="tab col s2"><a href="#cond-python">Python</a></li>
    </ul>
  </div>
  <div id="cond-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Syntax
var op1 = 'field_name';
var operator = '== | != | > | < | >= | <= | in | notIn';
var op2 = 'value';
const condition = cond(op1, operator, op2);

// Example
const condition = cond('_id', '==', 1);
db.delete('todos').where(condition).apply().then(res => ...);   
      </code>
    </pre>
  </div>
   <div id="cond-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
  <div id="cond-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
# Python client coming soon!      
      </code>
    </pre>
  </div>
</div>

The operators allowed are:

| Operator | Description                                       |
|:---------|:--------------------------------------------------|
| ==       | Passes if `op1` is equal to `op2`                 |
| !=       | Passes if `op1` is not equal to `op2`             |
| >        | Passes if `op1` is greater than `op2`             |
| <        | Passes if `op1` is lesser than `op2`              |
| >=       | Passes if `op1` is greater than or equal to `op2` |
| <=       | Passes if `op1` is lesser than or equal to `op2`  |
| in       | Passes if `op1` is in `op2`                       |
| notIn    | Passes if `op1` is not in `op2`                   |

### Combining multiple conditions

A single condition is often not enough to delete the data you desire. You might need to `and` / `or` several conditions together. For e.g., you want to delete only those posts which are of a particular author **and** of a particular category. The following code snippet shows how you can use and and or:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#multiple-cond-js">Javascript</a></li>
      <li class="tab col s2"><a href="#multiple-cond-java">Java</a></li>
      <li class="tab col s2"><a href="#multiple-cond-python">Python</a></li>
    </ul>
  </div>
  <div id="multiple-cond-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// AND syntax
const condition = and(...conditions)

// Example
const condition = and(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.delete('todos').where(condition).apply().then(res => ...)

// OR syntax
const condition = or(...conditions)

// Example
const condition = or(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.delete('todos').where(condition).apply().then(res => ...);      
      </code>
    </pre>
  </div>
   <div id="multiple-cond-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
  <div id="multiple-cond-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
# Python client coming soon!      
      </code>
    </pre>
  </div>
</div>

## Response

On response from the server, the callback passed to the `then` method is called with the response object as described below:

```
{
  "status": "number" // Status of the operation
}
```

## Next steps

Now you know all the operations of CRUD module. So let's take a deeper dive into configuring the database module of Space Cloud.

<div class="btns-wrapper">
  <a href="/docs/database/update" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/config" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
