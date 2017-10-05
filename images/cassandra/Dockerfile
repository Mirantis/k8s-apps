FROM cassandra:3.11.0

ARG JMX_VERSION=0.9

COPY docker-entrypoint.sh /docker-entrypoint.sh

RUN set -x && \
    apt-get update && apt-get install -y --no-install-recommends wget && \
    apt-get clean && \
    wget -P /etc/cassandra/ https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/$JMX_VERSION/jmx_prometheus_javaagent-$JMX_VERSION.jar && \
    echo 'JVM_OPTS="$JVM_OPTS $JMX_AGENT"' >> /etc/cassandra/cassandra-env.sh
