FROM python:2.7-jessie

COPY . /tweeviz-ui
RUN pip install /tweeviz-ui

WORKDIR /tweeviz-ui

ENTRYPOINT ["./tweeviz_ui.py"]
