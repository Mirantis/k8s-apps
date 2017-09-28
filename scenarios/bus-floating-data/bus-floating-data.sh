#!/bin/bash

set -eou pipefail

ARGS_NUMBER="$#"
ARGS="$@"

COLOR_GRAY='\033[0;37m'
COLOR_BLUE='\033[0;34m'
COLOR_RED='\033[0;31m'
COLOR_RESET='\033[0m'

function finish() {
    COMMAND="exit"
    log ""
    log "Total execution time: ${SECONDS} seconds"
    echo -e -n ${COLOR_RESET}
}
trap finish EXIT

echo -e -n ${COLOR_GRAY}

DEBUG=${DEBUG:-no}
if [ "yes" == "${DEBUG}" ]; then
    set -x
fi

COMMAND="main"

function log() {
    set +x
    echo -e "${COLOR_BLUE}[$(date +"%F %T")] bus-floating-data [${COMMAND}] ${COLOR_RED}|${COLOR_RESET}" $@${COLOR_GRAY}
    if [ "yes" == "${DEBUG}" ] ; then
        set -x
    fi
}

function log_error() {
    log ""
    log "ERROR: $1"
    exit 1
}

function header() {
    log "--------------------------------------------------"
    log $@
    log "--------------------------------------------------"
}

function ensure_command_exists() {
    if ! hash "$1" 2>/dev/null ; then
        log_error "'$1' required and not installed"
    fi
}

function check_dependencies() {
    ensure_command_exists ${BFD_KUBECTL_CMD}
    ensure_command_exists ${BFD_HELM_CMD}
    ensure_command_exists jq

    # TODO(slukjanov): check versions of kubectl, k8s and helm
}

function release_name() {
    chart="$1"
    echo "bfd-${BFD_NAME}-${chart}"
}

function _kubectl() {
    ${BFD_KUBECTL_CMD} "$@"
}

function _helm() {
    ${BFD_HELM_CMD} "$@"
}

WORKDIR="$(dirname ${PWD}/${0})"

function main() {
    header "Welcome to the Bus Floating Data!"

    log "Legend (available commands):"
    log "\tup - deploy or upgrade Bus Floating Data pipeline with all dependencies, not waiting for completion"
    log "\t\t(it's safe to re-run up many times to apply changes to the configs)"
    log "\twait - waiting for Bus Floating Data ready and returns stats (requires only NAMESPACE to be set)"
    log "\tdown - destroy Bus Floating Data deployment (requires only NAMESPACE to be set)"
    log "\ttest - deploy, wait and destroy Bus Floating Data"

    if [ "${ARGS_NUMBER}" -lt "1" ]; then
        log "Usage: bus-floating-data.sh <up | test | down >"
        log "Example: \"bus-floating-data.sh up test\" to get Bus Floating Data ready and serve stats successfully"
        log_error "Too few arguments"
    fi

    # Check that all commands exists
    for command in ${ARGS}; do
        if [ -n "$(type -t command_${command})" ] && [ "$(type -t command_${command})" = function ]; then
            if [ "yes" == "${DEBUG}" ] ; then
                log "\t${command}"
            fi
        else
            log_error "Command '${command}' not found."
        fi
    done

    log ""
    log "Parameters (could be set using environment variables):"

    export BFD_NAME=${BFD_NAME:-demo}
    if [ -z "${BFD_NAME}" ] || [ "${BFD_NAME}" == "rand" ] ; then
        export BFD_NAME=$(hexdump -n 4 -e '1/4 "%08x" 1 "\n"' /dev/random)
    fi
    log "\tBFD_NAME - name of the Bus Floating Data deployment"
    log "\t\tchange it to do multiple deployments even in single namespace (used in pod names, etc.)"
    log "\t\tuse 'rand' to generate random name"
    log "\t\tcurrent: ${BFD_NAME} (default: demo)"

    export BFD_NAMESPACE=${BFD_NAMESPACE:-}
    if [ -z "${BFD_NAMESPACE}" ] ; then
        export BFD_NAMESPACE=${BFD_NAME}
    fi
    log "\tBFD_NAMESPACE - K8s namespace that will be used for deployment"
    log "\t\tcurrent: ${BFD_NAMESPACE} (default: same as BFD_NAME)"

    export BFD_EXTERNAL_IP=${BFD_EXTERNAL_IP:-cluster.local}
    log "\tBFD_EXTERNAL_IP - External IP of one of the k8s cluster nodes."
    log "\tIf not specified, attempt will be made to automatically detect it."
    log "\t\tcurrent: $([ "${BFD_EXTERNAL_IP}" ] || echo unset) (default: unset)"

    export BFD_KUBERNETES_DOMAIN=${BFD_KUBERNETES_DOMAIN:-cluster.local}
    log "\tBFD_KUBERNETES_DOMAIN - K8s cluster domain"
    log "\t\tcurrent: ${BFD_KUBERNETES_DOMAIN} (default: cluster.local)"

    export BFD_DELETE_NS=${BFD_DELETE_NS:-yes}
    log "\tBFD_DELETE_NS - Delete or not K8s namespace when 'down' called"
    log "\t\tcurrent: ${BFD_DELETE_NS} (default: yes)"

    export BFD_USE_INTERNAL_IP=${BFD_USE_INTERNAL_IP:-}
    log "\tBFD_USE_INTERNAL_IP - if set use internal IPs of nodes if everything else fails"
    log "\t\tcurrent: $([ "${BFD_USE_INTERNAL_IP}" ] && echo set || echo unset) (default: unset)"

    export BFD_MODE=${BFD_MODE:-multi}
    log "\tBFD_MODE - single or multi-node deployment"
    log "\t\tcurrent: ${BFD_MODE} (default: multi)"
    if [ "${BFD_MODE}" != "single" ] && [ "${BFD_MODE}" != "multi" ] ; then
        log_error "BFD_MODE could be only 'single' or 'multi'"
    fi

    export BFD_RETRIES=${BFD_RETRIES:-60}
    log "\tBFD_RETRIES - number of retries for 'test' command"
    log "\t\tcurrent: ${BFD_RETRIES} (default: 60)"

    export BFD_RETRY_INTERVAL=${BFD_RETRY_INTERVAL:-10}
    log "\tBFD_RETRY_INTERVAL - amount of time in seconds to sleep between retries"
    log "\t\tcurrent: ${BFD_RETRY_INTERVAL} (default: 10)"

    export BFD_KUBECTL_CMD=${BFD_KUBECTL_CMD:-kubectl}
    log "\tBFD_KUBECTL_CMD - path to kubectl binary to run"
    log "\t\tcurrent: ${BFD_KUBECTL_CMD} (default: kubectl)"

    export BFD_HELM_CMD=${BFD_HELM_CMD:-helm}
    log "\tBFD_HELM_CMD - path to helm binary to run"
    log "\t\tcurrent: ${BFD_HELM_CMD} (default: helm)"

    export BFD_CHARTS=${BFD_CHARTS:-"zookeeper cassandra kafka spark bus-floating-data"}
    log "\tBFD_CHARTS - list of Helm charts to be deployed"
    log "\t\tcurrent: ${BFD_CHARTS} (default: zookeeper cassandra kafka spark bus-floating-data)"

    # Check that all dependencies installed
    check_dependencies

    # Calculated params
    export BFD_ZOOKEEPER_RELEASE=$(release_name zookeeper)
    export BFD_KAFKA_RELEASE=$(release_name kafka)
    export BFD_SPARK_RELEASE=$(release_name spark)
    export BFD_CASSANDRA_RELEASE=$(release_name cassandra)

    header "Following commands will be executed: ${ARGS}"

    for command in ${ARGS}; do
        COMMAND="$(printf %4s ${command})"
        command_${command}
    done

    log ""
    log "Successfully finished"
}

function command_up() {
    header "Deploying or Upgrading Bus Floating Data services"

    _kubectl get ns ${BFD_NAMESPACE} 1>/dev/null 2>/dev/null || _kubectl create ns ${BFD_NAMESPACE}
    local ok=0
    for i in $(seq 1 60); do
        if [[ $(_kubectl get ns ${BFD_NAMESPACE} -o jsonpath="{ .metadata.name }") == "${BFD_NAMESPACE}" ]]; then
            ok=1
            break
        else
            sleep 1
        fi
    done
    if [[ ${ok} != 1 ]]; then
        log_error "Failed to create ${BFD_NAMESPACE} namespace."
    fi

    if [[ "${BFD_EXTERNAL_IP}" == "" ]]; then
        node_ips=$(_kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="ExternalIP")].address }')
        if [ -z "${node_ips}" ] ; then
            node_ips=$(_kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="LegacyHostIP")].address }')
        fi
        if [ -z "${node_ips}" -a "${BFD_USE_INTERNAL_IP}" ] ; then
            node_ips=$(_kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="InternalIP")].address }')
        fi
        if [ -z "${node_ips}" ] ; then
            log_error "There are no External IPs for K8s nodes available, need at least one to access NodePort"
        fi
        export BFD_EXTERNAL_IP=$(echo $node_ips | cut -d" " -f1)
    fi

    local tmp=$(mktemp -d)
    log "Calculated configs and Helm logs: ${tmp}"

    cp -r ${WORKDIR}/${BFD_MODE}-node/* ${tmp}/

    pushd ${tmp} 1>/dev/null

    # Apply BFD_ env variables to files
    for param in $(env | grep BFD_); do
        local param_name=$(echo "${param}" | cut -d'=' -f1)
        local param_value=$(echo "${param}" | cut -d'=' -f2)

        if [ "yes" == "${DEBUG}" ] ; then
            log "Applying param ${param_name}=${param_value} to configs"
        fi

        # Replace $param_name and ${param_name} with param value
        # Not using envsubst to decrease number of dependencies
        find configs/ -type f -name "*.yaml" -exec sed -i="" "s/\$${param_name}/${param_value}/g" {} \;
        find configs/ -type f -name "*.yaml" -exec sed -i="" "s/\${${param_name}}/${param_value}/g" {} \;
    done

    log "Setting up Helm"
    _helm repo list | grep -q bfd-${BFD_NAME} && _helm repo remove bfd-${BFD_NAME} 1>/dev/null
    # TODO make URL configurable
    _helm repo add bfd-${BFD_NAME} https://mirantisworkloads.storage.googleapis.com 1>/dev/null
    _helm repo update 1>/dev/null

    for chart in ${BFD_CHARTS}; do
        local release=$(release_name ${chart})

        log "Deploying or upgrading chart ${chart} with release ${release}"

        _helm upgrade --install --wait --namespace ${BFD_NAMESPACE} ${release} bfd-${BFD_NAME}/${chart} -f configs/${chart}.yaml | tee ${chart}.log
    done

    _helm repo list | grep -q bfd-${BFD_NAME} && _helm repo remove bfd-${BFD_NAME} 1>/dev/null

    log "All Bus Floating Data services deployed or upgraded (async)"

    popd 1>/dev/null
}

function command_down() {
    header "Destroying Bus Floating Data"

    if [ "${BFD_DELETE_NS}" == "yes" ] ; then
        _kubectl get ns | grep -q ${BFD_NAMESPACE} && _kubectl delete ns ${BFD_NAMESPACE}
    fi

    for chart in ${BFD_CHARTS}; do
        local release=$(release_name ${chart})

        log "Destroying chart ${chart} with release ${release}"

        _helm delete --purge ${release}
    done
}

function command_test() {
    header "Wait for Bus Floating Data readiness"

    log "Will retry ${BFD_RETRIES} times with ${BFD_RETRY_INTERVAL} seconds intervals"

    retries=0
    until [ ${retries} -ge "${BFD_RETRIES}" ] ; do
        local release=$(release_name tweeviz)
        local service="tweeviz-${release}"

        url="http://${BFD_EXTERNAL_IP}:8000"

        if curl --silent -m 10 ${url} --output /dev/null ; then
            header "Deployed services endpoints"

            spark_url="http://${BFD_EXTERNAL_IP}:$(_kubectl -n ${BFD_NAMESPACE} get svc spark-master-bfd-${BFD_NAME}-spark-0 -o jsonpath='{ $.spec.ports[?(@.port==8080)].nodePort }')"
            log "Spark Web UI: ${spark_url}"

            header "Bus Floating Data is ready"
            log "Link: ${url}"

            return
        else
            log "Bus Floating Data isn't ready (yet), sleeping for ${BFD_RETRY_INTERVAL} seconds (${retries}/${BFD_RETRIES})"
        fi

        sleep ${BFD_RETRY_INTERVAL}

        retries=$[${retries}+1]
    done

    log_error "Bus Floating Data isn't ready (yet)"
}

main
