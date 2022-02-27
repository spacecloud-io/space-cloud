#!/bin/sh
set -e

sbt docker:stage
docker build --no-cache -t spacecloudio/dbevents:0.2.0 .
docker push spacecloudio/dbevents:0.2.0
