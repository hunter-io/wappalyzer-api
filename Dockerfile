FROM node:12-alpine

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD true
# ENV CHROME_BIN /usr/bin/chromium-browser

RUN apk update && apk add --no-cache \
	nodejs \
	nodejs-npm \
  udev \
  # chromium \
  ttf-freefont

ADD *.json /
ADD *.js /

RUN npm i

# RUN /usr/bin/chromium-browser --version

ENTRYPOINT ["node", "app.js"]
