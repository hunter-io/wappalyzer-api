const express = require('express')
const Wappalyzer = require('wappalyzer')
const morgan = require('morgan')

const PORT = process.env.PORT || 3000

const app = express()

if (process.env.DISABLE_REQUESTS_LOGGING == undefined) {
  app.use(morgan('combined'))
}

app.get('/', (req, res) => {
  res.send('Wappalyzer API is ready! 🚀')
})

app.get('/extract', (req, res) => {
  let url = req.query.url

  if (url == undefined || url == '') {
    res.status(400).send('missing url query parameter')
    return
  }

  const options = {
    // browser: 'puppeteer',
    debug: false,
    maxDepth: 1,
    recursive: false,
    maxWait: 20000,
    userAgent: 'Wappalyzer',
    htmlMaxCols: 2000,
    htmlMaxRows: 2000,
  }

  const wappalyzer = new Wappalyzer(url, options)
  wappalyzer.analyze()
    .then((json) => {
      res.json(json)
    })
    .catch((error) => {
      res.status(500).send(`${error}\n`)
    })
})

app.listen(PORT, () => console.log(`Starting Wappalyzer on http://0.0.0.0:${PORT}`))
