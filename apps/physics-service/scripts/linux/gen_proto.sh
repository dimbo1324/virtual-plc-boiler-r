#!/bin/bash
set -e

echo "Generating gRPC code..."
mkdir -p app/server/protos

uv run python -m grpc_tools.protoc \
    -I. \
    --python_out=. \
    --pyi_out=. \
    --grpc_python_out=. \
    proto/boiler.proto

echo "Done! Files generated in app/server/protos/"