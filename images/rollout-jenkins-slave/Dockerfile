FROM jenkins/jnlp-slave:2.62

ARG HELM_VERSION=v2.7.2
ARG KUBERNETES_VERSION=v1.8.5

USER root
ADD helm-apply /usr/local/bin/
ADD ssh/* /home/jenkins/.ssh/
ADD tweeviz-ui/* /tweeviz-ui/

RUN apt-get update && apt-get install -y jq libltdl7 gettext-base git && \
    apt-get clean && rm -rf /var/lib/apt/lists/* && \
    curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBERNETES_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl && \
    curl -LO https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    tar -zxvf helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -rf helm-${HELM_VERSION}-linux-amd64.tar.gz linux-amd64 && \
    chmod +x /usr/local/bin/kubectl /usr/local/bin/helm /usr/local/bin/helm-apply
