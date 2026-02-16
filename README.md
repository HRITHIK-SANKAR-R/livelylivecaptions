# LivelyLiveCaptions

LivelyLiveCaptions is a real-time, AI-powered speech-to-text application that displays live transcriptions from your microphone input directly in your terminal. It leverages the Sherpa-ONNX library for high-quality, streaming speech recognition.

![](./livelylivecaptions.jpeg)   

**Note**: This project has been tested and verified to work optimally on Arch Linux as of February 2026 with an RTX 3050 4GB GPU. Performance may vary on other systems.

## Contributing

Feel free to create issues if you want to suggest improvements or report problems. You can submit a pull request (PR) if you find any bugs or have updates to contribute to the project. Visit the Issues tab to see existing discussions or create a new one.

## Features

- **Real-Time Transcription**: Captures audio from your microphone and provides a live feed of the transcribed text.
- **GPU Acceleration**: Supports NVIDIA GPUs (via CUDA) for faster, more efficient transcription.
- **CPU Fallback**: Automatically falls back to CPU if a compatible GPU is not found or fails to initialize.
- **Customizable UI**: Features a terminal-based UI with an audio level meter and customizable colors.
- **Flexible Configuration**: Can be configured via a `config.yaml` file, environment variables, or command-line flags.

---

## Project Setup

### 1. Core Dependencies

First, ensure you have the necessary core dependencies for your operating system.

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

---

### 1.1. Python Environment for CUDA Libraries and Model Downloads (Optional)

If you intend to use GPU acceleration or the `download.sh` script, it is highly recommended to set up a Python virtual environment to manage CUDA library dependencies and the `hf` CLI tool. This provides a self-contained environment for the necessary Python components.

#### Setup Steps:

1.  **Ensure Python and pip are installed:**
    *   **Linux (Debian/Ubuntu):**
        ```bash
        sudo apt-get update
        sudo apt-get install -y python3-venv python3-pip
        ```
    *   **Linux (Arch):**
        ```bash
        sudo pacman -Syu --noconfirm
        sudo pacman -S --noconfirm python python-pip
        ```
    *   **Windows:** Install Python from the [official website](https://www.python.org/downloads/windows/). Ensure "Add Python to PATH" is checked during installation.

2.  **Create and activate a Python virtual environment:**
    ```bash
    python3 -m venv .venv
    source .venv/bin/activate # On Windows, use .venv\Scripts\activate
    ```
    *Note: The scripts (`nemotron.sh`, `sherpa.sh`) expect the virtual environment to be named `.venv` and located in the project root.*

    *   **Using an Existing Virtual Environment:** If you already have a virtual environment with CUDA-providing Python packages installed (e.g., `torch`, `tensorflow-gpu`), you can activate it instead of creating a new one. However, ensure that the CUDA version provided by your existing venv (e.g., CUDA 12.x) is compatible with the Sherpa-ONNX binaries used by this project.

3.  **Install CUDA-providing Python packages and the `hf` CLI:**
    ```bash
    pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu121 # For CUDA 12.1 (adjust as needed for your GPU/CUDA version)
    pip install hf
    ```
    *   **Important:** The `torch` package (or `tensorflow-gpu`) bundles the necessary CUDA libraries (like `libcublasLt.so.12`) that your Go application relies on. Ensure you install the version compatible with your NVIDIA GPU and the `cuda-12.x` binaries used by Sherpa-ONNX. You can find specific installation commands for PyTorch at [https://pytorch.org/get-started/locally/](https://pytorch.org/get-started/locally/).
    *   The `hf` CLI is used by `download.sh` to fetch models from Hugging Face.

4.  **Deactivate the virtual environment when done with setup:**
    ```bash
    deactivate
    ```

---

### 2. Clone Repository

```bash
git clone https://github.com/HRITHIK-SANKAR-R/LivelyLiveCaptions.git
cd LivelyLiveCaptions
```

### 3. Download Transcription Models

The application requires models for processing. We recommend using our automated download script to get all necessary models and libraries:

```bash
chmod +x download.sh
./download.sh
```
**Note**: The `download.sh` script now utilizes the `hf` command-line interface (CLI) for downloading models from Hugging Face. Ensure `hf` is installed as per the "Dependencies" section.

The script will guide you through the download process for your operating system and model preferences. The GPU libraries are required for optimal performance and will always be downloaded.

#### Manual Setup (Alternative)

If you prefer to set up models manually, follow these steps:

1.  **Create the necessary directory structure:**
    ```
    LivelyLiveCaptions/
    ├── models/
    │   ├── nemotron/
    │   ├── sherpa/
    │   └── sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu/
    ```

2.  **Download Nemotron model files:**
    *   Download the entire repository from [https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14](https://huggingface.co/csukuangfj/sherpa-onnx-nemotron-speech-streaming-en-0.6b-2026-01-14)
    *   Place all contents directly into the `models/nemotron/` directory.

3.  **Download Sherpa June 2023 model files:**
    *   Download the archive: [https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-en-2023-06-26.tar.bz2](https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-en-2023-06-26.tar.bz2)
    *   Extract the contents. From the extracted directory (e.g., `sherpa-onnx-streaming-zipformer-en-2023-06-26/`), copy the following files into the `models/sherpa/` directory:
        *   `encoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx`
        *   `decoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx`
        *   `joiner-epoch-99-avg-1-chunk-16-left-128.int8.onnx`
        *   `tokens.txt`
        *   (Optional: `encoder-epoch-99-avg-1-chunk-16-left-128.onnx`, `decoder-epoch-99-avg-1-chunk-16-left-128.onnx`, `joiner-epoch-99-avg-1-chunk-16-left-128.onnx` if you prefer non-INT8 models)

4.  **Download GPU Libraries:**
    *   **For Linux:**
        *   Download the archive: [https://github.com/k2-fsa/sherpa-onnx/releases/download/v1.12.24/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu.tar.bz2](https://github.com/k2-fsa/sherpa-onnx/releases/download/v1.12.24/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu.tar.bz2)
        *   Extract its contents directly into the `models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu/` directory.
    *   **For Windows:**
        *   Download the archive: [https://github.com/k2-fsa/sherpa-onnx/releases/download/v1.12.24/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-win-x64-cuda.tar.bz2](https://github.com/k2-fsa/sherpa-onnx/releases/download/v1.12.24/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-win-x64-cuda.tar.bz2)
        *   Extract its contents directly into the `models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu/` directory. (Note: The target directory name is Linux-specific, but it should contain the Windows binaries).

Remember that using the `./download.sh` script is highly recommended for a simpler and faster setup.

---

## Building and Running the Application

#### Linux
The application provides convenience scripts that both build and run the application with specific model configurations:

**Option 1: Nemotron-primary model (build and run)**
This option builds `LivelyLiveCaptions_Nemotron` (using your custom Nemotron model as primary) and then runs it.
```bash
chmod +x nemotron.sh
./nemotron.sh
```

**Option 2: Sherpa-only model (build and run)**
This option builds `LivelyLiveCaptions_Sherpa` (using the preferred Sherpa June 2023 model) and then runs it.
```bash
chmod +x sherpa.sh
./sherpa.sh
```

#### Windows
You must run the `go build` command manually. For GPU support, you will need to set up the environment to correctly link the Sherpa-ONNX CUDA libraries, which can be complex.

**CPU Build (Windows):**
```bash
go build -o livelylivecaptions.exe cmd/livelylivecaptions/main.go
```

**GPU Build (Windows - Advanced):**
Building with CUDA on Windows requires manually setting paths to the Sherpa-ONNX GPU libraries. For a simpler build experience with GPU support, using **Windows Subsystem for Linux (WSL)** is highly recommended.



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
