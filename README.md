docker-wappalyzer-api
=====

This repository contains a dockerized and 'API-fied' version of [Wappalyzer](https://github.com/AliasIO/Wappalyzer). It aims to make it available through an API endpoint you can call from anywhere. It requires a running Selenium Chrome instance and is built in Go.

## To build it:
```
env GOOS=linux GOARCH=amd64 go build -o server

docker build hunter-io/wappalyzer-api .
```

## To run it:
```
docker network create wappalyzer
docker run -p 4444:4444 --net=wappalyzer selenium/standalone-chrome
docker run -p 3001:3001 --net=wappalyzer hunter-io/wappalyzer-api -seleniumURL=http://{SELENIUM_CHROME_IP:SELENIUM_CHROME_PORT}/wd/hub
```

## To use it:

```
curl -XPOST 'localhost:3001/extract' -d 'url=https://google.com'
```

## License:

Derived work of [Wappalyzer](https://github.com/AliasIO/Wappalyzer/tree/master/src/drivers/npm), [Docker-Selenium](https://github.com/SeleniumHQ/docker-selenium) and [tebeka/selenium](https://github.com/tebeka/selenium) as the WebDriver client.

Licensed under [GPL-3.0](https://opensource.org/licenses/GPL-3.0).
