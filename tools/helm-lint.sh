#!/bin/bash

set -e

if ! which helm; then
    echo "Helm is not installed"
    exit 1
fi


for chart in $(find -name 'Chart.yaml'); do
    chart_dir=$(dirname "${chart}")
    chart_name=$(basename "${chart_dir}")
    echo "Inspect ${chart_name}"
    pushd "${chart_dir}"
    helm lint
    popd
done
