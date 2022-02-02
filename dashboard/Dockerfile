FROM node:12 as build

WORKDIR /usr/src/app

COPY rollup.config.js ./
COPY package*.json ./

RUN npm install

COPY ./src ./src
COPY ./public ./public
RUN ls -la ./src

# build webapp
RUN npm run build


FROM nginx:1.19.2-alpine as runtime
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build ./usr/src/app/public /var/share/nginx/html
# rewrite localhost:8080/todolist to https://sillytodonet-n8a9x.ondigitalocean.app/api/todolist (we need to escape (\) . and : and / in the urls)
# RUN sed -i "s/http\:\/\/localhost\:8080\/todolist/https\:\/\/sillytodonet-n8a9x\.ondigitalocean\.app\/api\/todolist/g" /var/share/nginx/html/build/bundle.js
EXPOSE 80
