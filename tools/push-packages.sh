#!/usr/bin/env bash

org="mirantisworkloads"

if [ $# -eq 1 ]; then
    org="$1"
fi

set -eux

if [ ! -f "${GCLOUD_KEYPATH:-}" ]; then
    echo "Keyfile not found. Please specify path to keyfile in GCLOUD_KEYPATH environtment variable!"
    exit 1
fi

gcloud auth activate-service-account --key-file "${GCLOUD_KEYPATH}"
repo_dir="$(dirname "$(dirname "$(realpath $0)")")"
pushd "${repo_dir}/dist/charts"
gsutil -m rsync ./ gs://${org}/
