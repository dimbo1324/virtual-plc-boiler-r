FROM python:3.12-slim

WORKDIR /app

COPY pyproject.toml .
RUN pip install --no-cache-dir uv && \
    uv pip install --system -e .

COPY src/ ./src/

CMD ["python", "/app/src/data_consumer/main.py"]