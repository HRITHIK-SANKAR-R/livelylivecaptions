#!/bin/bash

# Create the directory for the Nemotron model
mkdir -p models/nemotron

# Navigate to the directory
cd models/nemotron

echo "Downloading Nemotron speech model files..."

# Download the model files using wget
# These are the expected filenames based on the Hugging Face repository
wget https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14/resolve/main/encoder.onnx -O encoder.onnx
wget https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14/resolve/main/decoder.onnx -O decoder.onnx
wget https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14/resolve/main/joiner.onnx -O joiner.onnx
wget https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14/resolve/main/tokens.txt -O tokens.txt

echo "Nemotron model files downloaded successfully!"
echo "Files in models/nemotron/:"
ls -lh