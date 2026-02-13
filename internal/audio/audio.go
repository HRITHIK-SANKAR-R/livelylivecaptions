package audio

import (
	"fmt"
	"livelylivecaptions/internal/logger"
	"livelylivecaptions/internal/types"
	"math"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

// PortAudioProvider implements the AudioProvider interface using portaudio.
type PortAudioProvider struct{}

func (p PortAudioProvider) GetDevices() ([]types.AudioDevice, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize portaudio: %w", err)
	}
	// Note: This means Init/Terminate are called frequently. A global management
	// strategy would be better for a more complex app.

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to get portaudio devices: %w", err)
	}

	audioDevices := make([]types.AudioDevice, 0)
	for _, deviceInfo := range devices {
		if deviceInfo.MaxInputChannels > 0 {
			// Important: create a new variable for the pointer.
			dev := deviceInfo
			audioDevices = append(audioDevices, &PortAudioDevice{Info: dev})
		}
	}

	return audioDevices, nil
}

// PortAudioDevice implements the AudioDevice interface using portaudio.
type PortAudioDevice struct {
	Info        *portaudio.DeviceInfo
	stream      *portaudio.Stream
	audioBuffer chan []byte
	quitRead    chan struct{}
	isCapturing bool
	mutex       sync.Mutex
}

func (d *PortAudioDevice) Name() string {
	return d.Info.Name
}

func (d *PortAudioDevice) ID() interface{} {
	// Using the device name as a unique ID. PortAudio doesn't provide a stable ID string.
	return d.Info.Name
}

func (d *PortAudioDevice) Start() error {
	d.mutex.Lock()
	if d.isCapturing {
		d.mutex.Unlock()
		return fmt.Errorf("portaudio device already started")
	}
	d.audioBuffer = make(chan []byte, 100)
	d.quitRead = make(chan struct{})
	d.mutex.Unlock()

	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize portaudio for capture: %w", err)
	}

	processAudio := func(in []int16) {
		// This callback is executed by PortAudio's processing thread.
		d.mutex.Lock()
		if !d.isCapturing {
			d.mutex.Unlock()
			return
		}
		d.mutex.Unlock()

		// Convert int16 samples to a byte slice (Little Endian).
		byteBuf := make([]byte, len(in)*2)
		for i, sample := range in {
			byteBuf[i*2] = byte(sample)
			byteBuf[i*2+1] = byte(sample >> 8)
		}

		select {
		case d.audioBuffer <- byteBuf:
		default:
			// Non-blocking, so we don't hold up PortAudio's thread.
			// logger.Warn is commented out to avoid log spam on busy systems.
			// logger.Warn("PortAudio buffer full, dropping frame.")
		}
	}

	streamParams := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   d.Info,
			Channels: 1, // Mono
			Latency:  time.Millisecond * 100,
		},
		Output:     portaudio.StreamDeviceParameters{}, // No output
		SampleRate: 16000.0,
	}

	var err error
	d.stream, err = portaudio.OpenStream(streamParams, processAudio)
	if err != nil {
		portaudio.Terminate()
		return fmt.Errorf("failed to open portaudio stream: %w", err)
	}

	if err := d.stream.Start(); err != nil {
		d.stream.Close()
		portaudio.Terminate()
		return fmt.Errorf("failed to start portaudio stream: %w", err)
	}

	d.mutex.Lock()
	d.isCapturing = true
	d.mutex.Unlock()

	logger.Info("Audio capture started on device: %s", d.Name())
	return nil
}

func (d *PortAudioDevice) Read() ([]byte, error) {
	select {
	case <-d.quitRead:
		return nil, fmt.Errorf("portaudio device closed")
	case frame := <-d.audioBuffer:
		return frame, nil
	case <-time.After(100 * time.Millisecond):
		d.mutex.Lock()
		capturing := d.isCapturing
		d.mutex.Unlock()
		if !capturing {
			return nil, fmt.Errorf("portaudio device not capturing")
		}
		return nil, nil // Indicate no data yet, but still capturing
	}
}

func (d *PortAudioDevice) Close() error {
	d.mutex.Lock()
	if !d.isCapturing {
		d.mutex.Unlock()
		return nil
	}
	d.isCapturing = false
	close(d.quitRead)
	d.mutex.Unlock()

	var firstErr error
	if d.stream != nil {
		if err := d.stream.Stop(); err != nil {
			firstErr = fmt.Errorf("failed to stop portaudio stream: %w", err)
			logger.Warn(firstErr.Error())
		}
		if err := d.stream.Close(); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("failed to close portaudio stream: %w", err)
			}
			logger.Warn("Error closing portaudio stream: %v", err)
		}
	}

	if err := portaudio.Terminate(); err != nil {
		if firstErr == nil {
			firstErr = fmt.Errorf("failed to terminate portaudio: %w", err)
		}
		logger.Warn("Error terminating portaudio: %v", err)
	}

	logger.Info("Audio capture stopped on device: %s", d.Name())
	return firstErr
}

// CalculateRMS computes the RMS value for a slice of raw int16 LE bytes.
// This function is independent of the audio capture library.
func CalculateRMS(audioData []byte) float64 {
	if len(audioData) == 0 {
		return 0.0
	}
	var sumSquares float64
	numSamples := len(audioData) / 2
	for i := 0; i < numSamples; i++ {
		// Little-endian conversion
		val := int16(uint16(audioData[2*i]) | uint16(audioData[2*i+1])<<8)
		normalized := float64(val) / 32768.0
		sumSquares += normalized * normalized
	}
	return math.Sqrt(sumSquares / float64(numSamples))
}
