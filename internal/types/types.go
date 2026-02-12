package types

import "livelylivecaptions/internal/hardware"

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

// AppConfig represents the global application configuration.
// It is designed to be loaded by Viper, supporting layered configuration
// from file, environment variables, and CLI flags.
type AppConfig struct {
	Model struct {
		Provider hardware.Provider `mapstructure:"provider"` // cpu, cuda, (mock for testing)
		Path     string            `mapstructure:"path"`     // Base path for models
		Encoder  string            `mapstructure:"encoder"`
		Decoder  string            `mapstructure:"decoder"`
		Joiner   string            `mapstructure:"joiner"`
		Tokens   string            `mapstructure:"tokens"`
	} `mapstructure:"model"`
	Audio struct {
		SampleRate  int    `mapstructure:"sample_rate"`  // e.g., 16000
		DeviceID    string `mapstructure:"device_id"`    // Specific audio device ID or name
		MonitorMode bool   `mapstructure:"monitor_mode"` // Capture output audio (if supported)
	} `mapstructure:"audio"`
	Log struct {
		ToMemory bool `mapstructure:"to_memory"` // Log to in-memory ring buffer for UI display
		FilePath string `mapstructure:"file_path"` // Path to log file
		Level    string `mapstructure:"level"`     // info, debug, warn, error
	} `mapstructure:"log"`
	Debug struct {
		Enabled bool `mapstructure:"enabled"` // Enable general debug features
	} `mapstructure:"debug"`
	// Additional configuration fields can be added here
}
