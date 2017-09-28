FROM quay.io/deis/base:v0.3.7

RUN curl -s -L https://github.com/dcos/demos/releases/download/0.2/traffic-fetcher -o /traffic-fetcher && \
    chmod +x /traffic-fetcher

ENTRYPOINT ["/traffic-fetcher"]