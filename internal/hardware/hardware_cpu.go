//go:build !cuda && !mock_hardware
// +build !cuda,!mock_hardware

package hardware

// DetectBestProvider returns ProviderCPU by default when CUDA is not enabled.
func DetectBestProvider() Provider {
	return ProviderCPU
}
