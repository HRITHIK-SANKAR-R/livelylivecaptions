package transcriber_test

import (
	"context"
	// "fmt"
	"livelylivecaptions/internal/audio"
	"livelylivecaptions/internal/hardware"
	"livelylivecaptions/internal/transcriber"
	// "livelylivecaptions/internal/types"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/codycollier/wer"
)

// Helper function to get the root directory of the project
func getProjectRoot() string {
	// Assuming tests are run from the project root or from a subdirectory like internal/transcriber
	// This finds the go.mod file to determine the root
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found in any parent directory")
		}
		dir = parent
	}
}

func TestGoldenFileTranscription(t *testing.T) {
	projectRoot := getProjectRoot()
	goldenAudioPath := filepath.Join(projectRoot, audio.DefaultTestWavPath) // Using the sine wave audio
	goldenTranscriptPath := filepath.Join(projectRoot, "test_assets", "golden_transcript.txt")

	// 1. Load Golden Master Audio
	mockAudioDevice, err := audio.NewMockAudioDevice("Mock Test Device", "mock-01", goldenAudioPath)
	if err != nil {
		t.Fatalf("Failed to create mock audio device: %v", err)
	}

	// 2. Initialize Transcriber (using CPU provider for testing)
	tr, err := transcriber.NewTranscriber(hardware.ProviderCPU)
	if err != nil {
		t.Fatalf("Failed to initialize transcriber: %v", err)
	}
	defer tr.Close()

	// Context for managing the test goroutines
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10-second timeout for the test
	defer cancel()

	var transcribedTextBuilder strings.Builder
	var wg sync.WaitGroup

	// Goroutine to start mock audio device and feed data to transcriber
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(tr.InputChan) // Ensure input channel is closed when audio feeding stops

		if err := mockAudioDevice.Start(); err != nil {
			t.Errorf("Failed to start mock audio device: %v", err)
			return
		}
		defer mockAudioDevice.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				audioData, err := mockAudioDevice.Read()
				if err != nil {
					t.Errorf("Error reading from mock audio device: %v", err)
					return
				}
				if audioData == nil {
					// No data yet, wait a bit
					time.Sleep(10 * time.Millisecond)
					continue
				}
				select {
				case tr.InputChan <- audioData:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	// Goroutine to start transcriber and collect results
	wg.Add(1)
	go func() {
		defer wg.Done()
		tr.Start() // Start transcriber's internal processing loop
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-tr.OutputChan:
				if !ok { // Channel closed
					return
				}
				if event.IsFinal {
					transcribedTextBuilder.WriteString(event.Text + " ")
				}
			}
		}
	}()

	// Wait for goroutines to finish or context to be cancelled
	// Give some time for transcription to happen
	go func() {
		wg.Wait()
		cancel() // Signal test completion after all relevant goroutines are done
	}()

	<-ctx.Done() // Block until test finishes or times out

	finalTranscribedText := strings.TrimSpace(transcribedTextBuilder.String())

	// 3. Load Golden Transcript
	expectedTranscriptBytes, err := os.ReadFile(goldenTranscriptPath)
	if err != nil {
		t.Fatalf("Failed to read golden transcript file: %v", err)
	}
	expectedTranscript := strings.TrimSpace(string(expectedTranscriptBytes))

	// 4. Calculate WER
	// Tokenize words for WER calculation
	actualWords := strings.Fields(finalTranscribedText)
	expectedWords := strings.Fields(expectedTranscript)

	// Note: wer.CalculateWER returns 0 for identical empty slices, which is fine for our empty transcript test.
	errorRate, _ := wer.WER(expectedWords, actualWords)

	const maxAllowedWER = 0.5 // Adjust this threshold based on model performance
	if errorRate > maxAllowedWER {
		t.Errorf("Word Error Rate too high! Expected WER <= %.2f, Got %.2fActual: '%s'Expected: '%s'",maxAllowedWER, errorRate, finalTranscribedText, expectedTranscript)
	} else {
		t.Logf("Word Error Rate (WER) %.2f is within acceptable limits (<= %.2f)", errorRate, maxAllowedWER)
		t.Logf("Actual transcription: '%s'", finalTranscribedText)
		t.Logf("Expected transcription: '%s'", expectedTranscript)
	}
}

// Ensure the Transcriber is closed correctly when the main test function exits
// This defer ensures resources are cleaned up even if there are early exits due to fatal errors.
// func TestMain(m *testing.M) {
//     // setup here if needed
//     code := m.Run()
//     // teardown here if needed
//     os.Exit(code)
// }
