FROM golang:1.8.1 as builder

COPY . /go/src/
ENV PATH $PATH:/go/glide

RUN apt-get update \
    && apt-get install -y wget git \
    && wget https://github.com/Masterminds/glide/releases/download/0.10.2/glide-0.10.2-linux-386.tar.gz -P /opt/helm-broker \
    && tar -C /go -xzf /opt/helm-broker/glide-0.10.2-linux-386.tar.gz \
    && mv /go/linux-386/ /go/glide/ \
    && cd /go/src/helm-broker/ \
    && glide i --strip-vendor --strip-vcs \
    && GOOS=linux GOARCH=amd64 go build -ldflags -s -o /sbin/helm-broker

FROM quay.io/deis/base:v0.3.5

# Add user and group
RUN adduser --system \
    --shell /bin/bash \
    --disabled-password \
    --home /opt/helm-broker \
    --group \
    helm-broker

COPY --from=builder /sbin/helm-broker /opt/helm-broker/sbin/

#Fix some permission since we'll be running as a non-root user
RUN chown -R helm-broker:helm-broker /opt/helm-broker \
    && chmod +x /opt/helm-broker/sbin/*

USER helm-broker
WORKDIR /opt/helm-broker/sbin/
CMD ["./helm-broker"]

ARG VERSION
ARG BUILD_DATE
ENV VERSION $VERSION
ENV BUILD_DATE $BUILD_DATE
