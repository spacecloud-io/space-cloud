FROM golang:1.13.5-alpine3.10
WORKDIR /build

# Take the current space cloud version as a argument
ARG SC_VERSION=0.18.4

# Copy all the source files
COPY . .
# Install the required packages
RUN apk --no-cache add ca-certificates wget unzip

# Build SC
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o app .

# Download mission control
RUN echo $SC_VERSION && wget https://storage.googleapis.com/space-cloud/mission-control/mission-control-v$SC_VERSION.zip && unzip mission-control-v$SC_VERSION.zip

FROM alpine:3.10
ARG SC_VERSION=0.18.4

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

CMD ["./app", "run"]
