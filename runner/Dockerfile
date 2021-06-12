FROM golang:1.15.3-alpine3.12
WORKDIR /build
COPY . .
#RUN apk --no-cache add build-base
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o app .

FROM alpine:3.12
RUN apk --no-cache add ca-certificates && \
    addgroup -S spacecloud && adduser -S spacecloud -G spacecloud && \
    mkdir -p /home/spacecloud/.space-cloud/mission-control-v$SC_VERSION && \
    chown -R spacecloud:spacecloud /home/spacecloud/.space-cloud && \
    mkdir /app && \
    chown -R spacecloud:spacecloud /app 

COPY --chown=spacecloud:spacecloud --from=0 /build/build /home/spacecloud/.space-cloud/mission-control-v$SC_VERSION/build
COPY --chown=spacecloud:spacecloud --from=0 /build/app /app

USER spacecloud:spacecloud
WORKDIR /app
CMD ["./app", "start"]
