package types

// TranscriptionEvent represents a single update from the transcriber
type TranscriptionEvent struct {
	Text       string
	IsFinal    bool
	Confidence float64
}

// AudioLevelMsg carries the RMS value for UI updates
type AudioLevelMsg float64

// AudioDevice defines the interface for interacting with audio hardware
type AudioDevice interface {
	Name() string
	ID() interface{}
	// Capture starts recording and sends raw bytes to audioChan.
	// It should also calculate and send RMS levels to levelChan if provided.
	Capture(audioChan chan<- []byte, levelChan chan<- AudioLevelMsg, quitChan <-chan struct{}) error
}

// AudioProvider defines the interface for audio capture engines (e.g., malgo)
type AudioProvider interface {
	GetDevices() ([]AudioDevice, error)
}

// AppConfig represents the global configuration state managed by Viper
type AppConfig struct {
	ModelPath    string
	UseGPU       bool
	SampleRate   int
	SelectedMic  string
	LogToMemory  bool
	DebugEnabled bool
}
