FROM python:3.6

ENV DEPS_VERSION=2.5.0-beta
RUN curl https://www.dependencies.io/install.sh | bash -s -- -b /usr/local/bin $DEPS_VERSION

RUN git config --system user.email "bot@dependencies.io"
RUN git config --system user.name "Dependencies.io Bot"

ADD src/ /usr/src/app/

WORKDIR /repo

ENTRYPOINT ["python", "/usr/src/app/entrypoint.py"]
