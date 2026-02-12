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
	Start() error
	Read() ([]byte, error)
	Close() error
}

// AudioProvider defines the interface for audio capture engines (e.g., malgo)
type AudioProvider interface {
	GetDevices() ([]AudioDevice, error)
}
