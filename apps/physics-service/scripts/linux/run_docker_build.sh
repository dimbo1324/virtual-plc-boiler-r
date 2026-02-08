#!/bin/bash
set -e

IMAGE_NAME="physics-service:v1"
CONTAINER_NAME="boiler-sim"

echo "Building Docker Image..."
docker build -t $IMAGE_NAME .

echo "Stopping old container if exists..."
docker stop $CONTAINER_NAME 2>/dev/null || true
docker rm $CONTAINER_NAME 2>/dev/null || true

echo "Running new container..."
docker run -d \
    -p 50051:50051 \
    --name $CONTAINER_NAME \
    $IMAGE_NAME

echo "Container started. Logs:"
docker logs -f $CONTAINER_NAME