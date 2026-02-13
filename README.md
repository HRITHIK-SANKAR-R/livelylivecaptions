# LivelyLiveCaptions

LivelyLiveCaptions is a real-time, AI-powered speech-to-text application that displays live transcriptions from your microphone input directly in your terminal. It leverages the Sherpa-ONNX library for high-quality, streaming speech recognition.

## Features

- **Real-Time Transcription**: Captures audio from your microphone and provides a live feed of the transcribed text.
- **GPU Acceleration**: Supports NVIDIA GPUs (via CUDA) for faster, more efficient transcription.
- **CPU Fallback**: Automatically falls back to CPU if a compatible GPU is not found or fails to initialize.
- **Customizable UI**: Features a terminal-based UI with an audio level meter and customizable colors.
- **Flexible Configuration**: Can be configured via a `config.yaml` file, environment variables, or command-line flags.

---

## Project Setup

### 1. Dependencies

First, ensure you have the necessary dependencies for your operating system.

#### Linux (Debian/Ubuntu)
```bash
sudo apt-get update
sudo apt-get install -y portaudio19-dev wget tar git
```

#### Linux (Arch)
```bash
sudo pacman -Syu --noconfirm
sudo pacman -S --noconfirm portaudio wget tar git
```

#### Windows
- **Go**: Install Go for Windows from the [official website](https://golang.org/dl/).
- **Git**: Install Git for Windows from the [official website](https://git-scm.com/download/win). A shell environment like Git Bash is recommended.
- **PortAudio**: There is no official package manager for PortAudio on Windows. The `go get` process may handle this, but manual setup might be required if you encounter errors.
- **NVIDIA GPU Drivers & CUDA**: Ensure you have the latest NVIDIA drivers and the CUDA Toolkit installed if you plan to use GPU acceleration.

### 2. Clone Repository

```bash
git clone https://github.com/your-username/LivelyLiveCaptions.git
cd LivelyLiveCaptions
```
*(Replace `your-username` with the actual repository path)*

### 3. Download Transcription Models

The application requires different models for CPU and GPU processing. The following directory structure must be created:

```
LivelyLiveCaptions/
└── models/
    ├── cpu/
    │   └── sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/
    │       ├── ... model files ...
    └── gpu/
        └── sherpa-onnx-streaming-zipformer-en-2023-06-26/
            ├── ... model files ...
```

Create these directories:
```bash
mkdir -p models/cpu models/gpu
```

#### CPU Model

Download and extract the CPU model into the `models/cpu/` directory.

```bash
# Enter the cpu model directory
cd models/cpu

# Download the model
wget https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2

# Extract the model
tar -xvf sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2

# Clean up the archive
rm sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2

# Go back to the project root
cd ../..
```

#### GPU Model (Optional)

For GPU support, download the streaming zipformer model. You can download it manually from the [Hugging Face repository](https://huggingface.co/k2-fsa/sherpa-onnx-streaming-zipformer-en-2023-06-26/tree/main) or use `git` to clone it directly into the correct directory.

```bash
git clone https://huggingface.co/k2-fsa/sherpa-onnx-streaming-zipformer-en-2023-06-26 models/gpu/sherpa-onnx-streaming-zipformer-en-2023-06-26
```

---

## Building the Application

#### Linux
A helper script `build.sh` is provided. It sets the necessary CGO flags to link against the Sherpa-ONNX libraries.

```bash
chmod +x build.sh
./build.sh
```

#### Windows
You must run the `go build` command manually. For GPU support, you will need to set up the environment to correctly link the Sherpa-ONNX CUDA libraries, which can be complex.

**CPU Build (Windows):**
```bash
go build -o livelylivecaptions.exe cmd/livelylivecaptions/main.go
```

**GPU Build (Windows - Advanced):**
Building with CUDA on Windows requires manually setting paths to the Sherpa-ONNX GPU libraries. For a simpler build experience with GPU support, using **Windows Subsystem for Linux (WSL)** is highly recommended.

---

## Running the Application

#### Linux
It is **crucial** to use the `run.sh` script to start the application. This script sets the `LD_LIBRARY_PATH` so the program can find the required GPU libraries.

```bash
chmod +x run.sh
./run.sh
```

#### Windows
You must set the `PATH` environment variable to include the directory containing the required `.dll` files from the Sherpa-ONNX Windows package.

```powershell
# Example: Assuming the required DLLs are in a 'lib' folder
$env:PATH = "C:\\path\\to\\sherpa-onnx-libs;" + $env:PATH
./livelylivecaptions.exe
```

Upon first run, you will be prompted to select an audio input device from a list.

## Configuration

You can configure the application in three ways (from lowest to highest priority):

1.  **`config.yaml` file:** Create a `config.yaml` file in the root directory.
    ```yaml
    model:
      provider: "cuda" # "cuda" or "cpu". Auto-detects if empty.
    audio:
      device_id: "default" # Name or ID of your audio device.
    ```
2.  **Environment Variables:**
    ```bash
    # Linux
    export LIVELY_MODEL_PROVIDER="cpu"
    # Windows (PowerShell)
    $env:LIVELY_MODEL_PROVIDER="cpu"
    ```
3.  **Command-Line Flags:**
    ```bash
    # Linux
    ./run.sh --model.provider="cpu"
    # Windows
    ./livelylivecaptions.exe --model.provider="cpu"
    ```
