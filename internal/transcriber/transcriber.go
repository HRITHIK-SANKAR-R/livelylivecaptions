package transcriber

import (
	"fmt"
	"livelylivecaptions/internal/types"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

// Transcriber handles speech recognition
type Transcriber struct {
	recognizer *sherpa.OnlineRecognizer
	stream     *sherpa.OnlineStream
	InputChan  chan []byte
	OutputChan chan types.TranscriptionEvent
	QuitChan   chan struct{}
}

// NewTranscriber initializes the Sherpa-ONNX recognizer
func NewTranscriber() (*Transcriber, error) {
	// Configuration for the streaming zipformer model
	// Updated to match the nested structure of sherpa-onnx-go API
	config := sherpa.OnlineRecognizerConfig{
		FeatConfig: sherpa.FeatureConfig{
			SampleRate: 16000,
			FeatureDim: 80,
		},
		ModelConfig: sherpa.OnlineModelConfig{
			Transducer: sherpa.OnlineTransducerModelConfig{
				Encoder: "./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/encoder-epoch-99-avg-1.onnx",
				Decoder: "./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/decoder-epoch-99-avg-1.onnx",
				Joiner:  "./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/joiner-epoch-99-avg-1.onnx",
			},
			Tokens:     "./models/cpu/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/tokens.txt",
			NumThreads: 1,
			Debug:      0,
			ModelType:  "zipformer",
		},
		DecodingMethod: "greedy_search",
		MaxActivePaths: 4,
		EnableEndpoint: 1, // Enable endpoint detection
	}

	recognizer := sherpa.NewOnlineRecognizer(&config)
	if recognizer == nil {
		return nil, fmt.Errorf("failed to create recognizer")
	}

	stream := sherpa.NewOnlineStream(recognizer)
	if stream == nil {
		return nil, fmt.Errorf("failed to create stream")
	}

	return &Transcriber{
		recognizer: recognizer,
		stream:     stream,
		InputChan:  make(chan []byte, 10), // Buffered to prevent blocking audio capture
		OutputChan: make(chan types.TranscriptionEvent),
		QuitChan:   make(chan struct{}),
	}, nil
}

// Start begins processing audio from the input channel
func (t *Transcriber) Start() {
	go func() {
		defer close(t.OutputChan)
		
		// Expected sample rate (16kHz) and format (float32 required by Sherpa)
		// Audio capture is providing S16LE (int16), so we need to convert.
		
		for {
			select {
			case <-t.QuitChan:
				return
			case audioData := <-t.InputChan:
				// Skip empty buffers
				if len(audioData) < 2 {
					continue
				}
				// Convert []byte (int16 LE) to []float32
				samples := make([]float32, len(audioData)/2)
				for i := 0; i < len(samples); i++ {
					// Little-endian conversion
					val := int16(uint16(audioData[2*i]) | uint16(audioData[2*i+1])<<8)
					samples[i] = float32(val) / 32768.0
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
				if len(result.Text) > 0 {
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

// Close releases resources
func (t *Transcriber) Close() {
	if t.stream != nil {
		sherpa.DeleteOnlineStream(t.stream)
	}
	if t.recognizer != nil {
		sherpa.DeleteOnlineRecognizer(t.recognizer)
	}
}
