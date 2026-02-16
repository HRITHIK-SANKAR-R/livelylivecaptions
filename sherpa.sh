#!/bin/bash

# Get the absolute path to the project directory
PROJECT_DIR=$(pwd)
GPU_LIB_DIR="$PROJECT_DIR/models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu/lib"
# Define the Python virtual environment's NVIDIA library directory
VENV_LIB_DIR="$PROJECT_DIR/.venv/lib/python3.12/site-packages/nvidia" # Adjust python version if needed

echo "\nBuilding Sherpa-only version with GPU support linking against: $GPU_LIB_DIR\n"

# Set CGO flags for the GPU library
export CGO_LDFLAGS="-L$GPU_LIB_DIR -lsherpa-onnx-c-api -Wl,-rpath,$GPU_LIB_DIR"

# Build the application with Sherpa-only model selection
go build -tags cuda -o LivelyLiveCaptions_Sherpa cmd/livelylivecaptions/main.go

if [ $? -eq 0 ]; then
    echo "\n✓ Sherpa-only build successful: ./LivelyLiveCaptions_Sherpa\n"
    
    # --- Start of run.sh functionality ---
    # Construct LD_LIBRARY_PATH
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

    echo "\nRunning with enhanced LD_LIBRARY_PATH via Python venv\n"
    
    ./LivelyLiveCaptions_Sherpa --model.provider="sherpa_only" "$@"
    # --- End of run.sh functionality ---
    echo "\n"
    echo "To explicitly control the model provider (e.g., force CPU or CUDA), use the --model.provider flag:"
    echo "  - To try forcing CUDA: ./LivelyLiveCaptions_Sherpa --model.provider=\"cuda\""
    echo "  - To force CPU:       ./LivelyLiveCaptions_Sherpa --model.provider=\"cpu\""
    echo "If you expect GPU to work but it's falling back to CPU, ensure 'nvidia-smi' runs correctly in your environment and your CUDA/cuDNN installations are compatible."
    echo ""
else
    echo "\n✗ Sherpa-only build failed\n"
    exit 1
fi