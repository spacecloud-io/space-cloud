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
      <li class="tab col s2"><a href="#upload-golang">Golang</a></li>
    </ul>
  </div>
  <div id="upload-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from "space-api";

// Initialize api with the project name and url of the space cloud
const api = new API("todo-app", "http://localhost:4122");

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
API api = new API("books-app", "localhost", 4124);
FileStore fileStore = api.fileStore();
InputStream inputStream = new FileInputStream("input.txt");
fileStore.uploadFile("\\", "file.txt", inputStream, new Utils.ResponseListener() {
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
 <div id="upload-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:4124")

# Initialize file storage module
file_store = api.file_store()

# Upload a file (to be named "new.txt" [remote]) into location ("\\" [remote]) from a file ("a.txt" [local])
response = file_store.upload_file("\\", "new.txt", "a.txt")
if response.status == 200:
    print("Success")
else:
    print(response.error)
      </code>
    </pre>
  </div>
  <div id="upload-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
api, err := api.New("books-app", "localhost:4124", false)
if(err != nil) {
  fmt.Println(err)
}
filestore := api.Filestore()

file, err := os.Open("a.txt")
if err != nil {
  panic(err)
}
resp, err := filestore.UploadFile("\\Folder", "hello1.txt", file)
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
url = `http://localhost:4122/v1/api/$projectName/files/$path/$fileName`
```

The url is different for each file and has following variable parts to it:
- **$projectName** - This is the name of project with which you initialized the API
- **$path** - This is the path at which the file was uploaded
- **$fileName** - This is the name with which the file was uploaded

A file can also be downloaded directly into a stream (or file), using the Java, Python and Golang clients.
Here's a code snippet to download a file:

 <div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a href="#download-java">Java</a></li>
      <li class="tab col s2"><a href="#download-python">Python</a></li>
      <li class="tab col s2"><a href="#download-golang">Golang</a></li>
    </ul>
  </div>
  <div id="download-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
API api = new API("books-app", "localhost", 4124);
FileStore fileStore = api.fileStore();
OutputStream outputStream = new FileOutputStream("output.txt";);
fileStore.downloadFile("\\file.txt", outputStream, new Utils.ResponseListener() {
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
 <div id="download-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

# Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:4124")

# Initialize file storage module
file_store = api.file_store()

# Download the file ("\\a.txt" [remote]) into a file ("b.txt" [local])
response = file_store.download_file("\\a.txt", "b.txt")
if response.status == 200:
    print("Success")
else:
    print(response.error)
      </code>
    </pre>
  </div>
  <div id="download-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
api, err := api.New("books-app", "localhost:4124", false)
if(err != nil) {
  fmt.Println(err)
}
filestore := api.Filestore()

file, err := os.Create("test1.txt")
if err != nil {
  fmt.Println("Error:", err)
  return
}
defer file.Close()
err = filestore.DownloadFile("\\Folder\\text.txt", file)
if err != nil {
  fmt.Println("Error:", err)
}
      </code>
    </pre>
  </div>
</div>


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
