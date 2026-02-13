#!/bin/bash

# Get the absolute path to the project directory
PROJECT_DIR=$(pwd)
GPU_LIB_DIR="$PROJECT_DIR/models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu/lib"

echo "Building with GPU support linking against: $GPU_LIB_DIR"

# Set CGO flags for the GPU library
export CGO_LDFLAGS="-L$GPU_LIB_DIR -lsherpa-onnx-c-api -Wl,-rpath,$GPU_LIB_DIR"

# Build the application
go build -tags cuda -o livelylivecaptions cmd/livelylivecaptions/main.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful: ./livelylivecaptions"
else
    echo "✗ Build failed"
    exit 1
fi
