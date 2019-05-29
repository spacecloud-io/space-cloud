FROM ubuntu

LABEL maintainer="https://rawsanj.dev"

WORKDIR /root

RUN apt update && apt upgrade -y

RUN apt install curl wget unzip -y

RUN wget https://spaceuptech.com/downloads/linux/space-cloud.zip

RUN unzip space-cloud.zip && chmod +x space-cloud

RUN rm space-cloud.zip

RUN wget https://raw.githubusercontent.com/spaceuptech/space-cloud/master/examples/realtime-todo-app/config.yaml

ENTRYPOINT ./space-cloud run --config config.yaml
