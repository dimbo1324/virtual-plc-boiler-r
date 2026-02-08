FROM python:3.12-slim-bookworm

COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/uv

WORKDIR /app

COPY pyproject.toml uv.lock ./

RUN uv sync --frozen --no-cache

COPY . .

ENV PYTHONPATH=/app

EXPOSE 50051

CMD ["uv", "run", "cmd/main.py"]