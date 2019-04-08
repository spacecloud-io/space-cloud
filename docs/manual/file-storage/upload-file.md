# Upload and download files with ease!

You can easily allow users to upload and download files via the File Management module of Space Cloud.

## Upload a file

Uploading a file via Space Cloud from frontend is as simple as getting the reference of the file and calling `uploadFile` on the frontend. Here's a code snippet to upload the file:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js1">Javascript</a></li>
      <li class="tab col s2"><a href="#java1">Java</a></li>
      <li class="tab col s2"><a href="#python1">Python</a></li>
    </ul>
  </div>
  <div id="js1" class="col s12" style="padding:0">
    <pre>
      <code>
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:8080");

// Get the file to be uploaded
const myFile = document.querySelector("#your-file-input").files[0];

// Upload the file
api.FileStore()
  .uploadFile("/some-path", myFile)
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
  <div id="java1" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python1" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The `uploadFile` function takes two parameters to upload a file which are as follows:
- **path** - The path at which to upload the file.
- **file** - A file of the type HTML5 File API.

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

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#js2">Javascript</a></li>
      <li class="tab col s2"><a href="#java2">Java</a></li>
      <li class="tab col s2"><a href="#python2">Python</a></li>
    </ul>
  </div>
  <div id="js2" class="col s12" style="padding:0">
    <pre>
      <code>
const url = `http://localhost:8080/api/$projectName/files/$path`
      </code>
    </pre>
  </div>
  <div id="java2" class="col s12" style="padding:0">Java Client Coming Soon!</div>
  <div id="python2" class="col s12" style="padding:0">Python Client Coming Soon!</div>
</div>

The url is different for each file and has following variable parts to it:
- **$projectName** - This is the name of project with which you initialized the API
- **$path** - This is the path at which the file was uploaded


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
