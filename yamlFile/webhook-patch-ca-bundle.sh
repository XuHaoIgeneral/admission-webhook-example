#!/bin/bash
ROOT=$(cd $(dirname $0)/../../; pwd)
set -o errexit
set -o nounset
set -o pipefail
export CA_BUNDLE=$(kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
if command -v envsubst >/dev/null 2>&1; then
    envsubst
else
    sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g"
fi