const express     = require('express');
const wappalyzer  = require('wappalyzer');
const validUrl    = require('valid-url');
const bodyParser  = require('body-parser');

const PORT = 3001;

const app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

app.get('/', function(req, res) {
  res.send('OK');
});

app.post('/extract', function(req, res) {
  var url = req.body.url;

  console.log('Extracting technologies for ' + url);

  if (validUrl.isUri(url)) {
    wappalyzer.run([url, '--quiet'], function(stdout, stderr) {
      if (stdout) {
        res.send(stdout);
      }
      if (stderr) {
        res.status(400).send(stderr);
      }
    });
  } else {
    res.send(422);
  }
});

app.listen(PORT);

console.log('Listening on port ' + PORT);
