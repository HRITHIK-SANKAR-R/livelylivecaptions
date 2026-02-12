//go:build mock_hardware
// +build mock_hardware

package hardware

import "fmt"

const (
	// ProviderMock is a dummy provider for testing
	ProviderMock Provider = "mock"
)

// DetectBestProvider returns ProviderMock when mock_hardware tag is enabled.
func DetectBestProvider() Provider {
	return ProviderMock
}

// GetModelPaths returns dummy paths for the mock provider.
func GetModelPaths(p Provider) (encoder, decoder, joiner, tokens string) {
	if p != ProviderMock {
		panic(fmt.Sprintf("GetModelPaths (Mock implementation) called with non-Mock provider: %s", p))
	}
	modelDir := "test_models/mock" // A dummy directory for mock models

	encoder = fmt.Sprintf("%s/mock_encoder.onnx", modelDir)
	decoder = fmt.Sprintf("%s/mock_decoder.onnx", modelDir)
	joiner = fmt.Sprintf("%s/mock_joiner.onnx", modelDir)
	tokens = fmt.Sprintf("%s/mock_tokens.txt", modelDir)
	return
}
