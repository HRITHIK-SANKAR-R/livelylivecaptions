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
	// We use the CPU model files (2023-02-17) for both modes because they
	// contain the 'attention_dims' metadata required by the v1.12 engine.
	// The larger GPU model files are missing this metadata.
	return "./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/encoder-epoch-99-avg-1.onnx",
		"./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/decoder-epoch-99-avg-1.onnx",
		"./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/joiner-epoch-99-avg-1.onnx",
		"./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/tokens.txt"
}
