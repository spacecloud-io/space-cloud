# Deleting Data

You can delete data in your app by simply calling `db.delete` on the front end. Here's a code snippet to delete all documents matching a specific condition:

```js
import { API, and, or, cond } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The condition to be matched
const condition = cond("userId", "==", 'user1');

// Delete all todos of a particular user
db.delete("todos")
  .where(condition).all()
  .then(res => {
    if (res.status === 200) {
      // Documents were deleted successfully
      return;
    }
  })
  .catch(ex => {
    // Exception occured while processing request
  });
```

As you would have noticed, the above function is asynchronous in nature. `todos` is the name of the collection / table which contains the docs that you want to delete.

You can delete a single document using `one` method. Multiple documents can be deleted using `all` method.

> Note: `one` method is not available for SQL databases.

## Delete documents selectively

You can selectively delete only a few documents which you desire and leave the rest by using `where` clause. The `where` method accepts a `condition` object. After validation, `space-cloud` generates a database specific query. The documents or rows which match this query are deleted by the update operations described above.

### Specifying a single condition

The `cond` function is used to specify a single condition as shown below:

```js
// Syntax
var op1 = 'field_name';
var operator = '== | != | > | < | >= | <= | in | notIn';
var op2 = 'value';
const condition = cond(op1, operator, op2);

// Example
const condition = cond('_id', '==', 1);
db.delete('todos').where(condition).all().then(res => ...)
```

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

A single condition is often not enough to delete the data you desire. You might need to `and` / `or` several conditions together. For e.g., you want to delete only those posts which are of a particular author **and** of a particular category. The following code snippet shows how you can use and and or:

```js
// AND syntax
const condition = and(...conditions)

// Example
const condition = and(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.delete('todos').where(condition).all().then(res => ...)

// OR syntax
const condition = or(...conditions)

// Example
const condition = or(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.delete('todos').where(condition).all().then(res => ...)
```

## Response

A response object sent by the server contains the **status** fields explained below:

**status** : Number describing the status of the operation. Following values are possible:

- 200 - Operation was successful
- 401 - Request was unauthenticated
- 403 - Request was unauthorized
- 500 - Internal server error


## Deleting a single document:
```js
// Deleting a single todo
db.delete('todos').where(cond('_id', '==', 1)).one().then(res => ...).catch(ex => ...);
```

## Deleting multiple documents at once:
```js
// Deleting all todos of a particular user
db.delete('todos').where(cond("userId", "==", 'user1')).all().then(res => ...).catch(ex => ...);
```

## Next steps

Now you know all the operations of CRUD module. So let's take a deeper dive into configuring the database module of Space Cloud

<div class="btns-wrapper">
  <a href="/docs/database/update" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/config" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
