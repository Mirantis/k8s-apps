FROM rabbitmq:3.6.10 as builder

RUN apt-get update && apt-get install -y --no-install-recommends git make ca-certificates python erlang-dev erlang-src rsync zip curl
RUN git clone https://github.com/aweber/rabbitmq-autocluster.git
RUN cd rabbitmq-autocluster && make dist

FROM rabbitmq:3.6.10

ENV RABBITMQ_USE_LONGNAME=true \
    AUTOCLUSTER_LOG_LEVEL=debug \
    AUTOCLUSTER_CLEANUP=true \
    CLEANUP_INTERVAL=60 \
    CLEANUP_WARN_ONLY=false \
    AUTOCLUSTER_TYPE=k8s \
    LANG=en_US.UTF-8

COPY --from=builder /rabbitmq-autocluster/plugins/autocluster-0.6.1.ez /usr/lib/rabbitmq/lib/rabbitmq_server-3.6.10/plugins/
RUN rabbitmq-plugins enable --offline autocluster
