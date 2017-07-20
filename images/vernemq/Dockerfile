FROM erlio/docker-vernemq:latest

RUN apt update && apt install dnsutils -y

COPY run.sh /run.sh

CMD ["/run.sh"]
