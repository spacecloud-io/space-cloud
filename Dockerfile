FROM alpine:3.9
WORKDIR /space-cloud
# Copy the space-cloud binary from the build context to the container's working directory
COPY space-cloud .
RUN set -ex  \
  && apk add --no-cache ca-certificates \
  && chmod +x space-cloud
ENV PROD=false
ENV PATH="/space-cloud:${PATH}"
EXPOSE 8080
CMD ./space-cloud run