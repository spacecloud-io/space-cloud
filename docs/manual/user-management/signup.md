# Add basic email sign up to your app 
You can easily allow users to create a new user on your app via email using the `db.signUp` function. Here's a code snippet to do a basic email sign up:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#signup-js">Javascript</a></li>
      <li class="tab col s2"><a href="#signup-java">Java</a></li>
      <li class="tab col s2"><a href="#signup-python">Python</a></li>
      <li class="tab col s2"><a href="#signup-golang">Golang</a></li>
    </ul>
  </div>
  <div id="signup-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from 'space-api';

// Initialize api with the project name and url of the space cloud
const api = new API('todo-app', 'http://localhost:4122');

// Initialize database(s) you intend to use
const db = api.Mongo();

// SignUp
db.signUp('demo@example.com', 'User1', '1234', 'default').then(res => {
  if (res.status === 200) {
    // Set the token id to enable operations of other modules
    api.setToken(res.data.token)
    
    // res.data contains request payload
    console.log('Response:', res.data);
    return;
  }
  // Request failed
}).catch(ex => {
  // Exception occured while processing request
});
      </code>
    </pre>
  </div>
  <div id="signup-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
API api = new API("books-app", "localhost", 4124);
SQL db = api.MySQL();
db.signUp("email", "name", "password", "role", new Utils.ResponseListener() {
    @Override
    public void onResponse(int statusCode, Response response) {
        if (statusCode == 200) {
            try {
                System.out.println("Token: " + response.getResult(Map.class).get("token"));
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
 <div id="signup-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

// Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:4124")

// Initialize database(s) you intend to use
db = api.my_sql()

// Sign Up
response = db.sign_up("user_email", "user_name", "user_password", "user_role")
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
  <div id="signup-golang" class="col s12" style="padding:0">
    <pre>
      <code class="golang">
import (
	"github.com/spaceuptech/space-api-go/api"
	"fmt"
)

func main() {
	api, err := api.New("books-app", "localhost:4124", false)
	if(err != nil) {
		fmt.Println(err)
	}
	db := api.MySQL()
	resp, err := db.SignUp("email", "name", "password", "role")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if resp.Status == 200 {
			var v map[string]interface{}
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
}
      </code>
    </pre>
  </div>
</div>

As you would have noticed, the above function is asynchronous in nature. The `signUp` method takes 4 parameters and creates a new `user` with an auto generated unique id in the `users` collection / table. The 4 parameters used to create a new `user` are as follows:

- **email** - Email of the user (Used to log in)
- **name** - Name of the user
- **pass** - Password of the user (Used to log in)
- **role** - Role of the user (Comes handy in authorization to restrict a feature to specific set of users)

## Response

On getting the sign up request, `space-cloud` validates whether such an user exists already and then creates a new user. A response object sent by the server contains the **status** and **data** fields explained below.

**status:** Number describing the status of the operation. Following values are possible:
- 200 - Operation was successful
- 400 - User already exists
- 401 - Request was unauthenticated
- 403 - Request was unauthorized
- 500 - Internal server error

**data:** The data object consists of the following fields:
- **token** - The JWT token used for authentication and authorization
- **user** - Object / document of the created user

## Next steps

The next step would be fetching the profile of an user(s).

<div class="btns-wrapper">
  <a href="/docs/user-management/signin" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/user-management/profiles" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
