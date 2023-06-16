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

app.get('/extract', async (req, res, next) => {
  const url = req.query.url

  if (url == undefined || url == '') {
    return res.status(400).send('missing url query parameter')
  }

  const options = {
    debug: req.query.debug || false,
    maxDepth: req.query.maxDepth || 1,
    recursive: req.query.recursive || false,
    maxWait: req.query.maxWait || 20000,
    userAgent: req.query.userAgent || 'Wappalyzer',
    htmlMaxCols: req.query.htmlMaxCols || 2000,
    htmlMaxRows: req.query.htmlMaxRows || 2000,
  }

  const wappalyzer = new Wappalyzer(options)

  try {
    await wappalyzer.init()

    const site = await wappalyzer.open(url)

    await new Promise((resolve) =>
      setTimeout(resolve, parseInt(options.defer || 0, 10))
    )

    const results = await site.analyze()

    await wappalyzer.destroy()

    res.json(results)
  } catch (error) {
    res.status(500).send(`${error}\n`)
  }
})

app.listen(PORT, '0.0.0.0', () => console.log(`Starting Wappalyzer on http://0.0.0.0:${PORT}`))

process.on('uncaughtException', function (err) {
  console.error((new Date).toUTCString() + ' uncaughtException:', err.message)
  console.error(err.stack)
  process.exit(1)
})
