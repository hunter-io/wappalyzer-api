const express     = require('express');
const wappalyzer  = require('wappalyzer');
const validUrl    = require('valid-url');
const bodyParser  = require('body-parser');
const cluster     = require('cluster');

const NUM_CPUS    = require('os').cpus().length;
const PORT        = 3001;


if(cluster.isMaster) {
    console.log(`Master cluster setting up ${NUM_CPUS} workers`);

    for (let i = 0; i < NUM_CPUS; i++) {
      cluster.fork();
    }

    cluster.on('exit', (worker, code, signal) => {
      console.log(`Worker ${worker.process.pid} died`);
    });
} else {
    const app = express();
    app.use(bodyParser.json());
    app.use(bodyParser.urlencoded({ extended: true }));

    app.get('/', function(req, res) {
      res.send('OK');
    });

    app.post('/extract', function(req, res) {
      var url = req.body.url;

      console.log(`Extracting technologies for ${url}`);

      if (validUrl.isUri(url)) {
        wappalyzer.run([url, '--quiet'], function(stdout, stderr) {
          if (stdout) {
            res.send(stdout);
          }
          if (stderr) {
            res.sendStatus(400).send(stderr);
          }
        });
      } else {
        res.sendStatus(422);
      }
    });

    app.listen(PORT);

    console.log(`Worker ${process.pid} listening on port ${PORT}`)
}
