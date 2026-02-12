//go:build !cuda && !mock_hardware
// +build !cuda,!mock_hardware

package hardware

import "fmt"

// DetectBestProvider returns ProviderCPU by default when CUDA is not enabled.
func DetectBestProvider() Provider {
	return ProviderCPU
}

// GetModelPaths returns the paths to the ONNX model files for the CPU provider.
func GetModelPaths(p Provider) (encoder, decoder, joiner, tokens string) {
	if p != ProviderCPU {
		panic(fmt.Sprintf("GetModelPaths (CPU implementation) called with non-CPU provider: %s", p))
	}
	modelDir := "models/sherpa-onnx-v1.12.24-linux-x64-shared" // Original path for CPU

	encoder = fmt.Sprintf("%s/encoder-epoch-99-avg-1.int8.onnx", modelDir)
	decoder = fmt.Sprintf("%s/decoder-epoch-99-avg-1.int8.onnx", modelDir)
	joiner = fmt.Sprintf("%s/joiner-epoch-99-avg-1.int8.onnx", modelDir)
	tokens = fmt.Sprintf("%s/tokens.txt", modelDir)
	return
}
