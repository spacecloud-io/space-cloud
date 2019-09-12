# Updating Data

There are 3 ways of updating data in your app:
- [update](/docs/database/update#updating-all-documents) - Updating all documents matching a specific condition
- [updateOne](/docs/database/update#updating-a-single-document) - Updating a single document
- [upsert](/docs/database/update#upserting-a-document) - Update a document if present, else create it   

## <a name="updating-all-documents"></a>Updating all documents

You can update / modify all documents in your app matching a specific condition by simply calling `db.update` on the frontend. Here's a code snippet to do so:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#update-js">Javascript</a></li>
      <li class="tab col s2"><a href="#update-java">Java</a></li>
      <li class="tab col s2"><a href="#update-python">Python</a></li>
      <li class="tab col s2"><a href="#update-golang">Golang</a></li>
    </ul>
  </div>
  <div id="update-js" class="col s12" style="padding:0">
    <pre>
     <code class="javascript"> 
import { API, and, or, cond } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:4122");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The condition to be matched
const condition = cond("author", "==", 'author1');

// Update the todos
db.update("todos")
  .where(condition)
  .set({ text: "Fork Space Cloud on Github right now" }).apply()
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
 <div id="update-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
API api = new API("books-app", "localhost", 4124);
SQL db = api.MySQL();
HashMap<String, Object> set = new HashMap<>();
set.put("name", "Book1");
db.update("books").where(new Cond("id", "==", 1)).set(set).apply(new Utils.ResponseListener() {
    @Override
    public void onResponse(int statusCode, Response response) {
        if (statusCode == 200) {
            System.out.println("Success");
        } else {
            System.out.println(response.getError());
        }
    }

    @Override
    public void onError(Exception e) {
        System.out.println(e.getMessage());
    }
});
      </code>
    </pre>
  </div>
  <div id="update-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API, AND, OR, COND

api = API("books-app", "localhost:4124")
db = api.my_sql()

# The condition to be matched
condition = COND("author", "==", "author1")

# Update the books
response = db.update("books").where(condition).set({"name": "A book"}).apply()

if response.status == 200:
    print("Success")
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
  <div id="update-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
import (
	"github.com/spaceuptech/space-api-go/api"
	"github.com/spaceuptech/space-api-go/api/utils"
	"fmt"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	condition := utils.Cond("id", "==", 1)
	set := map[string]interface{}{"name":"ABook"}
	resp, err := db.Update("books").Where(condition).Set(set).Apply()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if resp.Status == 200 {
			fmt.Println("Success")
		} else {
			fmt.Println("Error Processing Request:", resp.Error)
		}
	}
}
      </code>
    </pre>
  </div>
</div>

As you would have noticed, the `update` method is asynchronous in nature. It takes the name of the concerned collection/table and updates the matching documents. The `apply` method actually triggers the given request to `space-cloud` and returns a promise.

## <a name="updating-a-single-document"></a>Updating a single document:

> **Note:** `updateOne` method is not available for SQL databases. 

`updateOne` finds and updates a single document. It returns an error if no matching document was found.

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#update-one-js">Javascript</a></li>
      <li class="tab col s2"><a href="#update-one-java">Java</a></li>
      <li class="tab col s2"><a href="#update-one-python">Python</a></li>
      <li class="tab col s2"><a href="#update-one-golang">Golang</a></li>
    </ul>
  </div>
  <div id="update-one-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Set text of todo by new value
db.updateOne('todos').where(cond('_id_', '==', 1))
  .set({text: 'Fork Space Cloud on Github'}).apply().then(res => ...);
      </code>
</pre>
  </div>
  <div id="update-one-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> set = new HashMap<>();
set.put("name", "Book1");
db.updateOne("books").where(new Cond("id", "==", 1)).set(set).apply(myResponseListener);
      </code>
    </pre>
  </div>
 <div id="update-one-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update_one("books").where(condition).set({"name": "A book"}).apply()
      </code>
    </pre>
  </div>
  <div id="update-one-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
set := map[string]interface{}{"name":"ABook"}
resp, err := db.UpdateOne("books").Where(condition).Set(set).Apply()
      </code>
    </pre>
  </div>
</div>

## <a name="upserting-a-document"></a>Upserting a document:

Sometimes you might want to create a document or update it if it already exists. You can acheive this by using the `upsert` method like below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#upsert-js">Javascript</a></li>
      <li class="tab col s2"><a href="#upsert-java">Java</a></li>
      <li class="tab col s2"><a href="#upsert-python">Python</a></li>
      <li class="tab col s2"><a href="#upsert-golang">Golang</a></li>
    </ul>
  </div>
  <div id="upsert-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// Create a todo or update it
db.upsert('todos').where(cond('_id', '==', 1))
  .set({text: 'Fork Space Cloud on Github'}).apply().then(res => ...).catch(ex => ...);
      </code>
</pre>
  </div>
  <div id="upsert-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> set = new HashMap<>();
set.put("name", "Book1");
db.upsert("books").where(new Cond("id", "==", 1)).set(set).apply(myResponseListener);
      </code>
    </pre>
  </div>
 <div id="upsert-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.upsert("books").where(condition).set({"name": "A book"}).apply()
      </code>
    </pre>
  </div>
  <div id="upsert-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
set := map[string]interface{}{"name":"ABook"}
resp, err := db.Upsert("books").Where(condition).Set(set).Apply()
      </code>
    </pre>
  </div>
</div>

The above example will update a todo of _id = 1 with the text - 'Fork Space Cloud on Github' if a todo with _id = 1 already exists. Otherwise it will create a new todo - { _id: 1, text: 'Fork Space Cloud on Github' }

## Update documents selectively

You can selectively update only a few documents which you desire and leave the rest by using `where` clause. The `where` method accepts a `condition` object. After validation, `space-cloud` generates a database specific query. The documents or rows which match this query are updated by the update operations described above.

### Specifying a single condition

The `cond` function is used to specify a single condition as shown below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#cond-js">Javascript</a></li>
      <li class="tab col s2"><a href="#cond-java">Java</a></li>
      <li class="tab col s2"><a href="#cond-python">Python</a></li>
      <li class="tab col s2"><a href="#cond-golang">Golang</a></li>
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
db.update('todos').where(condition)
  .set({text: 'Fork Space Cloud on Github'}).apply().then(res => ...); 
      </code>
    </pre>
  </div>
   <div id="cond-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> set = new HashMap<>();
set.put("name", "Book1");
db.update("books").where(new Cond("id", "==", 1)).set(set).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="cond-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).set({"name": "A book"}).apply()
      </code>
    </pre>
  </div>
  <div id="cond-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
set := map[string]interface{}{"name":"ABook"}
resp, err := db.Update("books").Where(condition).Set(set).Apply()
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

A single condition is often not enough to update the data you desire. You might need to `and` / `or` several conditions together. For e.g., you want to update only those posts which are of a particular author **and** of a particular category. The following code snippet shows how you can use and and or:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#multiple-cond-js">Javascript</a></li>
      <li class="tab col s2"><a href="#multiple-cond-java">Java</a></li>
      <li class="tab col s2"><a href="#multiple-cond-python">Python</a></li>
      <li class="tab col s2"><a href="#multiple-cond-golang">Golang</a></li>
    </ul>
  </div>
  <div id="multiple-cond-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
// AND syntax
const condition = and(...conditions)

// Example
const condition = and(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.update('todos').where(condition)
  .set({text: 'Fork Space Cloud on Github'}).apply().then(res => ...)

// OR syntax
const condition = or(...conditions)

// Example
const condition = or(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.update('todos').where(condition)
  .set({text: 'Fork Space Cloud on Github'}).apply().then(res => ...);
      </code>
    </pre>
  </div>
   <div id="multiple-cond-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> set = new HashMap<>();
set.put("name", "Book1");
db.update("books").where(Or.create(new Cond("id", "==", 1), new Cond("name", "==", "aBook"))).set(set).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="multiple-cond-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = AND(COND("author", "==", "author1"), COND("name", "==", "someBook"))
response = db.update("books").where(condition).set({"name": "A book"}).apply()
      </code>
    </pre>
  </div>
  <div id="multiple-cond-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition1 := utils.Cond("id", "==", 1)
condition2 := utils.Cond("id", "==", 2)
condition := utils.Or(condition1, condition2)
set := map[string]interface{}{"name":"ABook"}
resp, err := db.Update("books").Where(condition).Set(set).Apply()
      </code>
    </pre>
  </div>
</div>


## Update operations

You can perform different types of update operations like set, push, rename, etc. on your data. Following are the different types of update operations:

### Set operation

You can set the value of any field / column in your data by using `set` method like below: 
<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#set-js">Javascript</a></li>
      <li class="tab col s2"><a href="#set-java">Java</a></li>
      <li class="tab col s2"><a href="#set-python">Python</a></li>
      <li class="tab col s2"><a href="#set-golang">Golang</a></li>
    </ul>
  </div>
  <div id="set-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Set text of todo by new value
db.update('todos').where(cond('_id_', '==', 1))
  .set({text: 'Fork Space Cloud on Github'}).apply().then(res => ...);    
      </code>
    </pre>
  </div>
 <div id="set-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> set = new HashMap<>();
set.put("name", "Book1");
db.update("books").where(new Cond("id", "==", 1)).set(set).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="set-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).set({"name": "A book"}).apply()
      </code>
    </pre>
  </div>
  <div id="set-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
set := map[string]interface{}{"name":"ABook"}
resp, err := db.Update("books").Where(condition).Set(set).Apply()
      </code>
    </pre>
  </div>
</div>


The `set` method accepts an object of key-value pairs where key is the field / column name whereas value is the new value with which you want to update the field / column. 

### Push operation

You can push an element to an array in a document by using the `push` method like below:

> **Note:** `push` operation is only available in MongoDB.

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#push-js">Javascript</a></li>
      <li class="tab col s2"><a href="#push-java">Java</a></li>
      <li class="tab col s2"><a href="#push-python">Python</a></li>
      <li class="tab col s2"><a href="#push-golang">Golang</a></li>
    </ul>
  </div>
  <div id="push-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Add a new category for a todo
db.update('todos').where(cond('_id_', '==', 1))
  .push({categories: 'some-category'}).apply().then(res => ...);  
      </code>
    </pre>
  </div>
 <div id="push-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> push = new HashMap<>();
push.put("name", "Book1");
db.update("books").where(new Cond("id", "==", 1)).push(push).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="push-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).push({"name": "A book"}).apply()
      </code>
    </pre>
  </div>
  <div id="push-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
push := map[string]interface{}{"name":"ABook"}
resp, err := db.Update("books").Where(condition).Push(push).Apply()
      </code>
    </pre>
  </div>
</div>

The `push` method accepts an object of key-value pairs where key is the field name of the array whereas value is the new value which you want to push into that array.

### Remove operation

You can remove / delete a field in a document by using `remove` method like below:

> **Note:** `remove` operation is only available in MongoDB.

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#remove-js">Javascript</a></li>
      <li class="tab col s2"><a href="#remove-java">Java</a></li>
      <li class="tab col s2"><a href="#remove-python">Python</a></li>
      <li class="tab col s2"><a href="#remove-golang">Golang</a></li>
    </ul>
  </div>
  <div id="remove-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Remove categories and time field
db.update('todos').where(cond('_id_', '==', 1))
  .remove('categories', 'time').apply().then(res => ...);
      </code>
    </pre>
  </div>
 <div id="remove-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
db.update("books").where(new Cond("id", "==", 1)).remove("author").apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="remove-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).remove("author").apply()
      </code>
    </pre>
  </div>
  <div id="remove-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
resp, err := db.Update("books").Where(condition).Remove("author").Apply()
      </code>
    </pre>
  </div>
</div>

The `remove` method accepts multiple inputs each being the name of a field you want to remove.  

### Rename operation

You can rename the name of a field in a document by using `rename` method like below:

> **Note:** `rename` operation is only available in MongoDB.

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#rename-js">Javascript</a></li>
      <li class="tab col s2"><a href="#rename-java">Java</a></li>
      <li class="tab col s2"><a href="#rename-python">Python</a></li>
      <li class="tab col s2"><a href="#rename-golang">Golang</a></li>
    </ul>
  </div>
  <div id="rename-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Rename text field to 'content'
db.update('todos').where(cond('_id_', '==', 1))
  .rename({text: 'content'}).apply().then(res => ...);
      </code>
    </pre>
  </div>
 <div id="rename-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> rename = new HashMap<>();
rename.put("name", "bookName");
db.update("books").where(new Cond("id", "==", 1)).rename(rename).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="rename-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).rename({"writer":"author"}).apply()
      </code>
    </pre>
  </div>
  <div id="rename-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
rename := map[string]interface{}{"writer": "author"}
resp, err := db.Update("books").Where(condition).Rename(rename).Apply()
      </code>
    </pre>
  </div>
</div>

The `rename` method accepts an object of key-value pairs where key is the current name of a field while value (string) is the new name that you want to assign to that field.

### Increment operation

You can increment / decrement the value of a integer field in your data by using the `inc` method like below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#inc-js">Javascript</a></li>
      <li class="tab col s2"><a href="#inc-java">Java</a></li>
      <li class="tab col s2"><a href="#inc-python">Python</a></li>
      <li class="tab col s2"><a href="#inc-golang">Golang</a></li>
    </ul>
  </div>
  <div id="inc-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Increment 'likes' by 3
db.update('todos').where(cond('_id_', '==', 1))
  .inc({likes: 3}).apply().then(res => ...);

// Decrement 'likes' by 2
db.update('todos').where(cond('_id_', '==', 1))
  .inc({likes: -2}).apply().then(res => ...);
      </code>
    </pre>
  </div>
 <div id="inc-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> increment = new HashMap<>();
increment.put("likes", 1);
db.update("books").where(new Cond("id", "==", 1)).inc(increment).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="inc-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).inc({"likes":1}).apply()
      </code>
    </pre>
  </div>
  <div id="inc-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
inc := map[string]interface{}{"likes": 1}
resp, err := db.Update("books").Where(condition).Inc(inc).Apply()
      </code>
    </pre>
  </div>
</div>

The `inc` method accepts an object of key-value pairs where key is the name of the field whose value is to be incremented whereas value is the amount by which the value has to be incremented. As you would have noticed you can decrement a value by using negative integers.

### Multiply operation

You can multiply the value of a integer field in your data by using the `mul` method like below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#mul-js">Javascript</a></li>
      <li class="tab col s2"><a href="#mul-java">Java</a></li>
      <li class="tab col s2"><a href="#mul-python">Python</a></li>
      <li class="tab col s2"><a href="#mul-golang">Golang</a></li>
    </ul>
  </div>
  <div id="mul-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Multiply 'likes' by 10
db.update('todos').where(cond('_id_', '==', 1))
  .mul({likes: 10}).apply().then(res => ...);
      </code>
    </pre>
  </div>
 <div id="mul-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> mul = new HashMap<>();
mul.put("likes", 2);
db.update("books").where(new Cond("id", "==", 1)).mul(mul).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="mul-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).mul({"likes":10}).apply()
      </code>
    </pre>
  </div>
  <div id="mul-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
mul := map[string]interface{}{"likes": 10}
resp, err := db.Update("books").Where(condition).Mul(mul).Apply()
      </code>
    </pre>
  </div>
</div>

The `inc` method accepts an object of key-value pairs where key is the name of the field whose value is to be multiplied whereas value is the amount by which the value has to be multiplied.

### Max operation

Sometimes you might want to update a number in your document with a new value only if the new value is greater than the existing value. This can be acheived by using `max` method like below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#max-js">Javascript</a></li>
      <li class="tab col s2"><a href="#max-java">Java</a></li>
      <li class="tab col s2"><a href="#max-python">Python</a></li>
      <li class="tab col s2"><a href="#max-golang">Golang</a></li>
    </ul>
  </div>
  <div id="max-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Updates 'likes' if it was lesser than 50
db.update('todos').where(cond('_id_', '==', 1))
  .max({likes: 50}).apply().then(res => ...);
      </code>
    </pre>
  </div>
 <div id="max-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> max = new HashMap<>();
max.put("likes", 100);
db.update("books").where(new Cond("id", "==", 1)).max(max).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="max-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).max({"likes":100}).apply()
      </code>
    </pre>
  </div>
  <div id="max-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
max := map[string]interface{}{"likes": 100}
resp, err := db.Update("books").Where(condition).Max(max).Apply()
      </code>
    </pre>
  </div>
</div>

### Min operation

Sometimes you might want to update a number in your document with a new value only if the new value is lesser than the existing value. This can be acheived by using `min` method like below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#min-js">Javascript</a></li>
      <li class="tab col s2"><a href="#min-java">Java</a></li>
      <li class="tab col s2"><a href="#min-python">Python</a></li>
      <li class="tab col s2"><a href="#min-golang">Golang</a></li>
    </ul>
  </div>
  <div id="min-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Updates 'likes' if it was greater than 50
db.update('todos').where(cond('_id_', '==', 1))
  .min({likes: 50}).apply().then(res => ...);
      </code>
    </pre>
  </div>
 <div id="min-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
HashMap<String, Object> min = new HashMap<>();
min.put("likes", 100);
db.update("books").where(new Cond("id", "==", 1)).min(min).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="min-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).min({"likes":100}).apply()
      </code>
    </pre>
  </div>
  <div id="min-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
min := map[string]interface{}{"likes": 100}
resp, err := db.Update("books").Where(condition).Min(min).Apply()
      </code>
    </pre>
  </div>
</div>

### Current timestamp operation

You can update a field with the value of current timestamp by using the `currentTimestamp` method like below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#current-timestamp-js">Javascript</a></li>
      <li class="tab col s2"><a href="#current-timestamp-java">Java</a></li>
      <li class="tab col s2"><a href="#current-timestamp-python">Python</a></li>
      <li class="tab col s2"><a href="#current-timestamp-golang">Golang</a></li>
    </ul>
  </div>
  <div id="current-timestamp-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Update 'lastUpdated' with current timestamp 
db.update('todos').where(cond('_id_', '==', 1))
  .currentTimestamp('lastUpdated').apply().then(res => ...);
      </code>
    </pre>
  </div>
 <div id="current-timestamp-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
db.update("books").where(new Cond("id", "==", 1)).currentTimestamp("last_read").apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="current-timestamp-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).current_timestamp("last_read").apply()
      </code>
    </pre>
  </div>
  <div id="current-timestamp-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
resp, err := db.Update("books").Where(condition).CurrentTimestamp("lastRead").Apply()
      </code>
    </pre>
  </div>
</div>

The `currentTimestamp` method accepts multiple inputs each being the name of the field you want to update with current timestamp.

### Current date operation

You can update a field with the value of current date by using the `currentDate` method like below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#current-date-js">Javascript</a></li>
      <li class="tab col s2"><a href="#current-date-java">Java</a></li>
      <li class="tab col s2"><a href="#current-date-python">Python</a></li>
      <li class="tab col s2"><a href="#current-date-golang">Golang</a></li>
    </ul>
  </div>
  <div id="current-date-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">  
// Update 'lastUpdated' with current date 
db.update('todos').where(cond('_id_', '==', 1))
  .currentDate('lastUpdated').apply().then(res => ...);
      </code>
    </pre>
  </div>
 <div id="current-date-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
db.update("books").where(new Cond("id", "==", 1)).currentDate("last_read").apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="current-date-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("author", "==", "author1")
response = db.update("books").where(condition).current_date("last_read").apply()
      </code>
    </pre>
  </div>
  <div id="current-date-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
resp, err := db.Update("books").Where(condition).CurrentDate("lastRead").Apply()
      </code>
    </pre>
  </div>
</div>

The `currentDate` method accepts multiple inputs each being the name of the field you want to update with current date.



## Response

On response from the server, the callback passed to the `then` method is called with the response object as described below:

```
{
  "status": "number" // Status of the operation
}
```


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
