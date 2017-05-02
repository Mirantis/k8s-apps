#!/usr/bin/env bash

set -eux

if [ ! -f "${GCLOUD_KEYPATH:-}" ]; then
    echo "Keyfile not found. Please specify path to keyfile in GCLOUD_KEYPATH environtment variable!"
    exit 1
fi

gcloud auth activate-service-account --key-file "${GCLOUD_KEYPATH}"
repo_dir="$(dirname "$(dirname "$(readlink -f $0)")")"
pushd "${repo_dir}/dist/charts"
gsutil -m rsync ./ gs://mirantisworkloads/
