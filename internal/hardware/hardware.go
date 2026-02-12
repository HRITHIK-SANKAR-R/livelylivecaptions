package hardware

import (
	"os/exec"
)

// Provider represents the compute provider for Sherpa-ONNX
type Provider string

const (
	ProviderCPU  Provider = "cpu"
	ProviderCUDA Provider = "cuda"
)

// DetectBestProvider attempts to detect if a CUDA-capable GPU is available
func DetectBestProvider() Provider {
	// Check if nvidia-smi exists and returns successfully
	cmd := exec.Command("nvidia-smi")
	if err := cmd.Run(); err != nil {
		return ProviderCPU
	}

	// Double check for CUDA availability in environment or libraries if needed,
	// but nvidia-smi is usually a good indicator for Linux users with drivers installed.
	return ProviderCUDA
}

// GetModelPaths returns the appropriate model paths based on the provider
func GetModelPaths(p Provider) (encoder, decoder, joiner, tokens string) {
	if p == ProviderCUDA {
		return "./models/gpu/sherpa-onnx-streaming-zipformer-en-2023-06-26/encoder-epoch-99-avg-1-chunk-16-left-128.onnx",
			"./models/gpu/sherpa-onnx-streaming-zipformer-en-2023-06-26/decoder-epoch-99-avg-1-chunk-16-left-128.onnx",
			"./models/gpu/sherpa-onnx-streaming-zipformer-en-2023-06-26/joiner-epoch-99-avg-1-chunk-16-left-128.onnx",
			"./models/gpu/sherpa-onnx-streaming-zipformer-en-2023-06-26/tokens.txt"
	}

	// Default to CPU model
	return "./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/encoder-epoch-99-avg-1.onnx",
		"./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/decoder-epoch-99-avg-1.onnx",
		"./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/joiner-epoch-99-avg-1.onnx",
		"./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/tokens.txt"
}
