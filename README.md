# Wappalyzer API

This repository contains a dockerized and 'API-fied' version of [Wappalyzer](https://github.com/AliasIO/Wappalyzer). It aims to make it available through an API endpoint you can call from anywhere.

## To build it:
```
docker build -t hunterio/wappalyzer-api .
```

## To run it:
```
docker run -p 3000:3000 hunterio/wappalyzer-api
```

## To use it:
```
curl 'localhost:3000/extract?url=https://hunter.io'
```

## License:
Derived work of [Wappalyzer](https://github.com/AliasIO/Wappalyzer/).
Licensed under [GPL-3.0](https://opensource.org/licenses/GPL-3.0).