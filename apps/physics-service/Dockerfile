FROM ghcr.io/astral-sh/uv:latest AS uv_bin

FROM python:3.12-slim
WORKDIR /app

RUN apt-get update && apt-get install -y netcat-openbsd && rm -rf /var/lib/apt/lists/*

COPY --from=uv_bin /uv /uvx /bin/

COPY pyproject.toml uv.lock ./

RUN uv sync --frozen --no-dev --no-install-project

COPY . .

CMD ["uv", "run", "python", "-m", "physics.main"]