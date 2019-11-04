FROM alpine:3.9
WORKDIR /space-cloud
# Copy the space-cloud binary from the build context to the container's working directory
COPY space-cloud .
RUN set -ex  \
  && apk add --no-cache ca-certificates wget \
  && chmod +x space-cloud \ 
  && mkdir -p /root/.space-cloud/mission-control-v0.13.0

COPY mission-control /root/.space-cloud/mission-control-v0.13.0

ENV PROD=false
ENV PATH="/space-cloud:${PATH}"

# ports for the http and https servers
EXPOSE 4122 4126

CMD ./space-cloud run