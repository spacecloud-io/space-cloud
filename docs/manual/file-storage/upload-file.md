# Upload and download files with ease!

You can easily allow users to upload and download files via the File Management module of Space Cloud.

## Upload a file

Uploading a file via Space Cloud from frontend is as simple as getting the reference of the file and calling `uploadFile` on the frontend. Here's a code snippet to upload the file:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#upload-js">Javascript</a></li>
      <li class="tab col s2"><a href="#upload-java">Java</a></li>
      <li class="tab col s2"><a href="#upload-python">Python</a></li>
    </ul>
  </div>
  <div id="upload-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Get the file to be uploaded
const myFile = document.querySelector("#your-file-input").files[0];

// Upload the file
api.FileStore()
  .uploadFile("/some-path", myFile, "fileName")
  .then(res => {
    if (res.status === 200) {
      // File uploaded successfully
    }
    // Error uploading file
  })
  .catch(ex => {
    // Exception occured while processing request
  });
      </code>
    </pre>
  </div>
  <div id="upload-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
 <div id="upload-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
# Python client coming soon!
      </code>
    </pre>
  </div>
</div>

The `uploadFile` function takes two parameters to upload a file which are as follows:
- **path** - The path at which to upload the file.
- **file** - A file of the type HTML5 File API.
- **name** - Name of the file.

The `path` can be nested as well. For e.g a `path` - /folder1/folder2 would mean to upload the file inside folder2 which is in folder1. If any of the folders mentioned in the `path` were not present, they would be created before uploading the file.

## Response

A response object sent by the server contains the **status** fields explained below:

**status:** Number describing the status of the upload operation. Following values are possible:

- 200 - Successful upload
- 401 - Request was unauthenticatedoverview
- 403 - Request was unauthorized
- 500 - Internal server error


## Download a file

All files uploaded via File Management module are accessible on the following url:

```
url = `http://localhost:8080/v1/api/$projectName/files/$path/$fileName`
```

The url is different for each file and has following variable parts to it:
- **$projectName** - This is the name of project with which you initialized the API
- **$path** - This is the path at which the file was uploaded
- **$fileName** - This is the name with which the file was uploaded


## Next steps

Now you know how to upload and download a file. So let's check how to create a folder.

<div class="btns-wrapper">
  <a href="/docs/file-storage/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/file-storage/create-folder" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
