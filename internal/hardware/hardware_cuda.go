//go:build cuda
// +build cuda

package hardware

import (
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
