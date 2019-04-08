# Updating Data

You can update / modify existing data in your app by simply calling `db.update` on the frontend. Here's a code snippet to update all documents matching a specific condition:

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
import { API, and, or, cond } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The condition to be matched
const condition = cond("author", "==", 'author1');

// Update the todo
db.update("todos")
  .where(condition)
  .set({ text: "Fork Space Cloud on Github right now" }).all()
  .then(res => {
    if (res.status === 200) {
      // Documents were updated successfully
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

As you would have noticed, the above function is asynchronous in nature. `todos` is the name of the collection / table which contains the docs that you want to update.

You can update a single document using `one` method. Multiple documents can be updated using `all` method.

> Note: `one` method is not available for SQL databases.

## Update operations

You can perform different types of update operations like set, push, rename, etc. on your data. Following are the different types of update operations:

> Note: In SQL databases, only `set` operation is available.

### Set operation

You can set the value of any field / column in your data by using `set` method like below: 

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
// Set text of todo by new value
db.update('todos').where(cond('_id_', '==', 1))
  .set({text: 'Fork Space Cloud on Github'}).one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java2" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python2" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `set` method accepts an object of key-value pairs where key is the field / column name whereas value is the new value with which you want to update the field / column. 

### Push operation

You can push an element to an array in a document by using the `push` method like below:

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
// Add a new category for a todo
db.update('todos').where(cond('_id_', '==', 1))
  .push({categories: 'some-category'}).one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java3" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python3" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `push` method accepts an object of key-value pairs where key is the field name of the array whereas value is the new value which you want to push into that array.

### Remove operation

You can remove / delete a field in a document by using `remove` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js4">Javascript</a></li>
      <li class="tab col s2"><a href="#java4">Java</a></li>
      <li class="tab col s2"><a href="#python4">Python</a></li>
    </ul>
  </div>
  <div id="js4" class="col s12" style="padding:0">
    <pre>
      <code>
// Remove categories and time field
db.update('todos').where(cond('_id_', '==', 1))
  .remove('categories', 'time').one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java4" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python4" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `remove` method accepts multiple inputs each being the name of a field you want to remove.  

### Rename operation

You can rename the name of a field in a document by using `rename` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js5">Javascript</a></li>
      <li class="tab col s2"><a href="#java5">Java</a></li>
      <li class="tab col s2"><a href="#python5">Python</a></li>
    </ul>
  </div>
  <div id="js5" class="col s12" style="padding:0">
    <pre>
      <code>
// Rename 'text' field to 'content'
db.update('todos').where(cond('_id_', '==', 1))
  .rename({text: 'content'}).one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java5" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python5" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `rename` method accepts an object of key-value pairs where key is the current name of a field while value (string) is the new name that you want to assign to that field.

### Increment operation

You can increment / decrement the value of a integer field in your data by using the `inc` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js6">Javascript</a></li>
      <li class="tab col s2"><a href="#java6">Java</a></li>
      <li class="tab col s2"><a href="#python6">Python</a></li>
    </ul>
  </div>
  <div id="js6" class="col s12" style="padding:0">
    <pre>
      <code>
// Increment 'likes' by 3
db.update('todos').where(cond('_id_', '==', 1))
  .inc({likes: 3}).one().then(res => ...);

// Decrement 'likes' by 2
db.update('todos').where(cond('_id_', '==', 1))
  .inc({likes: -2}).one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java6" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python6" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `inc` method accepts an object of key-value pairs where key is the name of the field whose value is to be incremented whereas value is the amount by which the value has to be incremented. As you would have noticed you can decrement a value by using negative integers.

### Multiply operation

You can multiply the value of a integer field in your data by using the `mul` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js7">Javascript</a></li>
      <li class="tab col s2"><a href="#java7">Java</a></li>
      <li class="tab col s2"><a href="#python7">Python</a></li>
    </ul>
  </div>
  <div id="js7" class="col s12" style="padding:0">
    <pre>
      <code>
// Multiply 'likes' by 10
db.update('todos').where(cond('_id_', '==', 1))
  .mul({likes: 10}).one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java7" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python7" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `inc` method accepts an object of key-value pairs where key is the name of the field whose value is to be multiplied whereas value is the amount by which the value has to be multiplied.

### Max operation

Sometimes you might want to update a number in your document with a new value only if the new value is greater than the existing value. This can be acheived by using `max` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js8">Javascript</a></li>
      <li class="tab col s2"><a href="#java8">Java</a></li>
      <li class="tab col s2"><a href="#python8">Python</a></li>
    </ul>
  </div>
  <div id="js8" class="col s12" style="padding:0">
    <pre>
      <code>
// Updates 'likes' if it was lesser than 50
db.update('todos').where(cond('_id_', '==', 1))
  .max({likes: 50}).one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java8" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python8" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

### Min operation

Sometimes you might want to update a number in your document with a new value only if the new value is lesser than the existing value. This can be acheived by using `min` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js9">Javascript</a></li>
      <li class="tab col s2"><a href="#java9">Java</a></li>
      <li class="tab col s2"><a href="#python9">Python</a></li>
    </ul>
  </div>
  <div id="js9" class="col s12" style="padding:0">
    <pre>
      <code>
// Updates 'likes' if it was greater than 50
db.update('todos').where(cond('_id_', '==', 1))
  .min({likes: 50}).one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java9" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python9" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

### Current timestamp operation

You can update a field with the value of current timestamp by using the `currentTimestamp` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js10">Javascript</a></li>
      <li class="tab col s2"><a href="#java10">Java</a></li>
      <li class="tab col s2"><a href="#python10">Python</a></li>
    </ul>
  </div>
  <div id="js10" class="col s12" style="padding:0">
    <pre>
      <code>
// Update 'lastUpdated' with current timestamp 
db.update('todos').where(cond('_id_', '==', 1))
  .currentTimestamp('lastUpdated').one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java10" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python10" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `currentTimestamp` method accepts multiple inputs each being the name of the field you want to update with current timestamp.

### Current date operation

You can update a field with the value of current date by using the `currentDate` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js11">Javascript</a></li>
      <li class="tab col s2"><a href="#java11">Java</a></li>
      <li class="tab col s2"><a href="#python11">Python</a></li>
    </ul>
  </div>
  <div id="js11" class="col s12" style="padding:0">
    <pre>
      <code>
// Update 'lastUpdated' with current date 
db.update('todos').where(cond('_id_', '==', 1))
  .currentDate('lastUpdated').one().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="java11" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python11" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `currentDate` method accepts multiple inputs each being the name of the field you want to update with current date.

## Update documents selectively

You can selectively update only a few documents which you desire and leave the rest by using `where` clause. The `where` method accepts a `condition` object. After validation, `space-cloud` generates a database specific query. The documents or rows which match this query are updated by the update operations described above.

### Specifying a single condition

The `cond` function is used to specify a single condition as shown below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js12">Javascript</a></li>
      <li class="tab col s2"><a href="#java12">Java</a></li>
      <li class="tab col s2"><a href="#python12">Python</a></li>
    </ul>
  </div>
  <div id="js12" class="col s12" style="padding:0">
    <pre>
      <code>
// Syntax
var op1 = 'field_name';
var operator = '== | != | > | < | >= | <= | in | notIn';
var op2 = 'value';
const condition = cond(op1, operator, op2);

// Example
const condition = cond('_id', '==', 1);
db.update('todos').where(condition)
  .set({text: 'Fork Space Cloud on Github'}).all().then(res => ...)
      </code>
    </pre>
  </div>
  <div id="java12" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python12" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The operators allowed are:

- **==    :** Passes if `op1` is equal to `op2`
- **!=    :** Passes if `op1` is not equal to `op2`
- **>     :** Passes if `op1` is greater than `op2`
- **<     :** Passes if `op1` is lesser than `op2`
- **>=    :** Passes if `op1` is greater than or equal to `op2`
- **<=    :** Passes if `op1` is lesser than or equal to `op2`
- **in    :** Passes if `op1` is in `op2`
- **notIn :** Passes if `op1` is not in `op2`

### Combining multiple conditions

A single condition is often not enough to update the data you desire. You might need to `and` / `or` several conditions together. For e.g., you want to update only those posts which are of a particular author **and** of a particular category. The following code snippet shows how you can use and and or:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js13">Javascript</a></li>
      <li class="tab col s2"><a href="#java13">Java</a></li>
      <li class="tab col s2"><a href="#python13">Python</a></li>
    </ul>
  </div>
  <div id="js13" class="col s12" style="padding:0">
    <pre>
      <code>
// AND syntax
const condition = and(...conditions)

// Example
const condition = and(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.update('todos').where(condition)
  .set({text: 'Fork Space Cloud on Github'}).all().then(res => ...)

// OR syntax
const condition = or(...conditions)

// Example
const condition = or(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.update('todos').where(condition)
  .set({text: 'Fork Space Cloud on Github'}).all().then(res => ...)
      </code>
    </pre>
  </div>
  <div id="java13" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python13" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

## Response

A response object sent by the server contains the **status** fields explained below:

**status:** Number describing the status of the operation. Following values are possible:

- 200 - Operation was successful
- 401 - Request was unauthenticated
- 403 - Request was unauthorized
- 500 - Internal server error


## Updating a single document:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js14">Javascript</a></li>
      <li class="tab col s2"><a href="#java14">Java</a></li>
      <li class="tab col s2"><a href="#python14">Python</a></li>
    </ul>
  </div>
  <div id="js14" class="col s12" style="padding:0">
    <pre>
      <code>
// Updating a single todo
db.update('todos').where(cond('_id', '==', 1)).
  .set({text: 'Fork Space Cloud on Github'}).one().then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="java14" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python14" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

## Updating multiple documents at once:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js15">Javascript</a></li>
      <li class="tab col s2"><a href="#java15">Java</a></li>
      <li class="tab col s2"><a href="#python15">Python</a></li>
    </ul>
  </div>
  <div id="js15" class="col s12" style="padding:0">
    <pre>
      <code>
// Updating all todos of category 'frontend'
db.update('todos').where(cond('category', '==', 'frontend')).
  .set({text: 'Fork Space Cloud on Github'}).all().then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="java15" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python15" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

## Upserting a document:

Sometimes you might want to create a document or update it if it already exists. You can acheive this by using the `upsert` method like below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js16">Javascript</a></li>
      <li class="tab col s2"><a href="#java16">Java</a></li>
      <li class="tab col s2"><a href="#python16">Python</a></li>
    </ul>
  </div>
  <div id="js16" class="col s12" style="padding:0">
    <pre>
      <code>
// Create a todo or update it
db.update('todos').where(cond('_id', '==', 1))
  .set({text: 'Fork Space Cloud on Github'}).upsert().then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="java16" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python16" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

> Note: `upsert` method is only available for Mongo DB.

The above example will update a todo of _id = 1 with the text - 'Fork Space Cloud on Github' if a todo with _id = 1 already exists. Otherwise it will create a new todo - { _id: 1, text: 'Fork Space Cloud on Github' }


## Next steps

Now you know how to update data in a database using Space Cloud. So let's check how to query it back.

<div class="btns-wrapper">
  <a href="/docs/database/read" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/delete" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
