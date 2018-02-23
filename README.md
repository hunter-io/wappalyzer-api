docker-wappalyzer-api
=====

This repository contains a dockerized and 'API-fied' version of [Wappalyzer](https://github.com/AliasIO/Wappalyzer). It aims to make it available through an API endpoint you can call from anywhere. It uses Chrome Headless as execution engine and is built in Go.

## To build it:
```
env GOOS=linux GOARCH=amd64 go build -o server

docker build hunter-io/wappalyzer-api .
```

## To run it:
```
docker run --name wappalyzer-api --rm -p 3001:3001 hunter-io/wappalyzer-api
```

## To use it:

```
curl -XPOST 'localhost:3001/extract' -d 'url=https://google.com'
```

## License:

Derived work of [Wappalyzer](https://github.com/AliasIO/Wappalyzer/tree/master/src/drivers/npm) and [Automatic Google Chrome Debugger](https://github.com/wirepair/autogcd).

Licensed under [GPL-3.0](https://opensource.org/licenses/GPL-3.0).
