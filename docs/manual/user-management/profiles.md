# Reading user profiles
You can easily read the profiles of any user in your app by simply calling the `db.profiles` (to fetch all profiles) or `db.profile` (to fetch a single profile) functions on the frontend.

## Fetch profile of a single user
You can read the profile of a single user using `db.profile` function. It takes a single parameter - `id` (unique id of the user).

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#profile-js">Javascript</a></li>
      <li class="tab col s2"><a href="#profile-java">Java</a></li>
      <li class="tab col s2"><a href="#profile-python">Python</a></li>
    </ul>
  </div>
  <div id="profile-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from 'space-api';

// Initialize api with the project name and url of the space cloud
const api = new API('todo-app', 'http://localhost:8080');

// Initialize database(s) you intend to use
const db = api.Mongo();

// Read profile of an user
const userId = 'some-user-id'
db.profile(userId).then(res => {
  if (res.status === 200) {
    // res.data.user contains the profile of the user
    console.log('User profile', res.data.user)
    return;
  }
}).catch(ex => {
  // Exception occured while processing request
});
      </code>
    </pre>
  </div>
  <div id="profile-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
 <div id="profile-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

// Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

// Initialize database(s) you intend to use
db = api.my_sql()

// Read profile of an user
response = db.profile("user_id")
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

## Fetch profiles of all users

You can read the profiles of all users with the help of `profiles` function as shown below:

<div class="row tabs-wrapper">
  <div class="col s12" style="padding:0">
    <ul class="tabs">
      <li class="tab col s2"><a class="active" href="#profiles-js">Javascript</a></li>
      <li class="tab col s2"><a href="#profiles-java">Java</a></li>
      <li class="tab col s2"><a href="#profiles-python">Python</a></li>
    </ul>
  </div>
  <div id="profiles-js" class="col s12" style="padding:0">
    <pre>
      <code class="javascript">
import { API } from 'space-api';

// Initialize api with the project name and url of the space cloud
const api = new API('todo-app', 'http://localhost:8080');

// Initialize database(s) you intend to use
const db = api.Mongo();

// Read profiles of all users
db.profiles().then(res => {
  if (res.status === 200) {
    // res.data.users contains the profile of the users
    console.log('Profiles', res.data.users)
    return;
  }
}).catch(ex => {
  // Exception occured while processing request
});
      </code>
    </pre>
  </div>
  <div id="profiles-java" class="col s12" style="padding:0">
    <pre>
      <code class="java">
// Java client coming soon!      
      </code>
    </pre>
  </div>
 <div id="profiles-python" class="col s12" style="padding:0">
    <pre>
      <code class="python">
from space_api import API

// Initialize api with the project name and url of the space cloud
api = API("books-app", "localhost:8081")

// Initialize database(s) you intend to use
db = api.my_sql()

// Read profiles of all users
response = db.profiles()
if response.status == 200:
    print(response.result)
else:
    print(response.error)

api.close()
      </code>
    </pre>
  </div>
</div>

## Response
A response object sent by the server contains the **status**  and **data** fields explained below:

**status:** Number describing the status of the operation. Following values are possible:
- 200 - Operation was successful
- 401 - Request was unauthenticated
- 403 - Request was unauthorized
- 500 - Internal server error

**data:** The data object consists of one of the following fields:
- **user** (for `profile`) - User object
- **users** (for `profiles`) - Array of user objects 

<div class="btns-wrapper">
  <a href="/docs/user-management/signup" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/database/overview" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>
