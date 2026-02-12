#!/bin/bash

# Base directories
PROJECT_DIR=$(pwd)
GPU_LIB_DIR="$PROJECT_DIR/models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu/lib"
VENV_LIB_DIR="/home/hrithik/Desktop/.venv/lib/python3.11/site-packages/nvidia"

# Construct LD_LIBRARY_PATH
# 1. Sherpa-ONNX GPU libs
# 2. CUDA/cuDNN libs from venv
LD_PATHS="$GPU_LIB_DIR"

# Add all nvidia library paths from the venv if they exist
if [ -d "$VENV_LIB_DIR" ]; then
    for dir in "$VENV_LIB_DIR"/*/lib; do
        if [ -d "$dir" ]; then
            LD_PATHS="$LD_PATHS:$dir"
        fi
    done
fi

export LD_LIBRARY_PATH="$LD_PATHS:$LD_LIBRARY_PATH"

echo "Running with enhanced LD_LIBRARY_PATH"
# echo "Paths: $LD_PATHS"

./livelylivecaptions "$@"
