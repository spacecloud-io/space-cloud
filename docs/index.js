const fs = require('fs')
const showdown = require('showdown')
const express = require('express')
const handlebars = require('handlebars')

const converter = new showdown.Converter({simpleLineBreaks: true})
const app = express()

var config = []
var template = {}

const configFile = './config.json'
const templateFile = './public/index.handlebars'


const data1 = fs.readFileSync(configFile, 'utf8')
config = JSON.parse(data1)

const data2 = fs.readFileSync(templateFile, 'utf8')
template = handlebars.compile(data2)

const handleIndex = (_, res) => {
  fs.readFile(`./manual/index.md`, 'utf8', (err, data) => {
    if (err) {
      res.status(404)
      res.send()      
      return
    }

    const name = "Documentation"
    res.send(render(data, name, 'overview'))
  })
}

const handlePage = (req, res) => {
  const dir = req.params.dir
  const file = req.params.file === 'overview' ? 'index' : req.params.file

  fs.readFile(`./manual/${dir}/${file}.md`, 'utf8', (err, data) => {
    if (err) {
      res.status(404)
      res.send()
      return
    }

    const name = config.find(page => page.url == dir).name
    res.send(render(data, name, req.params.file))
  })
}

const render = (data, name, pageUrl) => {
  const html = converter.makeHtml(data)
  const pages = config.map(page => Object.assign({}, {isActive: page.name === name}, page, {
    files: [{ title: 'Overview', url: 'overview', isActive: pageUrl === 'overview' && page.name === name  }].concat(page.pages.map(p => ({ title: p[1], url: p[0], isActive: pageUrl === p[0] && page.name === name })))
  }))

  return template({ pages: pages, html: html, name: name })
}

app.use(express.static('public'))

app.get('/docs', handleIndex)
app.get('/docs/:dir', (req, res) => res.redirect(`/docs/${req.params.dir}/overview`))
app.get('/docs/:dir/:file', handlePage)


const port = process.env.POST | 3000
app.listen(port, () => console.log(`Example app listening on port ${port}!`))