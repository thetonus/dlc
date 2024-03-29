##################
# Mamba Layer
##################
FROM ubuntu:22.04 as mambaforge
ARG MAMBA_VERSION=23.1.0-1

RUN set -ex && \
  apt-get update && \
  apt-get install -y --no-install-recommends \
  build-essential \
  cmake \
  unzip \
  yasm \
  pkg-config \
  libtbb2 \
  libtbb-dev \
  curl \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/* \
  && curl -sSL -o /tmp/mambaforge.sh https://github.com/conda-forge/miniforge/releases/download/${MAMBA_VERSION}/Mambaforge-${MAMBA_VERSION}-Linux-x86_64.sh \
  && chmod +x /tmp/mambaforge.sh
# Install conda
RUN /tmp/mambaforge.sh -b -p /opt/conda

##################
# Base Layer - GPU
##################
FROM nvidia/cuda:11.8.0-runtime-ubuntu22.04 as base


# Set Env Variables for the images
ENV DEBIAN_FRONTEND=noninteractive \
  LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/usr/local/lib" \
  LD_LIBRARY_PATH="/opt/conda/lib:${LD_LIBRARY_PATH}" \
  PATH="/opt/conda/bin:$PATH" \
  CONDA_AUTO_UPDATE_CONDA=false \
  # Set MKL_THREADING_LAYER=GNU to prevent issues between torch and numpy/mkl
  MKL_THREADING_LAYER=GNU \
  # See http://bugs.python.org/issue19846
  LANG=C.UTF-8  \
  # Recommended for running python apps in containers
  PYTHONUNBUFFERED=1

RUN set -ex && \
  apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Grab Conda
COPY --from=mambaforge /opt/conda /opt/conda

# Ensure conda environment is always activated
RUN ln -s /opt/conda/etc/profile.d/conda.sh /etc/profile.d/conda.sh \
  && echo ". /opt/conda/etc/profile.d/conda.sh; conda activate base" >> /etc/skel/.bashrc \
  && echo ". /opt/conda/etc/profile.d/conda.sh; conda activate base" >> ~/.bashrc# Install Python and base dependencies

# Install Python and base dependencies
RUN mamba install -y -n base -c conda-forge python=3.10.11 \
  cython \
  typing \
  libgcc \
  # Instances are Intel-CPU machines, so add extra optimizations
  mkl \
  mkl-include \
  # Below 2 are included in miniconda base, but not mamba so need to install
  conda-content-trust \
  charset-normalizer \
  # Clean conda dependencies
  && mamba clean -ya

# Install Poetry
ENV POETRY_VIRTUALENVS_IN_PROJECT=true \
  POETRY_VIRTUALENVS_CREATE=true
RUN python --version \
  && pip install poetry==1.5.1 \
  && poetry --version


##################
# Base Layer
##################
FROM base as builder
ARG ARTIFACTORY_PYPI_USERNAME
ARG ARTIFACTORY_PYPI_PASSWORD

ENV POETRY_HTTP_BASIC_SHIPT_RESOLVE_USERNAME=$ARTIFACTORY_PYPI_USERNAME \
  POETRY_HTTP_BASIC_SHIPT_RESOLVE_PASSWORD=$ARTIFACTORY_PYPI_PASSWORD \
  PROJECT_HOME=/opt/receipt_ocr

WORKDIR ${PROJECT_HOME}

COPY pyproject.toml pyproject.toml
COPY poetry.lock poetry.lock

RUN poetry install --only main --no-root
COPY resources resources
COPY app app

RUN poetry install --only main


##################
# Dev Layer
##################
FROM builder as dev
RUN poetry install


##################
# Prod Layer
##################
FROM base as prod

ENV PROMETHEUS_MULTIPROC_DIR=/tmp \
  PROJECT_HOME=/opt/receipt_ocr \
  PROJECT_PYTHONPATH=${PROJECT_HOME}/.venv/lib/python3.10/site-packages \
  PROJECT_VENV_BIN=${PROJECT_HOME}/.venv/bin \
  # prefix the venv bin to the path
  PATH=${PROJECT_VENV_BIN}${PATH:+":$PATH"} \
  # allow invoking without poetry, e.g. can use simply `python -m` instead of `poetry run python -m`
  # this saves about a second of start up time in addition to being simpler to inspect later
  # e.g. /opt/my-project-name/.venv/lib/python3.9/site-packages
  PYTHONPATH=${PROJECT_PYTHONPATH}${PYTHONPATH:+":$PYTHONPATH"}

WORKDIR ${PROJECT_HOME}

COPY --from=builder ${PROJECT_HOME}/resources ${PROJECT_HOME}/resources
COPY --from=builder ${PROJECT_HOME}/app ${PROJECT_HOME}/app

