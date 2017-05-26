#!/usr/bin/env bash

org="mirantisworkloads"

if [ $# -eq 1 ]; then
    org="$1"
fi

set -eux

REPO_URL="https://${org}.storage.googleapis.com/"
repo_dir="$(dirname "$(dirname "$(realpath $0)")")"

# build dependencies
for d in ${repo_dir}/charts/*/; do
    pushd "${d}"
    helm dep up
    popd
done

# build packages
packages_dir="${repo_dir}/dist/charts"
mkdir -p "${packages_dir}"
pushd "${packages_dir}"
for d in ${repo_dir}/charts/*/; do
    helm package "${d}"
done

# generate index
wget "${REPO_URL}/index.yaml"
helm repo index --url "${REPO_URL}" --merge index.yaml .
