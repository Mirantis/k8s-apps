FROM openjdk:8

ENV NODE_DATA=true NODE_MASTER=true BOOTSTRAP_MLOCKALL=false ES_JAVA_OPTS=-Djava.net.preferIPv4Stack=true

RUN set -x \
    && curl -o /tmp/esc.tgz https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.2.2.tar.gz \
    && tar -zvxf /tmp/esc.tgz -C /usr/share/ \
    && rm /tmp/esc.tgz \
    && mv /usr/share/elasticsearch-5.2.2 /usr/share/elasticsearch \
    && useradd --user-group --home-dir /usr/share/elasticsearch elasticsearch \
    && ln -s /usr/share/elasticsearch/bin/elasticsearch /usr/bin/elasticsearch \
    && chown -R elasticsearch:elasticsearch /usr/share/elasticsearch /usr/bin/elasticsearch

# by default elasticsearch shell is /bin/false, we need
# /bin/bash to run elasticsearch as non-root
# https://discuss.elastic.co/t/running-as-non-root-user-service-wrapper-has-changed/7863
RUN usermod -s /bin/bash -d /usr/share/elasticsearch elasticsearch

RUN /usr/share/elasticsearch/bin/elasticsearch-plugin install io.fabric8:elasticsearch-cloud-kubernetes:5.2.2 \
    && rm /usr/share/elasticsearch/plugins/discovery-kubernetes/slf4j-* \
    && curl -o /usr/share/elasticsearch/plugins/discovery-kubernetes/slf4j-log4j12-1.7.13.jar http://central.maven.org/maven2/org/slf4j/slf4j-log4j12/1.7.13/slf4j-log4j12-1.7.13.jar \
    && curl -o /usr/share/elasticsearch/plugins/discovery-kubernetes/slf4j-api-1.7.13.jar http://central.maven.org/maven2/org/slf4j/slf4j-api/1.7.13/slf4j-api-1.7.13.jar


ADD elasticsearch.yml /usr/share/elasticsearch/config/elasticsearch.yml

USER elasticsearch
