# Reading Data

There are 4 ways of reading data in your app:
- [get](/docs/database/read#reading-all-documents) - Reading all documents matching a specific condition
- [getOne](/docs/database/read#reading-a-single-document) - Reading a single document
- [distinct](/docs/database/read#reading-only-distinct-values) - Reading the unique values of a field  
- [count](/docs/database/read#reading-count-of-number-of-documents) - Reading the count of documents matching a specific condition

## Reading all documents

You can query all documents from your database that matches a particular condition by simply calling `db.get` on the frontend. Here's a code snippet to do so:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#get-js">Javascript</a></li>
      <li class="tab col s2"><a href="#get-java">Java</a></li>
      <li class="tab col s2"><a href="#get-python">Python</a></li>
    </ul>
  </div>
  <div id="get-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API, and, or, cond } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The condition to be matched
const condition = cond("_id", "==", 1);

// Get the todos
db.get("todos").where(condition).apply().then(res => {
    if (res.status === 200) {
      // res.data.result contains the documents returned by the database
      console.log("Response:", res.data.result);
      return;
    }
  })
  .catch(ex => {
    // Exception occured while processing request
  });
      </code>
    </pre>
  </div>
 <div id="get-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!    
      </code>
    </pre>
  </div>
  <div id="get-python" class="col s12" style="padding:0">
    <pre>
     <code class="python">
from space_api import API, AND, OR, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = COND("id", "==", "1")

# Get the books
response = db.get("books").where(condition).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

As you would have noticed, the `get` method is asynchronous in nature. It takes the name of the concerned collection/table. The `apply` method actually triggers the given request to `space-cloud` and returns a promise where `res.data.result` is the array of documents fetched.
  

## <a name="reading-a-single-document"></a>Reading a single document:  

You can fetch a single document from your database that matches a particular condition by simply calling `db.getOne` on the frontend. Here's a code snippet to do so:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#get-one-js">Javascript</a></li>
      <li class="tab col s2"><a href="#get-one-java">Java</a></li>
      <li class="tab col s2"><a href="#get-one-python">Python</a></li>
    </ul>
  </div>
  <div id="get-one-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Reading a single todo
db.getOne('todos').where(cond('_id', '==', 1)).apply().then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="get-one-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!    
      </code>
    </pre>
  </div>
  <div id="get-one-python" class="col s12" style="padding:0">
    <pre>
     <code class="python">
from space_api import API, AND, OR, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = COND("author", "==", "SomeAuthor")

# Get the book
response = db.get_one("books").where(condition).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

As you would have noticed, the `getOne` method is asynchronous in nature. It takes the name of the concerned collection/table. The `getOne` method either returns a matching document in `res.data.result` or returns an error (`res.status` - 400). The `apply` method actually triggers the given request to `space-cloud` and returns a promise where `res.data.result` is the required document (object).

## <a name="reading-only-distinct-values"></a>Reading only distinct values:

You can read the distinct values of a given field by simply calling `db.distinct` on the frontend. Here's a code snippet to do so:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#distinct-js">Javascript</a></li>
      <li class="tab col s2"><a href="#distinct-java">Java</a></li>
      <li class="tab col s2"><a href="#distinct-python">Python</a></li>
    </ul>
  </div>
  <div id="distinct-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Read distinct values of categories of todos
db.distinct('todos').distinctKey('category').apply().then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="distinct-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!    
      </code>
    </pre>
  </div>
  <div id="distinct-python" class="col s12" style="padding:0">
    <pre>
     <code class="python">
from space_api import API, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.mongo()

# The condition to be matched
condition = COND("author", "==", "SomeAuthor")

# Get the books
response = db.distinct("books").where(condition).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

As you would have noticed, the `distinct` method is asynchronous in nature. It takes the name of the concerned collection/table. The `distinctKey` method takes the name of the `key` (field/column) for which you want to find unqiue values. The `apply` method actually triggers the given request to `space-cloud` and returns a promise where `res.data.result` is the array of the unique values for the given `key`.

## <a name="reading-count-of-number-of-documents"></a>Reading count of number of documents:

Sometimes, you might only want to fetch the number of documents for a given query but not the actual result or docs. In that case, you can use `db.count` method. Here's a code snippet to do so:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#count-js">Javascript</a></li>
      <li class="tab col s2"><a href="#count-java">Java</a></li>
      <li class="tab col s2"><a href="#count-python">Python</a></li>
    </ul>
  </div>
  <div id="count-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Read count of todos
db.count('todos').where(cond('categories', '==', 'some-category')).apply().then(res => ...).catch(ex => ...);
      </code>
    </pre>
  </div>
  <div id="count-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!    
      </code>
    </pre>
  </div>
  <div id="count-python" class="col s12" style="padding:0">
    <pre>
     <code class="python">
from space_api import API, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.mongo()

# The condition to be matched
condition = COND("author", "==", "SomeAuthor")

# Get the books
response = db.count("books").where(condition).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

> Note: `count` is only available in Mongo DB.

As you would have noticed, the `count` method is asynchronous in nature. It takes the name of the concerned collection/table. The `apply` method actually triggers the given request to `space-cloud` and returns a promise where `res.data.result` is an integer specifying the number of documents matching the given condition.

## Read documents selectively

You can selectively read only a few documents which you desire and filter out the rest by using `where` clause. The `where` method accepts a `condition` object. After validation, `space-cloud` generates a database specific query. The documents or rows which match this query are returned.

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
const condition = cond('_id', '==', '1');
db.get('todos').where(condition).apply().then(res => ...)
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
from space_api import API, AND, OR, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = COND("id", "==", "1")

# Get the books
response = db.get("books").where(condition).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
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

A single condition is often not enough to fetch the data you desire. You might need to `and` / `or` several conditions together. For e.g., you want to fetch only those posts which are of a particular author **and** of a particular category. The following code snippet shows how you can use `and` and `or`:

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
db.get('todos').where(condition).apply().then(res => ...);

// OR syntax
const condition = or(...conditions)

// Example
const condition = or(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.get('todos').where(condition).apply().then(res => ...);
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
from space_api import API, AND, OR, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = AND(COND("id", "==", "1"), COND("author", "==", "SomeAuthor"))

# Get the books
response = db.get("books").where(condition).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

## Dicing data as per your needs

There are variety of ways in which you can dice and slice the data you read, as per your requirements. For example you might want to read only a few selected fields of a document and not the entire document. Or you might want to limit the number of documents fetched. Here are the various ways of dicing your data:

### Selecting only a few fields

You can specify which fields to be returned for the docs in the result by using `select` method as shown below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#select-js">Javascript</a></li>
      <li class="tab col s2"><a href="#select-java">Java</a></li>
      <li class="tab col s2"><a href="#select-python">Python</a></li>
    </ul>
  </div>
  <div id="select-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Return only title and author field for each post
const selectClause = { title: 1, author: 1 }
db.get('posts').where(cond('category', '==', 'some-category'))
  .select(selectClause).apply().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="select-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
  <div id="select-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API, AND, OR, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = COND("author", "==", "SomeAuthor")

# Get the books
response = db.get("books").where(condition).select({"name":1}).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

### Sorting

You can receive a sorted result set by using the `sort` function. This is how you do it:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#sort-js">Javascript</a></li>
      <li class="tab col s2"><a href="#sort-java">Java</a></li>
      <li class="tab col s2"><a href="#sort-python">Python</a></li>
    </ul>
  </div>
  <div id="sort-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Sort syntax
db.get(collection).where(conditions).sort(...fields).apply().then()

// Example
db.get('posts').where(cond('category', '==', 'some-category'))
  .sort('title', '-author').apply().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="sort-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
  <div id="sort-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API, AND, OR, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = COND("author", "==", "SomeAuthor")

# Get the books
response = db.get("books").where(condition).sort("name", "-id").apply()
# "name" -> sort by name, ascending order
# "-name" -> sort by name, descending order
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

`sort` takes any number of `fields` as input parameters. `field` is a string corresponding to either field of the JSON document in case of document oriented databases like Mongo DB or name of column in case of SQL databases. The `sort` function sorts the result set in the order of the provided fields. For the above example, result will be sorted by title first and then by author. The minus sign in front of author means that the result set will sorted in a descending order for the author field.

### Skipping

You can skip n number of rows from the beginning of the result set by using `skip`. It takes an integer as an parameter. This is how you can skip docs:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#skip-js">Javascript</a></li>
      <li class="tab col s2"><a href="#skip-java">Java</a></li>
      <li class="tab col s2"><a href="#skip-python">Python</a></li>
    </ul>
  </div>
  <div id="skip-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Skip syntax
db.get(collection).where(conditions).skip(n).apply().then()

// Skip 20 rows / docs
db.get('posts').where(cond('category', '==', 'some-category'))
  .skip(20).apply().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="skip-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
  <div id="skip-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API, AND, OR, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = COND("author", "==", "SomeAuthor")

# Get the books
response = db.get("books").where(condition).skip(1).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

### Limiting

You can limit the number of docs / rows returned by the server by using `limit`. It takes an integer as an parameter. This is how you can limit result set:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#limit-js">Javascript</a></li>
      <li class="tab col s2"><a href="#limit-java">Java</a></li>
      <li class="tab col s2"><a href="#limit-python">Python</a></li>
    </ul>
  </div>
  <div id="limit-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Limit syntax
db.get(collection).where(conditions).limit(n).apply().then()

// Limit result up to 10 rows / docs
db.get('posts').where(cond('category', '==', 'some-category'))
  .limit(10).apply().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="limit-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
  <div id="limit-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API, AND, OR, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = COND("author", "==", "SomeAuthor")

# Get the books
response = db.get("books").where(condition).limit(2).apply()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

## Response

On response from the server, the callback passed to the `then` method is called with the response object as described below:

```
{
  "status": "number", // Status of the operation
  "data": {
    "result": "object | array | integer" // Result of the get operation
  }
}
```

The type of `data.result` depends on the operation. Its an array of objects for `get`, an object for `getOne`, an array for `distinct` and an integer for `count` operation.

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
