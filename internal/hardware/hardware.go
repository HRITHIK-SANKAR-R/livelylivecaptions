package hardware

// Provider represents a compute backend for the transcriber model.
type Provider string

const (
	ProviderCPU  Provider = "cpu"
	ProviderCUDA Provider = "cuda"
	// ProviderCoreML Provider = "coreml" // For future use on macOS
	// ProviderMock Provider = "mock" // For testing purposes
)

// DetectBestProvider checks for available hardware and returns the best provider.
// This function is implemented in separate files (e.g., hardware_cpu.go, hardware_cuda.go, hardware_mock.go)
// using build tags for conditional compilation.
func DetectBestProvider() Provider {
	panic("DetectBestProvider should be implemented via build tags")
}

// GetModelPaths returns the paths to the ONNX model files based on the provider.
// This function is implemented in separate files (e.g., hardware_cpu.go, hardware_cuda.go, hardware_mock.go)
// using build tags for conditional compilation.
func GetModelPaths(p Provider) (encoder, decoder, joiner, tokens string) {
	panic("GetModelPaths should be implemented via build tags")
}
