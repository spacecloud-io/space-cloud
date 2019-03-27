# List files and folder

You can list all files and folders within a specific folder by simply calling `listFiles` on the front end. Here's a code snippet that shows how to do it: 

```js
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Upload the file
api.FileStore()
  .listFiles("/some-path")
  .then(res => {
    if (res.status === 200) {
      // res.data.result contains list of files / folders
      console.log("Files: ", res.data.result)
    }
    // Error fetching list of files
  })
  .catch(ex => {
    // Exception occured while processing request
  });
```

As shown above, the `listFiles` function takes a single parameter `path` and lists all the files / folders located at that path.

The `path` can be nested as well. For e.g if you give a  `path` - /folder1/folder2, then `listFiles` method will list all the files / folders located in folder2 which is in folder1.

## Response
A response object sent by the server contains the **status** and **data** fields explained below:

**status** : Number describing the status of the upload operation. Following values are possible:

- 200 - Successfully fetched list of files / folders
- 401 - Request was unauthenticated
- 403 - Request was unauthorized
- 404 - Path does not exists
- 500 - Internal server error

**data**: The data is an object which contains `result` which is an array of objects wherein each object contains the following:
- **name** - Name of the file / folder
- **type** - A string whose value is "file" for a file and "dir" for a folder

## Next steps

Now you know how to list files / folders within a folder. So let's check how to delete a file / folder.

<div class="btns-wrapper">
  <a href="/docs/file-storage/create-folder" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/file-storage/delete" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
