package hardware

// Provider represents a compute backend for the transcriber model.
type Provider string

const (
	ProviderCPU     Provider = "cpu"
	ProviderCUDA    Provider = "cuda"
	ProviderNemotron Provider = "nemotron"
	ProviderSherpaJune2023 Provider = "sherpa_june_2023"  // For the 2023-06-26 model
	// ProviderCoreML Provider = "coreml" // For future use on macOS
	ProviderMock Provider = "mock" // For testing purposes
)

// The DetectBestProvider and GetModelPaths functions are implemented in
// separate files (e.g., hardware_cpu.go, hardware_cuda.go, hardware_mock.go)
// using build tags for conditional compilation.
