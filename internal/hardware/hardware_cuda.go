//go:build cuda
// +build cuda

package hardware

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	modelDir := filepath.Join(projectRoot, "models", "gpu", "sherpa-onnx-streaming-zipformer-en-2023-06-26")

	encoder = filepath.Join(modelDir, "encoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx")
	decoder = filepath.Join(modelDir, "decoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx")
	joiner = filepath.Join(modelDir, "joiner-epoch-99-avg-1-chunk-16-left-128.int8.onnx")
	tokens = filepath.Join(modelDir, "tokens.txt")
	return
}
