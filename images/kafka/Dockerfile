FROM java:8-jre

ARG KAFKA_VERSION=0.11.0.0
ARG KAFKA_SCALA_VERSION=2.12
ARG JMX_VERSION=0.9

COPY entrypoint.sh waitForZK.sh /usr/bin/
COPY check_status.py /

RUN set -x && \
    curl -Lo /tmp/kafka_$KAFKA_SCALA_VERSION-$KAFKA_VERSION.tgz http://archive.apache.org/dist/kafka/$KAFKA_VERSION/kafka_$KAFKA_SCALA_VERSION-$KAFKA_VERSION.tgz && \
    tar -xvzf /tmp/kafka_$KAFKA_SCALA_VERSION-$KAFKA_VERSION.tgz -C /opt && \
    rm /tmp/kafka_$KAFKA_SCALA_VERSION-$KAFKA_VERSION.tgz && \
    mv /opt/kafka_$KAFKA_SCALA_VERSION-$KAFKA_VERSION /opt/kafka && \
    useradd --user-group --create-home --home-dir /var/lib/kafka kafka && \
    apt-get update && \
    apt-get install -y --no-install-recommends dnsutils python3 python3-kazoo && \
    apt-get clean && \
    chmod +x /usr/bin/entrypoint.sh /usr/bin/waitForZK.sh && \
    chown -R kafka:kafka /opt/kafka/ && \
    wget -P /opt/kafka https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/$JMX_VERSION/jmx_prometheus_javaagent-$JMX_VERSION.jar

USER kafka
