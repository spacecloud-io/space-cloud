export const defaultTemplate = `// DB objects
const dbMongo = api.Mongo()
const dbMySQL = api.MySQL()
const dbPostgres = api.Postgres()
  
// Your logic
`

export const insertTemplate = `// DB object
const db = api.Mongo() // MySQL() || Postgres()

// Document to be inserted
const doc = { text: "Star Space Cloud!", time: new Date() }  

// Send request to Space Cloud
db.insert("todos").doc(doc).apply()
`

export const getTemplate = `// DB object
const db = api.Mongo() // MySQL() || Postgres()

// Send request to Space Cloud
db.get("todos").apply()
`

export const callTemplate = `// Service name
const serviceName = "my-service"

// Function to be triggered
const funcName = "my-func"

// Params to be sent to the function
const params = { msg: "Function Mesh is awesome!" }
// Send request to Space Cloud to trigger backend function
api.call(serviceName, funcName, params)
`