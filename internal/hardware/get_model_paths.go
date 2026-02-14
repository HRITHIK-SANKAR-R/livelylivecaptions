package hardware

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetModelPaths returns the paths to the ONNX model files for the specified provider.
func GetModelPaths(p Provider) (encoder, decoder, joiner, tokens string) {
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

	switch p {
	case ProviderCPU:
		modelDir := filepath.Join(projectRoot, "models", "cpu", "sherpa-onnx-streaming-zipformer-en-20M-2023-02-17")
		encoder = filepath.Join(modelDir, "encoder-epoch-99-avg-1.int8.onnx")
		decoder = filepath.Join(modelDir, "decoder-epoch-99-avg-1.int8.onnx")
		joiner = filepath.Join(modelDir, "joiner-epoch-99-avg-1.int8.onnx")
		tokens = filepath.Join(modelDir, "tokens.txt")
	case ProviderCUDA:
		modelDir := filepath.Join(projectRoot, "models", "gpu", "sherpa-onnx-streaming-zipformer-en-2023-06-26")
		encoder = filepath.Join(modelDir, "encoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx")
		decoder = filepath.Join(modelDir, "decoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx")
		joiner = filepath.Join(modelDir, "joiner-epoch-99-avg-1-chunk-16-left-128.int8.onnx")
		tokens = filepath.Join(modelDir, "tokens.txt")
	case ProviderNemotron:
		modelDir := filepath.Join(projectRoot, "models", "nemotron")
		encoder = filepath.Join(modelDir, "encoder.onnx")
		decoder = filepath.Join(modelDir, "decoder.onnx")
		joiner = filepath.Join(modelDir, "joiner.onnx")
		tokens = filepath.Join(modelDir, "tokens.txt")
	default:
		panic(fmt.Sprintf("unsupported provider: %s", p))
	}

	return
}