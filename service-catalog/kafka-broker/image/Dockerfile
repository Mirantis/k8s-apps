FROM golang:1.8.1 as builder

COPY . /go/src/

RUN apt-get update \
    && apt-get install -y wget git \
    && go get code.cloudfoundry.org/lager \
    && go get github.com/Shopify/sarama \
    && go get github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/server \
    && go get github.com/Shopify/sarama github.com/samuel/go-zookeeper/zk \
    && cd /go/src/kafka-broker/ \
    && GOOS=linux GOARCH=amd64 go build -ldflags -s -o /sbin/kafka-broker

FROM quay.io/deis/base:v0.3.5

# Add user and group
RUN adduser --system \
    --shell /bin/bash \
    --disabled-password \
    --home /opt/kafka-broker \
    --group \
    kafka-broker

COPY --from=builder /sbin/kafka-broker /opt/kafka-broker/sbin/

#Fix some permission since we'll be running as a non-root user
RUN chown -R kafka-broker:kafka-broker /opt/kafka-broker \
    && chmod +x /opt/kafka-broker/sbin/*

USER kafka-broker
WORKDIR /opt/kafka-broker/sbin/
CMD ["./kafka-broker"]

ARG VERSION
ARG BUILD_DATE
ENV PATH $PATH:/opt/discapi/sbin
ENV VERSION $VERSION
ENV BUILD_DATE $BUILD_DATE
