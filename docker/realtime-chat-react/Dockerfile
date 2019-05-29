FROM node AS node-build-env

LABEL maintainer="https://rawsanj.dev"

WORKDIR /root

RUN apt update && apt upgrade -y

RUN git clone https://github.com/spaceuptech/sample-apps.git

WORKDIR sample-apps/realtime-chat-react

RUN npm install -y

RUN npm run-script build

RUN cp -r build /app

FROM nginx:1.13.0-alpine

WORKDIR /usr/share/nginx/html

COPY --from=node-build-env app ./
