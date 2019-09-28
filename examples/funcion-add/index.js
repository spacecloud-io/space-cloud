const { API } = require('space-api');     // Import the space-api library

const projectId = 'my-adder';                   // The name of our project
const spaceCloudURL = 'http://localhost:4122';  // The url of Space Cloud

// Create an space api object
const api = new API(projectId, spaceCloudURL);

// Make a service object
const service = api.Service('arithmetic');

// Register function to a service
service.registerFunc('sum', (params, auth, cb) => {
  
  // params - contains the payload send by the client.
  // auth   - contains the JWT claims of the client making the request.
  // cb     - is the functions used to give back a response to the client
  
  // In our use case:
  // params = { num1: SOME_NUMBER, num2: ANOTHER_NUMBER }

  // Lets add the two numbers
  const sum = params.num1 + params.num2;

  // Return the response to the clinet
  cb('response', { sum: sum });
});

// Start the service
service.start() 

