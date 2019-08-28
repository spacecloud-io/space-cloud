# Creating data

You can add data to your app by simply calling `db.insert` on the frontend. Here's a code snippet to add a single document/record to your app:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#insert-js">Javascript</a></li>
      <li class="tab col s2"><a href="#insert-java">Java</a></li>
      <li class="tab col s2"><a href="#insert-python">Python</a></li>
      <li class="tab col s2"><a href="#insert-golang">Golang</a></li>
    </ul>
  </div>
  <div id="insert-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:4122");

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
API api = new API("books-app", "localhost", 4124);
SQL db = api.MySQL();
Map<String, String> document = new HashMap<>();
document.put("name", "aBook");
db.insert("books").doc(document).apply(new Utils.ResponseListener() {
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
 <div id="insert-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:4124")

# Initialize database(s) you intend to use
db = api.my_sql()

# The book to be created
document = {"name": "SomeAwesomeBook"}

result = db.insert("books").doc(document).apply()
if result.status == 200:
    print("Success")
else:
    print(result.error)

api.close()
      </code>
    </pre>
  </div>
  <div id="insert-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
import (
	"github.com/spaceuptech/space-api-go/api"
	"fmt"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	doc := map[string]interface{}{"name":"SomeBook"}
	resp, err := db.Insert("books").Doc(doc).Apply()
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

As you would have noticed, the `insert` method is asynchronous in nature. It takes the name of the concerned collection/table. The `doc` method takes an object to be inserted. The `apply` method actually triggers the given request to `space-cloud` and returns a promise.

## Adding multiple documents simultaneously:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#insertmany-js">Javascript</a></li>
      <li class="tab col s2"><a href="#insertmany-java">Java</a></li>
      <li class="tab col s2"><a href="#insertmany-python">Python</a></li>
      <li class="tab col s2"><a href="#insertmany-golang">Golang</a></li>
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
HashMap<String, String> document = new HashMap<>();
document.put("name", "aBook");
HashMap<String, String> document2 = new HashMap<>();
document2.put("name", "anotherBook");
HashMap[] docs = new HashMap[2];
docs[0] = document;
docs[1] = document2;
db.insert("books").docs(docs).apply(myResponseListener);
      </code>
    </pre>
  </div>
 <div id="insertmany-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
documents = [{"name": "SomeAwesomeBook"},{"name": "AnotherAwesomeBook"}]
result = db.insert("books").docs(documents).apply()
      </code>
    </pre>
  </div>
  <div id="insertmany-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
docs := make([]map[string]interface{}, 2)
docs[0] = map[string]interface{}{"name": "SomeBook"}
docs[1] = map[string]interface{}{"name": "SomeOtherBook"}
resp, err := db.Insert("books").Docs(docs).Apply()
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
