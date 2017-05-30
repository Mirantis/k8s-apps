FROM quay.io/deis/base:v0.3.5

# Add user and group
RUN adduser --system \
    --shell /bin/bash \
    --disabled-password \
    --home /opt/kafka-broker \
    --group \
    kafka-broker

COPY . /opt/kafka-broker/src/

ENV GOROOT /opt/kafka-broker/go
ENV PATH $PATH:/opt/kafka-broker/go/bin:/opt/discapi/sbin
ENV GOPATH /opt/kafka-broker/

RUN apt-get update \
    && apt-get install -y wget git \
    && wget https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz -P /opt/kafka-broker \
    && tar -C /opt/kafka-broker -xzf /opt/kafka-broker/go1.8.1.linux-amd64.tar.gz \
    && go get code.cloudfoundry.org/lager \
    && go get github.com/Shopify/sarama \
    && go get github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/server \
    && go get github.com/Shopify/sarama github.com/samuel/go-zookeeper/zk \
    && cd /opt/kafka-broker/src/kafka-broker/ \
    && GOOS=linux GOARCH=amd64 go build -ldflags -s -o /opt/kafka-broker/sbin/kafka-broker \
    && rm -rf /opt/kafka-broker/src/

#Fix some permission since we'll be running as a non-root user
RUN chown -R kafka-broker:kafka-broker /opt/kafka-broker \
    && chmod +x /opt/kafka-broker/sbin/*

USER kafka-broker
WORKDIR /opt/kafka-broker/sbin/
CMD ["./kafka-broker"]

ARG VERSION
ARG BUILD_DATE
ENV VERSION $VERSION
ENV BUILD_DATE $BUILD_DATE