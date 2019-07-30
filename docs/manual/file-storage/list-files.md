# List files and folder

You can list all files and folders within a specific folder by simply calling `listFiles` on the frontend. Here's a code snippet that shows how to do it:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#list-js">Javascript</a></li>
      <li class="tab col s2"><a href="#list-java">Java</a></li>
      <li class="tab col s2"><a href="#list-python">Python</a></li>
      <li class="tab col s2"><a href="#list-golang">Golang</a></li>
    </ul>
  </div>
  <div id="list-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:4122");

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
      </code>
    </pre>
  </div>
  <div id="list-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
API api = new API("books-app", "localhost", 4124);
FileStore fileStore = api.fileStore();

fileStore.listFiles("\\", new Utils.ResponseListener() {
    @Override
    public void onResponse(int statusCode, Response response) {
        System.out.println(statusCode);
        if (statusCode == 200) {
            try {
                Map[] files = response.getResults(Map[].class);
                for (Map file : files) {
                    System.out.println("Name: " + file.get("name"));
                    System.out.println("Type: " + file.get("type"));
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
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
 <div id="list-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:4124")

# Initialize file storage module
file_store = api.file_store()

# List all the files in a particular location ("\\" [remote])
response = file_store.list_files("\\")
if response.status == 200:
    print(response.result)
else:
    print(response.error)
      </code>
    </pre>
  </div>
  <div id="list-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
api, err := api.New("books-app", "localhost:4124", false)
if(err != nil) {
  fmt.Println(err)
}
filestore := api.Filestore()
resp, err := filestore.ListFiles("\\")
if err != nil {
  fmt.Println("Error:", err)
} else {
  if resp.Status == 200 {
    var v []map[string]interface{}
    err:= resp.Unmarshal(&v)
    if err != nil {
      fmt.Println("Error Unmarshalling:", err)
    } else {
      fmt.Println("Result:", v)
    }
  } else {
    fmt.Println("Error Processing Request:", resp.Error)
  }
}
      </code>
    </pre>
  </div>
</div>

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
