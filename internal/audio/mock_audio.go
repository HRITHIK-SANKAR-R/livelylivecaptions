package audio

import (
	// "bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"

	// "livelylivecaptions/internal/types"
)

const (
	// DefaultTestWavPath is the expected path for a test WAV file
	DefaultTestWavPath = "test_assets/speech_16k.wav"
	// MockChunkSize is the size of audio data returned by Read() in bytes (e.g., 100ms of 16kHz 16-bit mono audio)
	MockChunkSize = 16000 * 2 / 10 // 1600 samples * 2 bytes/sample = 3200 bytes
)

// MockAudioDevice implements the types.AudioDevice interface for testing purposes.
// It streams pre-loaded audio data (e.g., from a WAV file) in a loop.
type MockAudioDevice struct {
	name        string
	id          string
	audioData   []byte // The full audio data from the WAV file
	currentPos  int
	mu          sync.Mutex
	closed      bool
	sampleRate  int
	numChannels int
	bitDepth    int
}

// NewMockAudioDevice creates a new MockAudioDevice instance.
// It loads audio from the specified WAV file.
func NewMockAudioDevice(name, id, wavFilePath string) (*MockAudioDevice, error) {
	data, err := os.ReadFile(wavFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read WAV file %s: %w", wavFilePath, err)
	}

	// Basic WAV header parsing to get relevant audio properties.
	// This is a simplified parser and might not handle all WAV formats.
	if len(data) < 44 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" || string(data[12:16]) != "fmt " {
		return nil, fmt.Errorf("invalid WAV file format: %s", wavFilePath)
	}

	// Audio Format (bytes 20-21)
	audioFormat := binary.LittleEndian.Uint16(data[20:22])
	if audioFormat != 1 { // 1 means PCM
		return nil, fmt.Errorf("unsupported WAV audio format (only PCM is supported): %s", wavFilePath)
	}

	// Number of Channels (bytes 22-23)
	numChannels := binary.LittleEndian.Uint16(data[22:24])
	if numChannels != 1 { // Only mono for now
		return nil, fmt.Errorf("unsupported number of channels (only mono is supported): %s", wavFilePath)
	}

	// Sample Rate (bytes 24-27)
	sampleRate := binary.LittleEndian.Uint32(data[24:28])

	// Bits Per Sample (bytes 34-35)
	bitsPerSample := binary.LittleEndian.Uint16(data[34:36])
	if bitsPerSample != 16 { // Only 16-bit for now
		return nil, fmt.Errorf("unsupported bit depth (only 16-bit is supported): %s", wavFilePath)
	}

	// Find the "data" subchunk
	dataChunkOffset := -1
	dataChunkSize := 0
	for i := 36; i < len(data)-8; {
		chunkID := string(data[i : i+4])
		chunkLen := int(binary.LittleEndian.Uint32(data[i+4 : i+8]))
		if chunkID == "data" {
			dataChunkOffset = i + 8
			dataChunkSize = chunkLen
			break
		}
		i += 8 + chunkLen
	}

	if dataChunkOffset == -1 {
		return nil, fmt.Errorf("WAV file does not contain a 'data' chunk: %s", wavFilePath)
	}

	audioContent := data[dataChunkOffset : dataChunkOffset+dataChunkSize]

	return &MockAudioDevice{
		name:        name,
		id:          id,
		audioData:   audioContent,
		sampleRate:  int(sampleRate),
		numChannels: int(numChannels),
		bitDepth:    int(bitsPerSample),
	}, nil
}

// Name returns the name of the mock device.
func (m *MockAudioDevice) Name() string {
	return m.name
}

// ID returns the ID of the mock device.
func (m *MockAudioDevice) ID() interface{} {
	return m.id
}

// Start prepares the mock device for reading.
func (m *MockAudioDevice) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = false
	m.currentPos = 0 // Reset position on start
	fmt.Printf("MockAudioDevice '%s' started, streaming audio from '%s'.", m.name, DefaultTestWavPath)
	return nil
}

// Read reads a chunk of audio data from the pre-loaded WAV file.
// It simulates real-time audio by introducing a delay.
func (m *MockAudioDevice) Read() ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil, os.ErrClosed
	}

	// Calculate how many bytes correspond to MockChunkSize in the WAV file
	bytesPerSample := m.bitDepth / 8
	samplesPerChunk := MockChunkSize / bytesPerSample

	// Ensure we don't read beyond the end of the data, looping if necessary
	endPos := m.currentPos + samplesPerChunk*bytesPerSample
	if endPos > len(m.audioData) {
		// Loop back to the beginning if we've reached the end
		remaining := len(m.audioData) - m.currentPos
		chunk := make([]byte, samplesPerChunk*bytesPerSample)
		copy(chunk, m.audioData[m.currentPos:])
		copy(chunk[remaining:], m.audioData[0:samplesPerChunk*bytesPerSample-remaining])
		m.currentPos = samplesPerChunk*bytesPerSample - remaining
		return chunk, nil
	}

	chunk := m.audioData[m.currentPos:endPos]
	m.currentPos = endPos

	// Simulate real-time delay
	// Assuming 16kHz, 16-bit mono. MockChunkSize is 100ms of audio.
	time.Sleep(100 * time.Millisecond)

	return chunk, nil
}

// Close stops the mock device.
func (m *MockAudioDevice) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	fmt.Printf("MockAudioDevice '%s' closed.", m.name)
	return nil
}

// SampleRate returns the sample rate of the mock audio.
func (m *MockAudioDevice) SampleRate() int {
	return m.sampleRate
}

// NumChannels returns the number of channels of the mock audio.
func (m *MockAudioDevice) NumChannels() int {
	return m.numChannels
}

// BitDepth returns the bit depth of the mock audio.
func (m *MockAudioDevice) BitDepth() int {
	return m.bitDepth
}
