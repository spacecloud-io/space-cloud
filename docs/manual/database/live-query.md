# Listening to real-time `db.liveQuery`
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
const api = new API("todo-app", "http://localhost:4122");

// Initialize database(s) you intend to use
const db = api.Mongo();

// The condition to be matched
const condition = cond("category", "==", 'frontend');

// Callback for data changes:
const onSnapshot  = (docs, type, changedDoc) => {
   console.log(docs, snapshot, changedDoc)
}

// Callback for error while subscribing
const onError = (err) => {
   console.log('Live query error', err)
}

// Subscribe to any changes in posts of 'frontend' category
let subscription = db.liveQuery('posts').where(condition).subscribe(onSnapshot, onError) 

// Unsubscribe to changes
if (on some logic) {
  subscription.unsubscribe()
}
      </code>
    </pre>
  </div>
  <div id="live-query-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
API api = new API("books-app", "localhost", 4124);
SQL db = api.MySQL();
LiveQuerySubscription subscription = db.liveQuery("books").subscribe(new LiveDataListener() {
    @Override
    public void onSnapshot(LiveData data, String type, ChangedData changedData) {
        System.out.println(type);
        for (Book book : data.getValue(Book.class)) {
            System.out.printf("ID:%d, Name:%s, Author:%s\n", book.getId(), book.getName(), book.getAuthor());
        }
        Book book = changedData.getValue(Book.class);
        if (book!=null) {
            System.out.println("CHANGED: ");
            System.out.printf("ID:%d, Name:%s, Author:%s\n", book.getId(), book.getName(), book.getAuthor());
            System.out.println();
        }
        System.out.println();
    }
    @Override
    public void onError(String error) {
        System.out.println(error);
    }
});

// After some condition
subscription.unsubscribe();
      </code>
    </pre>
  </div>
 <div id="live-query-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

api = API('books-app', 'localhost:4124')
db = api.my_sql()


def on_snapshot(docs, kind, changedDoc):
    print("DOCS:", docs)
    print("KIND OF LIVE QUERY:", kind)
    print("CHANGED DOC:", changedDoc)


def on_error(error):
    print("ERROR:", error)


subscription = db.live_query('books').subscribe(on_snapshot, on_error)

# After some logic/condition
subscription.unsubscribe()
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
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	subscription := db.LiveQuery("books").Subscribe(func(liveData *model.LiveData, changeType string, changedData *model.ChangedData) () {
		fmt.Println("type", changeType)
		var v []interface{}
		liveData.Unmarshal(&v)
		fmt.Println("data", v)
		var v2 interface{}
		changedData.Unmarshal(&v2)
		fmt.Println("chagned", v2)
		fmt.Println()
	}, func(err error) () {
		fmt.Println(err)
	})

  // On some condition
  subscription.unsubscribe()
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
- **insert** - Whenever any data is added
- **update** - Whenever any data is updated
- **delete** - Whenever any data is deleted
**changedDoc:** The doc that changed.

`onError` function is called with the `error` if there was any error subscribing to data.

As you would have noticed the `subscribe` function returns an `unsubscribe` function. You should call this function whenever you want to unsubscribe to the changes.

## Setting the liveQuery options:
You can set the liveQuery options using the `options()` function.  
The function helps to set `changesOnly` to true or false (default).  
If `changesOnly` is false, it caches the docs. `onSnapshot` will be called with all 3 parameters set.  
If `changesOnly` is true, it does not cache the docs and also ignores the initial values. `onSnapshot` will be called with only the last 2 parameters set.  

Here's a code snippet to do this:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#live-query-options-js">Javascript</a></li>
      <li class="tab col s2"><a href="#live-query-options-java">Java</a></li>
      <li class="tab col s2"><a href="#live-query-options-python">Python</a></li>
      <li class="tab col s2"><a href="#live-query-options-golang">Golang</a></li>
    </ul>
  </div>
  <div id="live-query-options-js" class="col s12" style="padding:0">
    <pre>
      <code>
let subscription = db.liveQuery('posts').where({}).options({ changesOnly: true }).subscribe(onSnapshot, onError) 
      </code>
    </pre>
  </div>
  <div id="live-query-options-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
LiveQuerySubscription subscription = db.liveQuery("books")
    .options(LiveQueryOptions.Builder().setChangesOnly(true)).subscribe(new LiveDataListener() {
    @Override
    public void onSnapshot(LiveData data, String type, ChangedData changedData) {
        // ...
    }
    @Override
    public void onError(String error) {
        // ...
    }
});
      </code>
    </pre>
  </div>
 <div id="live-query-options-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
subscription = db.live_query('books').options(changes_only=True).subscribe(on_snapshot, on_error)
      </code>
    </pre>
  </div>
  <div id="live-query-options-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
subscription := db.LiveQuery("books").Options(&model.LiveQueryOptions{ChangesOnly: false}).
  Subscribe(func(liveData *model.LiveData, changeType string, changedData *model.ChangedData) () {
		// ...
	}, func(err error) () {
		// ...
	})
}
      </code>
    </pre>
  </div>
</div>


## Getting the current snapshot:
You can also get the current snapshot(temporary variable, not automatically updated when new changes come in), using the `subscription.getSnapshot()` function.  
The snapshot is empty if `changesOnly` is set to true.  

Here's a code snippet to do this:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#live-query-snapshot-js">Javascript</a></li>
      <li class="tab col s2"><a href="#live-query-snapshot-java">Java</a></li>
      <li class="tab col s2"><a href="#live-query-snapshot-python">Python</a></li>
      <li class="tab col s2"><a href="#live-query-snapshot-golang">Golang</a></li>
    </ul>
  </div>
  <div id="live-query-snapshot-js" class="col s12" style="padding:0">
    <pre>
      <code>
let subscription = db.liveQuery('posts').where({}).subscribe(onSnapshot, onError) 
snapshot = subscription.getSnapshot()
      </code>
    </pre>
  </div>
  <div id="live-query-snapshot-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
LiveQuerySubscription subscription = db.liveQuery("books")
    .options(LiveQueryOptions.Builder().setChangesOnly(true)).subscribe(new LiveDataListener() {
    @Override
    public void onSnapshot(LiveData data, String type, ChangedData changedData) {
        // ...
    }
    @Override
    public void onError(String error) {
        // ...
    }
});
LiveData snapshot = subscription.getSnapshot();
// This is a just temporary object, and will not be automatically updated when new changes come in.
      </code>
    </pre>
  </div>
 <div id="live-query-snapshot-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
subscription = db.live_query('books').options(changes_only=True).subscribe(on_snapshot, on_error)
snapshot = subscription.get_snapshot()
      </code>
    </pre>
  </div>
  <div id="live-query-snapshot-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
subscription := db.LiveQuery("books").Options(&model.LiveQueryOptions{ChangesOnly: false}).
  Subscribe(func(liveData *model.LiveData, changeType string, changedData *model.ChangedData) () {
		// ...
	}, func(err error) () {
		// ...
	})
}
var snapshot []interface{}
subscription.GetSnapshot().Unmarshal(&snapshot)
      </code>
    </pre>
  </div>
</div>


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
