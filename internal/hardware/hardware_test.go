//go:build !cuda && !mock_hardware
// +build !cuda,!mock_hardware

package hardware

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetModelPaths_CPU(t *testing.T) {
	// Create a temporary directory structure to simulate a project
	tmpDir, err := os.MkdirTemp("", "test-project-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a dummy go.mod to identify the project root
	if _, err := os.Create(filepath.Join(tmpDir, "go.mod")); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create the expected model directory
	modelDir := filepath.Join(tmpDir, "models", "cpu", "sherpa-onnx-streaming-zipformer-en-20M-2023-02-17")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		t.Fatalf("Failed to create model dir: %v", err)
	}

	// Create a subdirectory and change into it to test the upward search
	subDir := filepath.Join(tmpDir, "internal", "hardware")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd) // Clean up

	// Call the function
	encoder, decoder, joiner, tokens := GetModelPaths(ProviderCPU)

	// Define expected paths
	expectedEncoder := filepath.Join(modelDir, "encoder-epoch-99-avg-1.int8.onnx")
	expectedDecoder := filepath.Join(modelDir, "decoder-epoch-99-avg-1.int8.onnx")
	expectedJoiner := filepath.Join(modelDir, "joiner-epoch-99-avg-1.int8.onnx")
	expectedTokens := filepath.Join(modelDir, "tokens.txt")

	// Assert paths are correct
	if encoder != expectedEncoder {
		t.Errorf(`Encoder path mismatch.
Got:  %s
Want: %s`, encoder, expectedEncoder)
	}
	if decoder != expectedDecoder {
		t.Errorf(`Decoder path mismatch.
Got:  %s
Want: %s`, decoder, expectedDecoder)
	}
	if joiner != expectedJoiner {
		t.Errorf(`Joiner path mismatch.
Got:  %s
Want: %s`, joiner, expectedJoiner)
	}
	if tokens != expectedTokens {
		t.Errorf(`Tokens path mismatch.
Got:  %s
Want: %s`, tokens, expectedTokens)
	}
}

func TestDetectBestProvider_CPU(t *testing.T) {
	if provider := DetectBestProvider(); provider != ProviderCPU {
		t.Errorf("Expected DetectBestProvider to return ProviderCPU on a CPU build, but got %s", provider)
	}
}
