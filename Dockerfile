FROM node:boron

RUN mkdir -p /usr/src/wappalyzer-api
WORKDIR /usr/src/wappalyzer-api

COPY package.json /usr/src/wappalyzer-api/
RUN npm install

COPY . /usr/src/wappalyzer-api

EXPOSE 3001

CMD [ "npm", "start" ]
