FROM python:2.7-alpine

COPY setup.py tweepub.py /tweepub/
RUN pip install /tweepub

ENTRYPOINT ["tweepub"]
