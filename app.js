const express = require('express')
const app = express()
const port = 3000

app.get('/', (req, res) => {
  res.send('Wappalyzer API is ready! 🚀')
})

app.get('/extract', (req, res) => {
  const Wappalyzer = require('wappalyzer')

  // TODO: Handle missing URL

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
      res.send(`${JSON.stringify(json, null, 2)}`)
    })
    .catch((error) => {
      res.status(500).send(`${error}\n`)
    })
})

app.listen(port, () => console.log(`Starting Wappalyzer on http://localhost:${port}`))