#!/bin/bash

set -e

function usage {
    local base_name=$(basename $0)
    echo "Usage:"
    echo "  $base_name -s <address>"
    echo "  $base_name -n <namespace>"
    echo "  $base_name -i <node>"
    echo "  $base_name -u deploy ui"
}

export NAMESPACE="kube-system"
WORKDIR=$(dirname $0)

while getopts "s:n:i:u" opt; do
    case $opt in
        "s" )
            SRV_OPT=" -s $OPTARG"
            ;;
        "n" )
            NAMESPACE="$OPTARG"
            ;;
        "i" )
            NODE="$OPTARG"
            ;;
        "u" )
            DEPLOY_UI=true
            ;;
        * )
            usage
            exit 1
            ;;
    esac
done

which kubectl 1>/dev/null

function kube_cmd {
    kubectl $SRV_OPT --namespace $NAMESPACE "$@"
}

function await_readiness {
    rc_name=$1;
    template="{{.status.readyReplicas}}/{{.spec.replicas}}"
    get_rc="kube_cmd get rc $rc_name -o template --template "

    echo "Waiting for $rc_name readiness"

    while $get_rc $template | egrep -v '^(.*)/\1$' > /dev/null; do
        sleep 3
    done
    echo "The $rc_name is ready"
}

if [ -z $NODE ]; then
    NODE=$(kubectl get nodes -o template --template="{{ with index .items 0 }}{{ .metadata.name }}{{ end }}")
    echo "K8S node is not specified, using $NODE"
fi

kubectl label node $NODE app=ccp-registry --overwrite

function deploy_registry {
    kube_cmd apply -f $WORKDIR/registry-rc.yaml
    kube_cmd apply -f $WORKDIR/registry-service.yaml
    await_readiness registry
    kube_cmd get service registry
}

function deploy_registry_ui {
    cat $WORKDIR/registry-ui-rc.yaml | envsubst | kube_cmd apply -f -
    kube_cmd apply -f $WORKDIR/registry-ui-service.yaml
    await_readiness registry-ui
    kube_cmd get service registry-ui
}

deploy_registry
if [ $DEPLOY_UI ]; then
    deploy_registry_ui
fi
