FROM openjdk:8 as stage0
LABEL snp-multi-stage="intermediate"
LABEL snp-multi-stage-id="58f58868-292b-482f-8850-79b72faa9d52"
WORKDIR /opt/docker
COPY target/docker/stage/1/opt /1/opt
COPY target/docker/stage/2/opt /2/opt
USER root
RUN ["chmod", "-R", "u=rX,g=rX", "/1/opt/docker"]
RUN ["chmod", "-R", "u=rX,g=rX", "/2/opt/docker"]
RUN ["chmod", "u+x,g+x", "/1/opt/docker/bin/db-events-soruce"]

FROM golang:1.15.3-alpine3.12 as stage1
WORKDIR /build
COPY src/conn-string-parser .
#RUN apk --no-cache add build-base
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o app .

FROM openjdk:8 as mainstage
USER root
RUN id -u demiourgos728 1>/dev/null 2>&1 || (( getent group 0 1>/dev/null 2>&1 || ( type groupadd 1>/dev/null 2>&1 && groupadd -g 0 root || addgroup -g 0 -S root )) && ( type useradd 1>/dev/null 2>&1 && useradd --system --create-home --uid 1001 --gid 0 demiourgos728 || adduser -S -u 1001 -G root demiourgos728 ))
WORKDIR /opt/docker
COPY --from=stage0 --chown=demiourgos728:root /1/opt/docker /opt/docker
COPY --from=stage0 --chown=demiourgos728:root /2/opt/docker /opt/docker
COPY --from=stage1 --chown=demiourgos728:root /build/app /usr/local/bin/conn-string-parser
COPY --chown=demiourgos728:root src/main/resources/application.conf /config/application.conf
RUN chmod -R 0777 /opt/docker
USER 1001:0
ENTRYPOINT ["/opt/docker/bin/db-events-soruce"]
CMD []
