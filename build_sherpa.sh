#!/bin/bash

# Get the absolute path to the project directory
PROJECT_DIR=$(pwd)
GPU_LIB_DIR="$PROJECT_DIR/models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu/lib"

echo "Building Sherpa-only version with GPU support linking against: $GPU_LIB_DIR"

# Set CGO flags for the GPU library
export CGO_LDFLAGS="-L$GPU_LIB_DIR -lsherpa-onnx-c-api -Wl,-rpath,$GPU_LIB_DIR"

# Build the application with Sherpa-only model selection
go build -tags cuda -o LivelyLiveCaptions_Sherpa cmd/livelylivecaptions/main.go

if [ $? -eq 0 ]; then
    echo "✓ Sherpa-only build successful: ./LivelyLiveCaptions_Sherpa"
else
    echo "✗ Sherpa-only build failed"
    exit 1
fi