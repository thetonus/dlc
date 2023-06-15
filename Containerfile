FROM ubuntu:22.04
WORKDIR /opt/app
COPY pyproject.toml pyproject.toml
COPY poetry.lock poetry.lock
