FROM redis:4.0-alpine

COPY redis-cluster.sh /run.sh
COPY kill.sh /kill.sh
ADD https://raw.githubusercontent.com/antirez/redis/4.0.1/src/redis-trib.rb /usr/local/bin/
RUN apk add --update --no-cache bash perl ruby ruby-rdoc ruby-irb && gem install redis
CMD ["/run.sh"]
