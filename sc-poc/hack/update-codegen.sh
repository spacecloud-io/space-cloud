#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go mod vendor


# SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
# CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

SCRIPT_ROOT=".."
CODEGEN_PKG="./vendor/k8s.io/code-generator"

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
bash "${CODEGEN_PKG}"/generate-groups.sh "client,informer,lister" \
  github.com/spacecloud-io/space-cloud/pkg/client github.com/spacecloud-io/space-cloud/pkg/apis \
  core:v1alpha1 \
  --output-base "./hack" \
  --go-header-file "./hack/boilerplate.go.txt"

rm -rf ./pkg/client
cp -rf ./hack/github.com/spacecloud-io/space-cloud/pkg ./
rm -rf ./hack/github.com