FROM mongo:3.4

COPY bootstrap.py test.py /

RUN apt update && \
    apt install -y python python-pip build-essential python-dev && \
    pip install pymongo dumb-init && \
    apt purge -y build-essential python-dev
