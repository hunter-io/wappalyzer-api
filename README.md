docker-wappalyzer-api
=====

This repository contains a dockerized version of [the npm driver for Wappalyzer](https://github.com/AliasIO/Wappalyzer/tree/master/src/drivers/npm).

## To run it:
```
docker run --name wappalyzer-api --rm -p 3001:3001 bastienl/wappalyzer-api
```

## To use it:

`x-www-form-urlencoded` style:
```
curl -XPOST 'localhost:3001/extract?pretty' -d 'url=https://google.com'
```

JSON style:
```
curl -XPOST 'localhost:3001/extract?pretty' -H "Content-Type: application/json" -d '{"url": "https://google.com"}'
```

## License:

Derived work of [Wappalyzer](https://github.com/AliasIO/Wappalyzer/tree/master/src/drivers/npm).

Licensed under [GPL-3.0](https://opensource.org/licenses/GPL-3.0).

