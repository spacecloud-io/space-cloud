# Listening to real-time`db.liveQuery
You can listen / subscribe to changes happening in your app's data in real time by simply calling `db.liveQuery` on the frontend. Here's a code snippet to do this:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#live-query-js">Javascript</a></li>
      <li class="tab col s2"><a href="#live-query-java">Java</a></li>
      <li class="tab col s2"><a href="#live-query-python">Python</a></li>
      <li class="tab col s2"><a href="#live-query-golang">Golang</a></li>
    </ul>
  </div>
  <div id="live-query-js" class="col s12" style="padding:0">
    <pre>
      <code>
import { API, cond, or, and } from 'space-api';

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The condition to be matched
const condition = cond("category", "==", 'frontend');

// Callback for data changes:
const onSnapshot  = (docs, type) => {
   console.log(docs, snapshot)
}

// Callback for error while subscribing
const onError = (err) => {
   console.log('Live query error', err)
}

// Subscribe to any changes in posts of 'frontend' category
let unsubscribe = db.liveQuery('posts').where(condition).subscribe(onSnapshot, onError) 

// Unsubscribe to changes
if (on some logic) {
  unsubscribe()
}
      </code>
    </pre>
  </div>
  <div id="live-query-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
API api = new API("books-app", "localhost", 8081);
SQL db = api.MySQL();
LiveQueryUnsubscribe unsubscribe = db.liveQuery("books").subscribe(new LiveDataListener() {
    @Override
    public void onSnapshot(LiveData data, String type) {
        System.out.println(type);
        for (Book book : data.getValue(Book.class)) {
            System.out.printf("ID:%d, Name:%s, Author:%s\n", book.getId(), book.getName(), book.getAuthor());
        }
        System.out.println();
    }

    @Override
    public void onError(String error) {
        System.out.println(error);
    }
});

// After some condition
unsubscribe.unsubscribe();
      </code>
    </pre>
  </div>
 <div id="live-query-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

api = API('books-app', 'localhost:8081')
db = api.my_sql()


def on_snapshot(docs, kind):
    print("DOCS:", docs)
    print("KIND OF LIVE QUERY:", kind)


def on_error(error):
    print("ERROR:", error)


unsubscribe = db.live_query('books').subscribe(on_snapshot, on_error)

# After some logic/condition
unsubscribe()
api.close()
      </code>
    </pre>
  </div>
  <div id="live-query-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
import (
	"github.com/spaceuptech/space-api-go/api"
	"github.com/spaceuptech/space-api-go/api/model"
	"fmt"
)

func main() {
	api, err := api.Init("books-app", "localhost", "8081", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	db.LiveQuery("books").Subscribe(func(liveData *model.LiveData, changeType string) () {
		fmt.Println(changeType)
		var v []interface{}
		liveData.Unmarshal(&v)
		fmt.Println(v)
	}, func(err error) () {
		fmt.Println(err)
	})
	for {}
}
      </code>
    </pre>
  </div>
</div>

`liveQuery` function takes the name of the collection / table on which you want to subscribe. `subscribe` takes two functions `onSnapshot` and `onError` functions as it's input paramters and makes a request to subscribe for the given collection / table and `where` clause. 

`onSnapshot` function is called for the first time when you have successfully subscribed with the initial data and on consequent data changes (i.e. whenever new data is added, removed or updated within the where clause). The `onSnapshot` function is always called with the following two params: 
**docs:** An array of latest result set.
**type:** Type of operation due to which the `onSnapshot` is called. It can have one of the following values:
- **initial** - Called only once for the initial data on successful subscription
- **write** - Whenever any data is added or updated
- **delete** - Whenever any data is deleted

`onError` function is called with the `error` if there was any error subscribing to data.

As you would have noticed the `subscribe` function returns an `unsubscribe` function. You should call this function whenever you want to unsubscribe to the changes.

## Next steps

Now you know how to subscribe to realtime changes in data. The next step would be to have a look at how to update data from your frontend.

<div class="btns-wrapper">
  <a href="/docs/realtime/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/update" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
