# ===== STAGE =====
FROM python:3.10-slim-buster AS base

ARG ENV
ENV ENV $ENV
ARG USER=turing

RUN groupadd "${USER}" \
    && useradd -rm -d "/home/${USER}" -s /bin/bash -g "${USER}" "${USER}" \
    && mkdir -p "/home/${USER}" \
    && chown "${USER}:${USER}" "/home/${USER}" \
    && PIP_NO_CACHE_DIR=1 pip install "poetry==1.1.*" \
    && poetry config virtualenvs.create false


# ===== STAGE =====
FROM base AS requirements
WORKDIR /
COPY application/pyproject.toml .
COPY application/poetry.lock .
RUN poetry export > requirements.txt


# ===== STAGE =====
FROM sys-deps AS build
WORKDIR /
COPY --from=requirements ["/requirements.txt", "."]

RUN PIP_NO_CACHE_DIR=1 pip install -r requirements.txt \
    && rm requirements.txt


# ===== STAGE =====
FROM build AS run
COPY application /home/${USER}/app
RUN chown -R "${USER}:${USER}" "/home/${USER}/app"

USER ${USER}
WORKDIR /home/${USER}/app
CMD ["python", "-m", "webflow"]
