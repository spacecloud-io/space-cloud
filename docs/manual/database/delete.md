# Deleting Data

You can delete data in your app by simply calling `db.delete` on the frontend. Here's a code snippet to delete all documents matching a specific condition:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#delete-js">Javascript</a></li>
      <li class="tab col s2"><a href="#delete-java">Java</a></li>
      <li class="tab col s2"><a href="#delete-python">Python</a></li>
      <li class="tab col s2"><a href="#delete-golang">Golang</a></li>
    </ul>
  </div>
  <div id="delete-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API, and, or, cond } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:4122");

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
API api = new API("books-app", "localhost", 4124);
SQL db = api.MySQL();
db.delete("books").where(new Cond("name", "==", "aBook")).apply(new Utils.ResponseListener() {
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
  <div id="delete-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API, COND

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:4124")

# Initialize database(s) you intend to use
db = api.my_sql()

# The condition to be matched
condition = COND("name", "==", "SomeAwesomeBook")

# Delete all books which match a particular condition
response = db.delete("books").where(condition).apply()
if response.status == 200:
    print("Success")
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
  <div id="delete-golang" class="col s12" style="padding:0">
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
	resp, err := db.Delete("books").Where(condition).Apply()
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

As you would have noticed, the `delete` method is asynchronous in nature. It takes the name of the concerned collection/table and deletes all the matching documents. The `apply` method actually triggers the given request to `space-cloud` and returns a promise.

## Deleting a single document:

> **Note:** `deleteOne` method is not available for SQL databases.

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#delete-one-js">Javascript</a></li>
      <li class="tab col s2"><a href="#delete-one-java">Java</a></li>
      <li class="tab col s2"><a href="#delete-one-python">Python</a></li>
      <li class="tab col s2"><a href="#delete-one-golang">Golang</a></li>
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
db.deleteOne("books").where(new Cond("name", "==", "aBook")).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="delete-one-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("name", "==", "SomeAwesomeBook")
response = db.delete_one("books").where(condition).apply()
      </code>
    </pre>
  </div>
  <div id="delete-one-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("author", "==", "anAuthor")
resp, err := db.DeleteOne("books").Where(condition).Apply()
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
db.delete('todos').where(condition).apply().then(res => ...);
      </code>
    </pre>
  </div>
  <div id="cond-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
db.delete("books").where(new Cond("name", "==", "aBook")).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="cond-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = COND("name", "==", "SomeAwesomeBook")
response = db.delete("books").where(condition).apply()
      </code>
    </pre>
  </div>
  <div id="cond-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition := utils.Cond("id", "==", 1)
resp, err := db.Delete("books").Where(condition).Apply()
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
db.delete('todos').where(condition).apply().then(res => ...)

// OR syntax
const condition = or(...conditions)

// Example
const condition = or(cond('author', '==', 'some-author'), cond('category', '==', 'some-category'));
db.delete('todos').where(condition).apply().then(res => ...)
      </code>
    </pre>
  </div>
   <div id="multiple-cond-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
db.delete("books").where(And.create(new Cond("name", "==", "aBook"), new Cond("author", "==", "myelf"))).apply(myResponseListener);
      </code>
    </pre>
  </div>
  <div id="multiple-cond-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
condition = AND(COND("name", "==", "SomeAwesomeBook"), COND("author", "==", "SomeAuthor"))
response = db.delete("books").where(condition).apply()
      </code>
    </pre>
  </div>
  <div id="multiple-cond-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
condition1 := utils.Cond("id", "==", 1)
condition2 := utils.Cond("id", "==", 2)
condition := utils.Or(condition1, condition2)
resp, err := db.Delete("books").Where(condition).Apply()
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

Now you know how to perform all the basic operations of CRUD module. Now let's see how to perform batched operations and transactions.

<div class="btns-wrapper">
  <a href="/docs/database/update" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/transactions" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
