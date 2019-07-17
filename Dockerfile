FROM alpine:3.9
WORKDIR /space-cloud
# Copy the space-cloud binary from the build context to the container's working directory
COPY space-cloud .
RUN set -ex  \
  && apk add --no-cache ca-certificates wget \
  && chmod +x space-cloud \ 
  && mkdir -p /root/.space-cloud/mission-control-v0.10.0

COPY mission-control /root/.space-cloud/mission-control-v0.10.0

ENV PROD=false
ENV PATH="/space-cloud:${PATH}"

# ports for the http and https servers
EXPOSE 4122 4124
EXPOSE 4126 4128

# ports for nats
EXPOSE 4222 4248

# ports for gossip and raft
EXPOSE 4232 4234

CMD ./space-cloud run