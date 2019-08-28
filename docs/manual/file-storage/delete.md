# Delete a file or folder
You can easily allow users to delete a file or folder via the File Management module of Space Cloud by calling a simple function as shown below:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#delete-js">Javascript</a></li>
      <li class="tab col s2"><a href="#delete-java">Java</a></li>
      <li class="tab col s2"><a href="#delete-python">Python</a></li>
      <li class="tab col s2"><a href="#delete-golang">Golang</a></li>
    </ul>
  </div>
  <div id="delete-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:4122");

// Delete a file
api.FileStore()
  .delete("/some-path/some-file")
  .then(res => {
    if (res.status === 200) {
      // File deleted successfully
    }
    // Error deleting file
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
FileStore fileStore = api.fileStore();

fileStore.deleteFile("\\aNewFolder", new Utils.ResponseListener() {
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
from space_api import API

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:4124")

# Initialize file storage module
file_store = api.file_store()

# Delete a file ("\\a.txt" [remote])
response = file_store.delete_file("\\a.txt")
if response.status == 200:
    print("Success")
else:
    print(response.error)
      </code>
    </pre>
  </div>
  <div id="delete-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
api, err := api.New("books-app", "localhost:4124", false)
if(err != nil) {
  fmt.Println(err)
}
filestore := api.Filestore()
resp, err := filestore.DeleteFile("\\app.txt")
if err != nil {
  fmt.Println("Error:", err)
} else {
  if resp.Status == 200 {
    fmt.Println("Success")
  } else {
    fmt.Println("Error Processing Request:", resp.Error)
  }
}
      </code>
    </pre>
  </div>
</div>

The `delete` function takes a `path` of the file or folder to be deleted. The `path` consists of the path at which the file / folder was uploaded / created and the name of the file / folder as well. The `delete` function will recursively delete all files / folders in a folder if the `path` was for a folder. The `path` can be nested as well. For e.g a `path` - /folder1/folder2/file1 would mean to delete the file inside folder2 which is in folder1.

## Response

A response object sent by the server contains the **status** fields explained below:

**status:** Number describing the status of the upload operation. Following values are possible:

- 200 - Successfully deleted file / folder
- 401 - Request was unauthenticated
- 403 - Request was unauthorized
- 500 - Internal server error

## Next steps

Now you know all the operations of file storage module. So let's take a deeper dive into configuring the file storage module of Space Cloud.

<div class="btns-wrapper">
  <a href="/docs/file-storage/list-files" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/file-storage/config" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
