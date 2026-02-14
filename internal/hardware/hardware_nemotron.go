package hardware

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetNemotronModelPaths returns the paths to the ONNX model files for the Nemotron provider.
func GetNemotronModelPaths() (encoder, decoder, joiner, tokens string) {
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

	modelDir := filepath.Join(projectRoot, "models", "nemotron")

	encoder = filepath.Join(modelDir, "encoder.onnx")
	decoder = filepath.Join(modelDir, "decoder.onnx")
	joiner = filepath.Join(modelDir, "joiner.onnx")
	tokens = filepath.Join(modelDir, "tokens.txt")
	return
}