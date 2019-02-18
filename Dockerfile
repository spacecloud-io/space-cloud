FROM alpine:3.9
WORKDIR /space-cloud
RUN set -ex \
  && apk add --no-cache ca-certificates \
  && apk add --no-cache unzip \
  && wget https://spaceuptech.com/downloads/linux/space-cloud.zip \
  && unzip space-cloud.zip \
  && rm space-cloud.zip \
  && chmod +x space-cloud
ENV PROD=false
ENV PATH="/space-cloud:${PATH}"
EXPOSE 8080
CMD space-cloud run