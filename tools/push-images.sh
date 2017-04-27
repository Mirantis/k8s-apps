#!/usr/bin/env bash

set -eux

repo_dir="$(dirname "$(dirname "$(readlink -f $0)")")"
for d in ${repo_dir}/images/*/; do
    pushd "${d}"
    image_name="$(basename "${d}")"
    if [ ! -f "${d}/.version" ]; then
        echo "Image version is not specified for ${image_name}, skipping..."
        continue
    fi
    version="$(cat ${d}/.version | xargs)"
    echo "Push ${image_name} image..."
    docker push "mirantisworkloads/${image_name}:${version}"
    popd
done
