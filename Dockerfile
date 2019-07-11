FROM alpine:3.9
WORKDIR /space-cloud
# Copy the space-cloud binary from the build context to the container's working directory
# COPY space-cloud .
# COPY space-cloud.yaml .
RUN set -ex  \
  && apk add --no-cache ca-certificates wget \
  && wget http://192.168.43.226:8000/space-cloud \
  && wget http://192.168.43.226:8000/config.yaml \
  && chmod +x space-cloud
ENV PROD=false
ENV PATH="/space-cloud:${PATH}"
EXPOSE 8080
CMD ./space-cloud run --config config.yaml