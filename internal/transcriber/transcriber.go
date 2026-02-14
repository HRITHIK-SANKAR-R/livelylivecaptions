package transcriber

import (
	"fmt"
	"livelylivecaptions/internal/hardware"
	"livelylivecaptions/internal/logger" // Added import
	"livelylivecaptions/internal/types"
	"sync" // Import sync package

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

// OnlineRecognizer defines the subset of Sherpa-ONNX methods used by Transcriber
type OnlineRecognizer interface {
	IsReady(s *sherpa.OnlineStream) bool
	Decode(s *sherpa.OnlineStream)
	GetResult(s *sherpa.OnlineStream) *sherpa.OnlineRecognizerResult
	IsEndpoint(s *sherpa.OnlineStream) bool
	Reset(s *sherpa.OnlineStream)
}

// Transcriber handles speech recognition
type Transcriber struct {
	recognizer OnlineRecognizer
	stream     *sherpa.OnlineStream
	InputChan  chan []byte
	OutputChan chan types.TranscriptionEvent
	QuitChan   chan struct{}
	wg         sync.WaitGroup // Add WaitGroup for graceful shutdown
}

// NewTranscriber initializes the Sherpa-ONNX recognizer with hardware-specific configuration.
// It includes a panic-recovery mechanism to handle CGO errors safely.
func NewTranscriber(p hardware.Provider) (tr *Transcriber, err error) {
	// Defer a function to recover from panics, which can happen with CGO calls
	// if libraries are missing or there's a hardware mismatch.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred during recognizer initialization with provider '%s': %v", p, r)
		}
	}()

	encoderPath, decoderPath, joinerPath, tokensPath := hardware.GetModelPaths(p)

	// Configuration for the streaming zipformer model
	config := sherpa.OnlineRecognizerConfig{
		FeatConfig: sherpa.FeatureConfig{
			SampleRate: 16000,
			FeatureDim: 80,
		},
		ModelConfig: sherpa.OnlineModelConfig{
			Transducer: sherpa.OnlineTransducerModelConfig{
				Encoder: encoderPath,
				Decoder: decoderPath,
				Joiner:  joinerPath,
			},
			Tokens:     tokensPath,
			NumThreads: 1,
			Provider:   string(p),
			Debug:      0,
		},
		DecodingMethod: "modified_beam_search",
		MaxActivePaths: 4,
		EnableEndpoint: 1, // Enable endpoint detection
	}

	recognizer := sherpa.NewOnlineRecognizer(&config)
	if recognizer == nil {
		// This path is taken if Sherpa-ONNX returns nil without panicking.
		return nil, fmt.Errorf("failed to create recognizer with provider %s (returned nil)", p)
	}

	stream := sherpa.NewOnlineStream(recognizer)
	if stream == nil {
		// If stream creation fails, we must clean up the successfully created recognizer.
		sherpa.DeleteOnlineRecognizer(recognizer)
		return nil, fmt.Errorf("failed to create stream")
	}

	logger.Debug("Sherpa-ONNX transcriber resources allocated for provider: %s", p)

	return &Transcriber{
		recognizer: recognizer,
		stream:     stream,
		InputChan:  make(chan []byte, 10), // Buffered to prevent blocking audio capture
		OutputChan: make(chan types.TranscriptionEvent),
		QuitChan:   make(chan struct{}),
	}, nil
}

// BytesToSamples converts raw int16 LE bytes to float32 samples.
// It returns nil if the input is too small (less than 2 bytes).
func BytesToSamples(audioData []byte) []float32 {
	if len(audioData) < 2 {
		return nil
	}
	samples := make([]float32, len(audioData)/2)
	for i := 0; i < len(samples); i++ {
		// Little-endian conversion
		val := int16(uint16(audioData[2*i]) | uint16(audioData[2*i+1])<<8)
		samples[i] = float32(val) / 32768.0
	}
	return samples
}

// Start begins processing audio from the input channel
func (t *Transcriber) Start() {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		defer close(t.OutputChan)

		for {
			select {
			case <-t.QuitChan:
				return
			case audioData, ok := <-t.InputChan:
				if !ok {
					// InputChan was closed, exit loop
					return
				}
				samples := BytesToSamples(audioData)
				if samples == nil {
					continue
				}

				// Accept samples
				t.stream.AcceptWaveform(16000, samples)

				// Decode
				for t.recognizer.IsReady(t.stream) {
					t.recognizer.Decode(t.stream)
				}

				// Get result
				result := t.recognizer.GetResult(t.stream)

				// Only send if there's text (partial or final)
				if result != nil && len(result.Text) > 0 {
					event := types.TranscriptionEvent{
						Text: result.Text,
					}

					if t.recognizer.IsEndpoint(t.stream) {
						t.recognizer.Reset(t.stream)
						event.IsFinal = true
					}

					t.OutputChan <- event
				}
			}
		}
	}()
}

// Close releases resources and waits for the processing goroutine to finish.
func (t *Transcriber) Close() {
	// Signal the processing goroutine to stop by closing the quit channel.
	close(t.QuitChan)

	// Wait for the goroutine to finish its current work and exit.
	t.wg.Wait()

	// Now that the goroutine is stopped, it's safe to release CGO resources.
	if t.stream != nil {
		sherpa.DeleteOnlineStream(t.stream)
	}
	if t.recognizer != nil {
		// If it's the real recognizer, delete it using the Sherpa-ONNX deleter.
		// Mocks won't need this step.
		if r, ok := t.recognizer.(*sherpa.OnlineRecognizer); ok {
			sherpa.DeleteOnlineRecognizer(r)
		}
	}
}

