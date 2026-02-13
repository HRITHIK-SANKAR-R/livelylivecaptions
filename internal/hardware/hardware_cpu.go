//go:build !cuda && !mock_hardware
// +build !cuda,!mock_hardware

package hardware

import (
	"fmt"
	"os"
	"path/filepath"
)

// DetectBestProvider returns ProviderCPU by default when CUDA is not enabled.
func DetectBestProvider() Provider {
	return ProviderCPU
}

// GetModelPaths returns the paths to the ONNX model files for the CPU provider.
func GetModelPaths(p Provider) (encoder, decoder, joiner, tokens string) {
	if p != ProviderCPU {
		panic(fmt.Sprintf("GetModelPaths (CPU implementation) called with non-CPU provider: %s", p))
	}

	// Find project root by looking for go.mod
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get current working directory: %v", err))
	}
	var projectRoot string
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			projectRoot = dir
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found in any parent directory")
		}
		dir = parent
	}

	modelDir := filepath.Join(projectRoot, "models", "cpu", "sherpa-onnx-streaming-zipformer-en-20M-2023-02-17")

	encoder = filepath.Join(modelDir, "encoder-epoch-99-avg-1.int8.onnx")
	decoder = filepath.Join(modelDir, "decoder-epoch-99-avg-1.int8.onnx")
	joiner = filepath.Join(modelDir, "joiner-epoch-99-avg-1.int8.onnx")
	tokens = filepath.Join(modelDir, "tokens.txt")
	return
}
