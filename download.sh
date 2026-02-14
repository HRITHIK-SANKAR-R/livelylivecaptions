#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

echo "==========================================="
echo "LivelyLiveCaptions Setup Script"
echo "==========================================="

# Create necessary directories
echo "Creating directory structure..."
mkdir -p models/nemotron
mkdir -p models/sherpa
mkdir -p models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu

echo "Directory structure created."

# Function to download Nemotron model
download_nemotron() {
    echo "Downloading Nemotron model files..."
    cd models/nemotron
    
    echo "Downloading encoder.onnx..."
    wget https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14/resolve/main/encoder.onnx -O encoder.onnx
    
    echo "Downloading decoder.onnx..."
    wget https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14/resolve/main/decoder.onnx -O decoder.onnx
    
    echo "Downloading joiner.onnx..."
    wget https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14/resolve/main/joiner.onnx -O joiner.onnx
    
    echo "Downloading tokens.txt..."
    wget https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14/resolve/main/tokens.txt -O tokens.txt
    
    echo "Nemotron model files downloaded."
    cd ../..
}

# Function to download Sherpa June 2023 model (from the original GPU model location)
download_sherpa() {
    echo "Downloading Sherpa June 2023 model files..."
    cd /tmp  # Use temp directory to download and extract
    
    echo "Downloading Sherpa June 2023 GPU model archive..."
    wget https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-en-2023-06-26.tar.bz2 -O sherpa-june-2023.tar.bz2
    
    # Extract to temp directory
    mkdir -p temp_extract
    tar -xjf sherpa-june-2023.tar.bz2 -C temp_extract
    extracted_dir=$(ls -1d temp_extract/*/)
    
    # Copy specific files to our sherpa directory
    cp "${extracted_dir}"encoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx ../models/sherpa/
    cp "${extracted_dir}"decoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx ../models/sherpa/
    cp "${extracted_dir}"joiner-epoch-99-avg-1-chunk-16-left-128.int8.onnx ../models/sherpa/
    cp "${extracted_dir}"tokens.txt ../models/sherpa/
    
    # Also copy other model files that might be needed
    if [ -f "${extracted_dir}encoder-epoch-99-avg-1-chunk-16-left-128.onnx" ]; then
        cp "${extracted_dir}"encoder-epoch-99-avg-1-chunk-16-left-128.onnx ../models/sherpa/
    fi
    if [ -f "${extracted_dir}decoder-epoch-99-avg-1-chunk-16-left-128.onnx" ]; then
        cp "${extracted_dir}"decoder-epoch-99-avg-1-chunk-16-left-128.onnx ../models/sherpa/
    fi
    if [ -f "${extracted_dir}joiner-epoch-99-avg-1-chunk-16-left-128.onnx" ]; then
        cp "${extracted_dir}"joiner-epoch-99-avg-1-chunk-16-left-128.onnx ../models/sherpa/
    fi
    
    # Clean up temp files
    rm -rf temp_extract sherpa-june-2023.tar.bz2
    
    echo "Sherpa June 2023 model files downloaded."
    cd ..
}

# Function to download GPU libraries for Linux
download_gpu_libs_linux() {
    echo "Downloading Sherpa-ONNX GPU libraries for Linux..."
    cd models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu
    
    # Download the GPU library archive
    echo "Downloading GPU library archive..."
    wget https://github.com/k2-fsa/sherpa-onnx/releases/download/v1.12.24/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu.tar.bz2 -O sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu.tar.bz2
    
    # Extract the archive
    echo "Extracting GPU library archive..."
    tar -xjf sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu.tar.bz2 --strip-components=1
    
    # Clean up the archive file
    rm sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu.tar.bz2
    
    echo "GPU libraries downloaded and extracted."
    cd ../..
}

# Function to download GPU libraries for Windows
download_gpu_libs_windows() {
    echo "Downloading Sherpa-ONNX GPU libraries for Windows..."
    cd models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu  # Using same directory, but with Windows binaries
    
    # Download the Windows GPU library archive
    echo "Downloading Windows GPU library archive..."
    wget https://github.com/k2-fsa/sherpa-onnx/releases/download/v1.12.24/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-win-x64-cuda.tar.bz2 -O sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-win-x64-cuda.tar.bz2
    
    # Extract the archive
    echo "Extracting Windows GPU library archive..."
    tar -xjf sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-win-x64-cuda.tar.bz2 --strip-components=1
    
    # Clean up the archive file
    rm sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-win-x64-cuda.tar.bz2
    
    echo "Windows GPU libraries downloaded and extracted."
    cd ../..
}

# Ask user for their OS
echo "Select your operating system:"
echo "1) Linux"
echo "2) Windows"
read -p "Enter choice (1 or 2): " os_choice

# Ask what model they want to download
echo "Which model would you like to download?"
echo "1) Both Nemotron and Sherpa models (recommended)"
echo "2) Only Nemotron model"
echo "3) Only Sherpa model"
read -p "Enter choice (1-3): " model_choice

# Always download GPU libraries first
echo "Downloading required GPU libraries..."
case $os_choice in
    1) download_gpu_libs_linux ;;
    2) download_gpu_libs_windows ;;
    *) echo "Invalid OS choice"; exit 1 ;;
esac

# Then download the selected models
case $model_choice in
    1)
        download_nemotron
        download_sherpa
        ;;
    2)
        download_nemotron
        ;;
    3)
        download_sherpa
        ;;
    *)
        echo "Invalid model choice"
        exit 1
        ;;
esac

echo "==========================================="
echo "Download completed!"
echo "==========================================="

echo ""
echo "Next steps:"
echo "1. Build the application using either:"
echo "   - './build_nemotron.sh' for Nemotron-primary build"
echo "   - './build_sherpa.sh' for Sherpa-only build"
echo "2. Run the application with './run.sh'"
echo ""
echo "For more information, check the README.md file."