const express = require('express')
const app = express()
const port = 3000

app.get('/docs', (req, res) => res.redirect('/docs/getting-started/overview'))
app.get('/docs/:dir/:file', (req, res) => {
  
})


app.listen(port, () => console.log(`Example app listening on port ${port}!`))