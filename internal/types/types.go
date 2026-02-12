package types

// TranscriptionEvent represents a single update from the transcriber
type TranscriptionEvent struct {
	Text       string
	IsFinal    bool
	Confidence float64
}

// AudioLevelMsg carries the RMS value for UI updates
type AudioLevelMsg float64
