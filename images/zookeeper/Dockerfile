FROM java:8-jre

ARG JMX_VERSION=0.9
ARG ZOOKEEPER_VERSION=3.5.3-beta

COPY entrypoint.sh zkReconfigAdd.sh zkReconfigRemove.sh zkCheck.sh /usr/bin/

RUN set -x && \
    curl https://storage.googleapis.com/mirantisworkloads/bin/zookeeper/zookeeper-${ZOOKEEPER_VERSION}.tar.gz -o zookeeper.tar.gz && \
    tar xvf zookeeper.tar.gz && \
    mv zookeeper-${ZOOKEEPER_VERSION} /opt/zookeeper && \
    rm zookeeper.tar.gz && \
    useradd --user-group --create-home --home-dir /var/lib/zookeeper zookeeper && \
    mkdir -p /var/log/zookeeper && chown -R zookeeper:zookeeper /var/log/zookeeper && \
    apt-get update && \
    apt-get install -y --no-install-recommends dnsutils netcat-openbsd && \
    apt-get clean && \
    chmod +x /usr/bin/entrypoint.sh /usr/bin/zkReconfigAdd.sh /usr/bin/zkReconfigRemove.sh /usr/bin/zkCheck.sh && \
    wget -P /opt/zookeeper https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/$JMX_VERSION/jmx_prometheus_javaagent-$JMX_VERSION.jar

USER zookeeper
