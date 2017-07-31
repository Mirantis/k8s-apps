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
    echo -e "${COLOR_BLUE}[$(date +"%F %T")] twitter-stats [${COMMAND}] ${COLOR_RED}|${COLOR_RESET}" $@${COLOR_GRAY}
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
    ensure_command_exists ${TS_KUBECTL_CMD}
    ensure_command_exists ${TS_HELM_CMD}
    ensure_command_exists jq

    # TODO(slukjanov): check versions of kubectl, k8s and helm
}

function release_name() {
    chart="$1"
    echo "ts-${TS_NAME}-${chart}"
}

function _kubectl() {
    ${TS_KUBECTL_CMD} "$@"
}

function _helm() {
    ${TS_HELM_CMD} "$@"
}

WORKDIR="$(dirname ${PWD}/${0})"

function main() {
    header "Welcome to the Twitter Stats!"

    log "Legend (available commands):"
    log "\tup - deploy or upgrade Twitter Stats pipeline with all dependencies, not waiting for completion"
    log "\t\t(it's safe to re-run up many times to apply changes to the configs)"
    log "\twait - waiting for Twitter Stats ready and returns stats (requires only NAMESPACE to be set)"
    log "\tdown - destroy Twitter Stats deployment (requires only NAMESPACE to be set)"
    log "\ttest - deploy, wait and destroy Twitter Stats"

    if [ "${ARGS_NUMBER}" -lt "1" ]; then
        log "Usage: twitter-stats.sh <up | test | down >"
        log "Example: \"twitter-stats.sh up test\" to get Twitter Stats ready and serve stats successfully"
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

    export TS_NAME=${TS_NAME:-demo}
    if [ -z "${TS_NAME}" ] || [ "${TS_NAME}" == "rand" ] ; then
        export TS_NAME=$(hexdump -n 4 -e '1/4 "%08x" 1 "\n"' /dev/random)
    fi
    log "\tTS_NAME - name of the Twitter Stats deployment"
    log "\t\tchange it to do multiple deployments even in single namespace (used in pod names, etc.)"
    log "\t\tuse 'rand' to generate random name"
    log "\t\tcurrent: ${TS_NAME} (default: demo)"

    export TS_NAMESPACE=${TS_NAMESPACE:-}
    if [ -z "${TS_NAMESPACE}" ] ; then
        export TS_NAMESPACE=${TS_NAME}
    fi
    log "\tTS_NAMESPACE - K8s namespace that will be used for deployment"
    log "\t\tcurrent: ${TS_NAMESPACE} (default: same as TS_NAME)"

    export TS_DELETE_NS=${TS_DELETE_NS:-yes}
    log "\tTS_DELETE_NS - Delete or not K8s namespace when 'down' called"
    log "\t\tcurrent: ${TS_DELETE_NS} (default: yes)"

    log "\tTS_APP_KEY, TS_APP_SECRET, TS_TOKEN_KEY, TS_TOKEN_SECRET - Twitter API credentials (must be read+write)"
    export TS_APP_KEY=${TS_APP_KEY:-}
    export TS_APP_SECRET=${TS_APP_SECRET:-}
    export TS_TOKEN_KEY=${TS_TOKEN_KEY:-}
    export TS_TOKEN_SECRET=${TS_TOKEN_SECRET:-}

    if [ -z "${TS_APP_KEY}" ] || [ -z "${TS_APP_SECRET}" ] || [ -z "${TS_TOKEN_KEY}" ] || [ -z "${TS_TOKEN_SECRET}" ] ; then
        log_error "TS_APP_KEY, TS_APP_SECRET, TS_TOKEN_KEY, TS_TOKEN_SECRET must be specified"
    fi

    export TS_MODE=${TS_MODE:-multi}
    log "\tTS_MODE - single or multi-node deployment"
    log "\t\tcurrent: ${TS_MODE} (default: multi)"
    if [ "${TS_MODE}" != "single" ] && [ "${TS_MODE}" != "multi" ] ; then
        log_error "TS_MODE could be only 'single' or 'multi'"
    fi

    export TS_STORAGE=${TS_STORAGE:-hdfs}
    log "\tTS_STORAGE - hdfs or cassandra for storing processed data"
    log "\t\tcurrent: ${TS_STORAGE} (default: hdfs)"
    if [ "${TS_STORAGE}" != "hdfs" ] && [ "${TS_STORAGE}" != "cassandra" ] ; then
        log_error "TS_STORAGE could be only 'hdfs' or 'cassandra'"
    fi

    export TS_RETRIES=${TS_RETRIES:-60}
    log "\tTS_RETRIES - number of retries for 'test' command"
    log "\t\tcurrent: ${TS_RETRIES} (default: 60)"

    export TS_RETRY_INTERVAL=${TS_RETRY_INTERVAL:-10}
    log "\tTS_RETRY_INTERVAL - amount of time in seconds to sleep between retries"
    log "\t\tcurrent: ${TS_RETRY_INTERVAL} (default: 10)"

    export TS_KUBECTL_CMD=${TS_KUBECTL_CMD:-kubectl}
    log "\tTS_KUBECTL_CMD - path to kubectl binary to run"
    log "\t\tcurrent: ${TS_KUBECTL_CMD} (default: kubectl)"

    export TS_HELM_CMD=${TS_HELM_CMD:-helm}
    log "\tTS_HELM_CMD - path to helm binary to run"
    log "\t\tcurrent: ${TS_HELM_CMD} (default: helm)"

    export TS_CHARTS=${TS_CHARTS:-"zookeeper ${TS_STORAGE} kafka spark tweepub tweetics tweeviz"}
    log "\tTS_CHARTS - list of Helm charts to be deployed"
    log "\t\tcurrent: ${TS_CHARTS} (default: zookeeper \${TS_STORAGE} kafka spark tweepub tweetics tweeviz)"

    # Check that all dependencies installed
    check_dependencies

    # Calculated params
    export TS_ZOOKEEPER_RELEASE=$(release_name zookeeper)
    export TS_KAFKA_RELEASE=$(release_name kafka)
    export TS_SPARK_RELEASE=$(release_name spark)
    export TS_HDFS_RELEASE=$(release_name hdfs)
    export TS_CASSANDRA_RELEASE=$(release_name cassandra)

    header "Following commands will be executed: ${ARGS}"

    for command in ${ARGS}; do
        COMMAND="$(printf %4s ${command})"
        command_${command}
    done

    log ""
    log "Successfully finished"
}

function command_up() {
    header "Deploying or Upgrading Twitter Stats services"

    _kubectl get ns ${TS_NAMESPACE} 1>/dev/null 2>/dev/null || _kubectl create ns ${TS_NAMESPACE}

    local tmp=$(mktemp -d)
    log "Calculated configs and Helm logs: ${tmp}"

    cp -r ${WORKDIR}/${TS_MODE}-node/* ${tmp}/

    pushd ${tmp} 1>/dev/null

    # Apply TS_ env variables to files
    for param in $(env | grep TS_); do
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
    _helm repo list | grep -q ts-${TS_NAME} && _helm repo remove ts-${TS_NAME} 1>/dev/null
    # TODO make URL configurable
    _helm repo add ts-${TS_NAME} https://mirantisworkloads.storage.googleapis.com 1>/dev/null
    _helm repo update 1>/dev/null

    for chart in ${TS_CHARTS}; do
        local release=$(release_name ${chart})

        log "Deploying or upgrading chart ${chart} with release ${release}"

        _helm upgrade --install --namespace ${TS_NAMESPACE} ${release} ts-${TS_NAME}/${chart} -f configs/${chart}.yaml | tee ${chart}.log
    done

    _helm repo list | grep -q ts-${TS_NAME} && _helm repo remove ts-${TS_NAME} 1>/dev/null

    log "All Twitter Stats services deployed or upgraded (async)"

    popd 1>/dev/null
}

function command_down() {
    header "Destroying Twitter Stats"

    if [ "${TS_DELETE_NS}" == "yes" ] ; then
        _kubectl get ns | grep -q ${TS_NAMESPACE} && _kubectl delete ns ${TS_NAMESPACE}
    fi

    for chart in ${TS_CHARTS}; do
        local release=$(release_name ${chart})

        log "Destroying chart ${chart} with release ${release}"

        _helm delete --purge ${release}
    done
}

function command_test() {
    header "Wait Twitter Stats ready and serve stats"

    log "Will retry ${TS_RETRIES} times with ${TS_RETRY_INTERVAL} seconds intervals"

    retries=0
    until [ ${retries} -ge "${TS_RETRIES}" ] ; do
        local release=$(release_name tweeviz)
        local service="tweeviz-${release}"

        node_ips=$(_kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="ExternalIP")].address }')
        if [ -z "${node_ips}" ] ; then
            node_ips=$(_kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="LegacyHostIP")].address }')
        fi
        if [ -z "${node_ips}" ] ; then
            log_error "There are no External IPs for K8s nodes available, need at least one to access NodePort"
        fi

        for node_ip in ${node_ips} ; do
            node_port=$(_kubectl -n ${TS_NAMESPACE} get svc ${service} -o jsonpath='{ $.spec.ports[?(@.port==8589)].nodePort }')

            url="http://${node_ip}:${node_port}"

            if [ "$(curl -m 10 -f ${url}/stats 2>/dev/null | jq -r '.popularity[0].weight')" -gt "0" 2>/dev/null ] ; then
                header "Deployed services endpoints"

                spark_url="http://${node_ip}:$(_kubectl -n ${TS_NAMESPACE} get svc spark-master-ext-ts-${TS_NAMESPACE}-spark -o jsonpath='{ $.spec.ports[?(@.port==8080)].nodePort }')"
                log "Spark Web UI: ${spark_url}"

                hdfs_url="http://${node_ip}:$(_kubectl -n ${TS_NAMESPACE} get svc hdfs-ui-ts-${TS_NAMESPACE}-hdfs -o jsonpath='{ $.spec.ports[?(@.port==50070)].nodePort }')"
                log "HDFS Web UI: ${hdfs_url}"

                header "Twitter Stats ready and serve stats successfully"
                log "Link: ${url}"

                return
            else
                log "Twitter Stats isn't ready or not serving stats (yet), sleeping for ${TS_RETRY_INTERVAL} seconds (${retries}/${TS_RETRIES})"
            fi
        done

        sleep ${TS_RETRY_INTERVAL}

        retries=$[${retries}+1]
    done

    log_error "Twitter Stats isn't ready or not serving stats (yet)"
}

main
