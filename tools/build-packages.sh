#!/usr/bin/env bash

org="mirantisworkloads"

if [ $# -eq 1 ]; then
    org="$1"
fi

set -eux

HELM_CMD=${HELM_CMD:-helm}

REPO_URL="https://${org}.storage.googleapis.com"
repo_dir="$(dirname "$(dirname "${PWD}/$0")")"

# build dependencies
for i in {1..3}; do
    for d in ${repo_dir}/charts/*/; do
        pushd "${d}"
        rm -fr ./charts
        $HELM_CMD dep up
        popd
    done
done

# build packages
packages_dir="${repo_dir}/dist/charts"
mkdir -p "${packages_dir}"
pushd "${packages_dir}"
for d in ${repo_dir}/charts/*/; do
    $HELM_CMD package "${d}"
done

# generate index
wget "${REPO_URL}/index.yaml" -O index.yaml
$HELM_CMD repo index --url "${REPO_URL}" --merge index.yaml .
