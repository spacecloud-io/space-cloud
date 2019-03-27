# Reading Data

You can query data from your database by simply calling `db.get` on the front end. Here's a code snippet to fetch a single document:

```js
import { API, and, or, cond } from 'space-api';

// Initialize api with the project name and url of the space cloud
const api = new API('todo-app', 'http://localhost:8080');

// Initialize database(s) you intend to use
const db = api.Mongo();

// The condition to be matched
const condition = cond('_id', '==', 1)

// Get the todo
db.get('todos').where(condition).one().then(res => {
  if (res.status === 200) {
    // res.data contains the documents returned by the database
    console.log('Response:', res.data);
    return;
  }
}).catch(ex => {
  // Exception occured while processing request
});
```

As you would have noticed, the above function is asynchronous in nature. `todos` is the name of the collection / table which contains all your todos.

You can fetch a single document using `one` method. If the specified document is not present, `res.status` will be 404. Multiple documents can be fetched using `all` method. When you use `all` method, the `res.status` will be 200 even if no matching documents were found.

## Read documents selectively

You can selectively read only a few documents which you desire and filter out the rest by using `where` clause. The `where` method accepts a `condition` object. After validation, `space-cloud` generates a database specific query. The documents or rows which match this query are returned.

### Specifying a single condition

The `cond` function is used to specify a single condition as shown below:

```js
// Syntax
var op1 = 'field_name';
var operator = '== | != | > | < | >= | <= | in | notIn';
var op2 = 'value';
const condition = cond(op1, operator, op2);

// Example
const condition = cond('_id', '==', '1');
db.get('todos').where(condition).one().then(res => ...)
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

A single condition is often not enough to fetch the data you desire. You might need to `and` / `or` several conditions together. For e.g., you want to fetch only those posts which are of a particular author **and** of a particular category. The following code snippet shows how you can use and and or:

```js
// AND syntax
const condition = and(...conditions)

// Example
const condition = and(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.get('todos').where(condition).one().then(res => ...);

// OR syntax
const condition = or(...conditions)

// Example
const condition = or(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.get('todos').where(condition).one().then(res => ...);

```

## Response

A response object sent by the server contains the **status** and **data** fields explained below:

**status** : Number describing the status of the operation. Following values are possible:

- 200 - Operation was successful
- 401 - Request was unauthenticated
- 403 - Request was unauthorized
- 500 - Internal server error

**data** - An object containing the following fields:

- result - A single object in case of `one`, array of objects in case of `all`, an array in case of `distinct` and an integer in case of `count`

## Selecting only a few fields
You can specify which fields to be returned for the docs in the result by using `select` method as shown below:

```js
// Return only title and author field for each post
const selectClause = { title: 1, author: 1 }
db.get('posts').where(cond('category', '==', 'some-category'))
  .select(selectClause).all().then(res => ...);
```

## Sorting, skipping and limiting

Many a times, you will require to receive the data in an sorted order, skip a few docs or limit the result set to a small number or perhaps all the three.

```js
// Sort posts first by title and then by author, skip first 20 posts and fetch only 10 after those 20
db.get('posts').where(cond('category', '==', 'some-category'))
  .sort('title', 'author').skip(20).limit(10).all().then(res => ...);
```

### Sorting
You can receive a sorted result set by using the `sort` function. This is how you do it:
```js
// Sort syntax
db.get(collection).where(conditions).sort(...fields).all().then()

// Example
db.get('posts').where(cond('category', '==', 'some-category'))
  .sort('title', '-author').all().then(res => ...);
```

`sort` takes any number of `fields` as input parameters. `field` is a string corresponding to either field of the JSON document in case of document oriented databases like Mongo DB or name of column in case of SQL databases. The `sort` function sorts the result set in the order of the provided fields. For the above example, result will be sorted by title first and then by author. The minus sign in front of author means that the result set will sorted in a descending order for the author field.

### Skipping
You can skip n number of rows from the beginning of the result set by using `skip`. It takes an integer as an parameter. This is how you can skip docs:

```js
// Skip syntax
db.get(collection).where(conditions).skip(n).all().then()

// Skip 20 rows / docs
db.get('posts').where(cond('category', '==', 'some-category'))
  .skip(20).all().then(res => ...);
```

### Limiting
You can limit the number of docs / rows returned by the server by using `limit`. It takes an integer as an parameter. This is how you can limit result set:
```js
// Limit syntax
db.get(collection).where(conditions).limit(n).all().then()

// Limit result up to 10 rows / docs
db.get('posts').where(cond('category', '==', 'some-category'))
  .limit(10).all().then(res => ...);
``` 

## Reading a single document:

```js
// Reading a single todo
db.get('todos').where(cond('_id', '==', 1)).one().then(res => ...).catch(ex => ...);
```

## Reading multiple documents at once:
```js
// Read multiple todos at once!
db.get('todos').where(cond('categories', '==', 'some-category')).all().then(res => ...).catch(ex => ...);
```

## Reading count of number of documents:
Sometimes, you might only want to fetch the number of documents for a given query but not the actual result or docs. In that case, you can use `count` method. This is how you can fetch just the count:
```js
// Read count of todos
db.get('todos').where(cond('categories', '==', 'some-category')).count().then(res => ...).catch(ex => ...);
```

> Note: `count` is only available in Mongo DB.

## Reading only distinct values:
You can read the distinct values of a given field by using the `distinct` method as shown below:

```js
// Read distinct values of categories of todos
db.get('todos').distinct('category').then(res => ...).catch(ex => ...);
```
`res.data.result` will be an array of the unique values for the given field.

## Next steps

So you know how to read data from a database using Space Cloud. Now let's check how to update it.

<div class="btns-wrapper">
  <a href="/docs/database/create" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/update" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
