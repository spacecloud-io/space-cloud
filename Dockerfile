FROM golang:1.12 as builder
WORKDIR /space-cloud
COPY . /space-cloud
RUN GOOS=linux GOARCH=amd64 go install -ldflags '-s -w -extldflags "-static"'

FROM alpine:3.9
COPY --from=builder /go/bin/space-cloud /usr/local/bin
RUN chmod -R ugo=rx /usr/local/bin/
RUN chmod ugo=rx /usr/local/bin/space-cloud
EXPOSE 8080

