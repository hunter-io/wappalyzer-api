docker-wappalyzer-api
=====

This repository contains a dockerized and 'API-fied' version of [Wappalyzer](https://github.com/AliasIO/Wappalyzer). It aims to make it available through an API endpoint you can call from anywhere. It requires Chrome headless as an execution engine and is built in Go.

## To build it:
```
env GOOS=linux GOARCH=amd64 go build -o server

docker build hunter-io/wappalyzer-api .
```

## To run it:
```
docker network create wappalyzer
docker run -p 9222:9222 --net=wappalyzer knqz/chrome-headless (or any Chrome-headless Docker image)
docker run -p 3001:3001 --net=wappalyzer hunter-io/wappalyzer-api -chromeURL=http://{CHROME_HEADLESS_IP:CHROME_HEADLESS_PORT}/json
```

## To use it:

```
curl -XPOST 'localhost:3001/extract' -d 'url=https://google.com'
```

## License:

Derived work of [Wappalyzer](https://github.com/AliasIO/Wappalyzer/tree/master/src/drivers/npm) and [chromedp](https://github.com/chromedp/chromedp).

Licensed under [GPL-3.0](https://opensource.org/licenses/GPL-3.0).
