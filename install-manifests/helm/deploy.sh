#!/bin/sh

# Run this script using command bash deploy.sh
# Ensure that you are running this command from (space-cloud/install-manifests/helm) directory
# This script assumes that you have configured gcloud command, test if gcloud is configured using command (gcloud auth print-access-token)
# If gcloud is not configured comment the curl command & upload the generated .tgz files manually

arr=(mongo mysql postgres sqlserver space-cloud)

for name in ${arr[@]}; do
  cd ${name}
  fileName=${name}.tgz
  echo ${fileName}
  tar -cvzf ${fileName} .
  curl -v --upload-file ${fileName} \
    -H "Authorization: Bearer `gcloud auth print-access-token`" \
    'https://storage.googleapis.com/space-cloud/helm/'
  rm ${fileName}
  cd ..
done