FROM alpine:3.9

WORKDIR /space-cloud

# Copy the space-cloud binary from the build context to the container's working directory
COPY space-cloud .

# COPY space-cloud.yaml .
RUN set -ex  \
  && apk add --no-cache ca-certificates \
  && chmod +x space-cloud \ 
  && mkdir -p /root/.space-cloud/mission-control-v0.12.1

COPY mission-control /root/.space-cloud/mission-control-v0.12.1

ENV PATH="/space-cloud:${PATH}"

# ports for the http and https servers
EXPOSE 4122 4124 4126 4128 4222 4248 4232 4234

CMD ./space-cloud run