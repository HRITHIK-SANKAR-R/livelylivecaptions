//go:build cuda
// +build cuda

package hardware

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// DetectBestProvider checks for CUDA-capable GPU and returns ProviderCUDA if found.
func DetectBestProvider() Provider {
	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		if _, err := exec.LookPath("nvidia-smi"); err == nil {
			cmd := exec.Command("nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
			output, err := cmd.Output()
			if err == nil && len(strings.TrimSpace(string(output))) > 0 {
				return ProviderCUDA
			}
		}
	}
	// Fallback if CUDA is not detected, though this file should only be built with cuda tag.
	// In practice, this shouldn't be reached if DetectBestProvider is called after initial capability check.
	return ProviderCPU 
}

// GetModelPaths returns the paths to the ONNX model files for the CUDA provider.
func GetModelPaths(p Provider) (encoder, decoder, joiner, tokens string) {
	if p != ProviderCUDA {
		panic(fmt.Sprintf("GetModelPaths (CUDA implementation) called with non-CUDA provider: %s", p))
	}
	modelDir := "models/sherpa-onnx-v1.12.24-cuda-12.x-cudnn-9.x-linux-x64-gpu"

	encoder = fmt.Sprintf("%s/encoder-epoch-99-avg-1.int8.onnx", modelDir)
	decoder = fmt.Sprintf("%s/decoder-epoch-99-avg-1.int8.onnx", modelDir)
	joiner = fmt.Sprintf("%s/joiner-epoch-99-avg-1.int8.onnx", modelDir)
	tokens = fmt.Sprintf("%s/tokens.txt", modelDir)
	return
}
