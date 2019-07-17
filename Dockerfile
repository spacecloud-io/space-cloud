FROM alpine:3.9

WORKDIR /space-cloud

# Copy the space-cloud binary from the build context to the container's working directory
COPY space-cloud .

# COPY space-cloud.yaml .
RUN set -ex  \
  && apk add --no-cache ca-certificates wget \
  && chmod +x space-cloud
  
ENV PROD=false
ENV PATH="/space-cloud:${PATH}"

# ports for the http and https servers
EXPOSE 4242 4244
EXPOSE 4343 4245

# ports for nats
EXPOSE 4222 4248

# ports for gossip and raft
EXPOSE 4232 4234

CMD ./space-cloud run