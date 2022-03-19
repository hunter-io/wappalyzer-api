FROM node:14-alpine

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD true
ENV CHROME_BIN /usr/bin/chromium-browser

RUN apk update
RUN apk add --no-cache \
	nodejs \
	npm \
  udev \
  chromium \
  ttf-freefont \
  python3 \
  make \
  build-base \
  git

RUN mkdir /app && chown node /app
USER 1000
WORKDIR /app

ADD *.json /app/
RUN npm i

ADD *.js /app/

RUN /usr/bin/chromium-browser --version

ENTRYPOINT ["node", "app.js"]
