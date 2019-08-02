# Transactions and batched mutations

The crud module supports atomic CRUD operations. In a set of atomic operations, either all of the operations succeed, or none of them are applied. There are two types of atomic operations in the crud module of Space Cloud:

- Transactions (coming soon): A set of read, insert, update and delete operations on snapshot of a database.
- Batched Mutations: A set of inserts, updates and deletes on one or more documents.

## Batched mutations
A batch operation is used to batch multiple insert, update and delete operations in a single request. Either all of them succeed or none of them. You should use it when you don't want to read data between any of the mutations in the batch.

Here's a code snippet to batch multiple mutations in your app:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#batch-js">Javascript</a></li>
      <li class="tab col s2"><a href="#batch-java">Java</a></li>
      <li class="tab col s2"><a href="#batch-python">Python</a></li>
      <li class="tab col s2"><a href="#batch-golang">Golang</a></li>
    </ul>
  </div>
  <div id="batch-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:4122");

// Initialize database(s) you intend to use
const db = api.Mongo();

// Start a batch operation
const batch = db.beginBatch()

// Add operations to be batched
batch.add(db.insert('todos').doc({ _id: 1, text: "Star Space Cloud on Github!" }))
batch.add(db.update('todos').where(cond('id', ==, '1')).set({ text: "Fork Space Cloud on Github!" }))
batch.add(db.delete('some-other-collection').where(cond('id', ==, '1')))

// Trigger the batch request
batch.apply()
  .then(res => {
    if (res.status === 200) {
      // Batch operation was successful

    }
  })
  .catch(ex => {
    // Exception occured while processing request
  });
    </code>
</pre>
  </div>
  <div id="batch-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
API api = new API("books-app", "localhost", 4124);
SQL db = api.MySQL();
Batch batch = db.beginBatch();
Map<String, String> document = new HashMap<>();
document.put("name", "aBook");
batch.add(db.insert("books").doc(document));
HashMap<String, Object> set = new HashMap<>();
set.put("name", "Book1");
batch.add(db.update("books").where(new Cond("id", "==", 1)).set(set));
batch.add(db.delete("books").where(new Cond("id", "==", 1)));
batch.apply(new Utils.ResponseListener() {
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
 <div id="batch-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API, COND

api = API('grpc', 'localhost:4124')
db = api.my_sql()

b = db.begin_batch()
b.add(db.insert('books').doc({"name": "MyBook", "author": "John Doe"}))
b.add(db.insert('books').docs([{"name": "BookName"}, {"name": "BookName"}]))
b.add(db.delete('books').where(COND('name', '!=', 'Book_name')))
response = b.apply()
if response.status == 200:
  print("Success")
else:
  print(response.error)

api.close()
      </code>
    </pre>
  </div>
  <div id="batch-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
from space_api import API, COND

api = API('grpc', 'localhost:4124')
db = api.my_sql()

b = db.begin_batch()
b.add(db.insert('books').doc({"name": "MyBook", "author": "John Doe"}))
b.add(db.insert('books').docs([{"name": "BookName"}, {"name": "BookName"}]))
b.add(db.delete('books').where(COND('name', '!=', 'Book_name')))
response = b.apply()
if response.status == 200:
  print("Success")
else:
  print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

The `add` method is used to add multiple db operations to a batch request. It can be called n number of times to add n number of operations to the batch. As you would have noticed, the batch is triggered by calling `batch.apply` after adding all the mutations. 

> **Note:** `apply` is not to be used on individual operation in a batch request.

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
  <a href="/docs/database/delete" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/config" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>

