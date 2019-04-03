# Listening to real-time`db.liveQuery
You can listen / subscribe to changes happening in your app's data in real time by simply calling `db.liveQuery` on the frontend. Here's a code snippet to do this:

```js
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
```

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

Now you know the basics of realtime module. The next step would be to have a look at the file storage module so that you can upload and download files directly from your frontend.

<div class="btns-wrapper">
  <a href="/docs/realtime/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/file-storage/overview" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
