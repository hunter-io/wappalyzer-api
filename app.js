const express = require('express')
const Wappalyzer = require('wappalyzer')
const morgan = require('morgan')

const PORT = process.env.PORT || 3000

const app = express()

if (process.env.DISABLE_REQUESTS_LOGGING == undefined) {
  app.use(morgan('combined'))
}

function modifyURLs(urls) {
  for (let urlKey in urls) {
    let url = urls[urlKey]

    if (url.hasOwnProperty('error') && typeof url.error === 'string') {
      url.error = {
        type: url.error,
        message: "Response was not ok"
      }
    }

    if (url.status === 200) {
      url.statusAfterRedirects = 200
      url.urlAfterRedirects = urlKey
    }

    if (url.status === 0) {
      if (!url.error) url.error = {}
      url.error.type = "RESPONSE_NOT_OK"
    }
  }

  return urls
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
    debug: false,
    maxDepth: 1,
    recursive: false,
    maxWait: 20000,
    userAgent: 'Wappalyzer',
    htmlMaxCols: 2000,
    htmlMaxRows: 2000,
  }

  const wappalyzer = new Wappalyzer(options)

  try {
    await wappalyzer.init()

    const site = await wappalyzer.open(url)

    await new Promise((resolve) =>
      setTimeout(resolve, parseInt(options.defer || 0, 10))
    )

    let analyzeResult = await site.analyze()
    let results = analyzeResult

    if (req.query.backward_compatible === 'true') {
      let { technologies: applications, ...rest } = analyzeResult
      results = { applications, ...rest }

      results.applications = results.applications.map(app => {
        app.categories = app.categories.map(category => ({
          ...category,
          id: String(category.id)
        }))

        return app
      })

      let urls = results.urls
      urls = modifyURLs(urls)
      results.urls = urls

      results.urls = urls

    }

    wappalyzer.destroy()

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
