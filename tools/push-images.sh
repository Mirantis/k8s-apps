#!/usr/bin/env bash

org="mirantisworkloads"

if [ $# -eq 1 ]; then
    org="$1"
fi

set -eux

repo_dir="$(dirname "$(dirname "$(realpath $0)")")"
for d in ${repo_dir}/images/*/; do
    pushd "${d}"
    image_name="$(basename "${d}")"
    if [ ! -f "${d}/.version" ]; then
        echo "Image version is not specified for ${image_name}, skipping..."
        continue
    fi
    version="$(cat ${d}/.version | xargs)"
    echo "Push ${image_name} image..."
    docker push "${org}/${image_name}:${version}"
    popd
done
