const express     = require('express');
const wappalyzer  = require('wappalyzer');
const validUrl    = require('valid-url');
const bodyParser  = require('body-parser');
const cluster     = require('cluster');

const PORT        = 3001;

const app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

app.get('/', function(req, res) {
  res.send('OK');
});

app.post('/extract', function(req, res) {
  var url = req.body.url;
  var timeout = parseInt(req.body.timeout);

  if(isNaN(timeout) || timeout < 5000) {
    timeout = 5000;
  }

  console.log(`Extracting technologies for ${url} with a ${timeout}ms timeout`);

  if (validUrl.isUri(url)) {
    wappalyzer.run([url, '--quiet', `--resource-timeout=${timeout}`], function(stdout, stderr) {
      status = 200;
      response = stdout;

      if(stderr) {
        status = 400;
        response = stderr;
      }

      if(res.headersSent == false) {
        res.set("Connection", "close").status(status).send(response);
      }
      else {
        res.end();
      }
    });
  } else {
    res.set("Connection", "close").status(422).end();
  }
});

app.listen(PORT);

console.log(`Listening on port ${PORT}`)
