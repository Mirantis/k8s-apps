FROM java:8-jre

ARG ZOOKEEPER_VERSION=3.5.2-alpha

COPY entrypoint.sh zkReconfigAdd.sh zkReconfigRemove.sh /usr/bin/

RUN set -x && \
    curl -Lo /tmp/zookeeper-$ZOOKEEPER_VERSION.tar.gz http://apache.mirrors.pair.com/zookeeper/zookeeper-$ZOOKEEPER_VERSION/zookeeper-$ZOOKEEPER_VERSION.tar.gz && \
    tar -xvzf /tmp/zookeeper-$ZOOKEEPER_VERSION.tar.gz -C /opt && \
    rm /tmp/zookeeper-$ZOOKEEPER_VERSION.tar.gz && \
    mv /opt/zookeeper-$ZOOKEEPER_VERSION /opt/zookeeper && \
    useradd --user-group --create-home --home-dir /var/lib/zookeeper zookeeper && \
    mkdir -p /var/log/zookeeper && chown -R zookeeper:zookeeper /var/log/zookeeper && \
    apt-get update && \
    apt-get install -y --no-install-recommends dnsutils netcat-openbsd && \
    apt-get clean && \
    chmod +x /usr/bin/entrypoint.sh /usr/bin/zkReconfigAdd.sh /usr/bin/zkReconfigRemove.sh

USER zookeeper
