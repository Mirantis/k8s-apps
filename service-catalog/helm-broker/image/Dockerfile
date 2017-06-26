FROM quay.io/deis/base:v0.3.5

# Add user and group
RUN adduser --system \
    --shell /bin/bash \
    --disabled-password \
    --home /opt/helm-broker \
    --group \
    helm-broker

COPY . /opt/helm-broker/src/

ENV GOROOT /opt/helm-broker/go
ENV PATH $PATH:/opt/helm-broker/go/bin:/opt/discapi/sbin:/opt/helm-broker/glide
ENV GOPATH /opt/helm-broker/

RUN apt-get update \
    && apt-get install -y wget git \
    && wget https://github.com/Masterminds/glide/releases/download/0.10.2/glide-0.10.2-linux-386.tar.gz -P /opt/helm-broker \
    && tar -C /opt/helm-broker -xzf /opt/helm-broker/glide-0.10.2-linux-386.tar.gz \
    && mv /opt/helm-broker/linux-386/ /opt/helm-broker/glide/ \
    && wget https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz -P /opt/helm-broker \
    && tar -C /opt/helm-broker -xzf /opt/helm-broker/go1.8.1.linux-amd64.tar.gz \
    && cd /opt/helm-broker/src/helm-broker/ \
    && glide i --strip-vendor --strip-vcs \
    && GOOS=linux GOARCH=amd64 go build -ldflags -s -o /opt/helm-broker/sbin/helm-broker \
    && rm -rf /opt/helm-broker/src/

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
