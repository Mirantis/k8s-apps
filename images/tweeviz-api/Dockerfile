FROM python:2.7-alpine

COPY . /tweeviz-api
RUN pip install /tweeviz-api

WORKDIR /tweeviz-api

ENTRYPOINT ["./tweeviz-api.py"]
