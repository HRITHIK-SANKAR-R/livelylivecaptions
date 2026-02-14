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

The application requires models for processing. We recommend using our automated download script to get all necessary models and libraries:

```bash
chmod +x download.sh
./download.sh
```

The script will guide you through the download process for your operating system and model preferences. The GPU libraries are required for optimal performance and will always be downloaded.

#### Manual Setup (Alternative)

If you prefer to set up models manually, the following directory structure is used:

```
LivelyLiveCaptions/
├── models/
│   ├── nemotron/          # Nemotron model files
│   ├── sherpa/            # Sherpa June 2023 model files  
│   └── sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu/  # GPU libraries (required)
```

The automated script handles all necessary downloads and extraction for you.

---

## Building the Application

#### Linux
The application now supports two different build configurations:

**Option 1: Nemotron-primary build (uses your custom model as primary)**
```bash
chmod +x build_nemotron.sh
./build_nemotron.sh
```
This creates `livelylivecaptions-nemotron` with the following hierarchy:
1. Nemotron model (primary) - located in `models/nemotron/`
2. Sherpa June 2023 model (fallback) - located in `models/sherpa/`

**Option 2: Sherpa-only build (using the preferred June 2023 model)**
```bash
chmod +x build_sherpa.sh
./build_sherpa.sh
```
This creates `livelylivecaptions-sherpa` with the following hierarchy:
1. Sherpa June 2023 model (primary) - located in `models/sherpa/`
2. Sherpa June 2023 model (fallback) - same model used for both GPU and CPU

**Legacy build (original behavior)**
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
      provider: "" # "", "nemotron_only", "sherpa_only", "cuda", or "cpu". Auto-detects if empty.
    audio:
      device_id: "default" # Name or ID of your audio device.
    ```
    
The `provider` field now accepts additional values:
- `""` (empty) or `"nemotron_only"`: Use Nemotron model as primary with Sherpa fallbacks
- `"sherpa_only"`: Use Sherpa models only (GPU primary, CPU fallback)
- `"cuda"` or `"cpu"`: Traditional behavior with hardware-specific models

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
