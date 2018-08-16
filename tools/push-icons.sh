#!/usr/bin/env bash

org="mirantisworkloads"

if [ $# -eq 1 ]; then
    org="$1"
fi

set -eux

if [ ! -f "${GCLOUD_KEYPATH:-}" ]; then
    echo "Keyfile not found. Please, specify path to keyfile in GCLOUD_KEYPATH environment variable!"
    exit 1
fi

if [ -z "${GCLOUD_PROJECT:-}" ]; then
    echo "Project id not found. Please, specify project id in GCLOUD_PROJECT environment variable!"
    exit 1
fi

gcloud auth activate-service-account --key-file "${GCLOUD_KEYPATH}"
gcloud config set project "${GCLOUD_PROJECT}"

repo_dir="$(dirname "$(dirname "$(realpath $0)")")"
pushd "${repo_dir}/icons"

gsutil -m rsync ./ gs://${org}/icons/
gsutil -m acl set -R -a public-read gs://${org}
