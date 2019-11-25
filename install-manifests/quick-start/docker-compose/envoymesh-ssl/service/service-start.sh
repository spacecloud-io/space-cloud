#!/bin/bash

/usr/local/bin/space-cloud run --admin-user="admin" --admin-pass="admin" --admin-secret="topsecret" & 
## /usr/local/bin/space-cloud run --dev &
/usr/local/bin/envoy -c /etc/service-envoy-sc.yaml --service-cluster service${SERVICE_NAME}