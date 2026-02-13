package audio

import (
	"fmt"
	"livelylivecaptions/internal/logger" // Added import
	"livelylivecaptions/internal/types"
	"math"
	"sync"
	"time"

	"github.com/gen2brain/malgo"
)

// MalgoDevice implements types.AudioDevice using the malgo library.
// It wraps malgo.DeviceInfo and manages the underlying malgo device context.
type MalgoDevice struct {
	Info        malgo.DeviceInfo
	device      *malgo.Device
	context     *malgo.AllocatedContext
	audioBuffer chan []byte    // Channel to push audio frames from malgo callback to Read()
	quitRead    chan struct{}  // Signal to stop the internal goroutine feeding audioBuffer
	bufferMutex sync.Mutex     // Protects access to audioBuffer and related state
	isCapturing bool           // Tracks if the device is actively capturing
}

// Name returns the name of the malgo device.
func (d *MalgoDevice) Name() string {
	return d.Info.Name()
}

// ID returns the ID of the malgo device.
func (d *MalgoDevice) ID() interface{} {
	return d.Info.ID
}

// Start initializes and starts the malgo audio capture device.
func (d *MalgoDevice) Start() error {
	d.bufferMutex.Lock()
	if d.isCapturing {
		d.bufferMutex.Unlock()
		return fmt.Errorf("malgo device already started")
	}
	d.isCapturing = true
	d.audioBuffer = make(chan []byte, 100) // Buffered channel for audio frames
	d.quitRead = make(chan struct{})
	d.bufferMutex.Unlock()

	var err error
	d.context, err = malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		d.isCapturing = false
		return fmt.Errorf("failed to initialize malgo context: %w", err)
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16 // 16-bit signed integer
	deviceConfig.Capture.Channels = 1             // Mono
	deviceConfig.SampleRate = 16000               // 16 kHz sample rate
	deviceConfig.Capture.DeviceID = d.Info.ID.Pointer()

	onRecvFrames := func(pSample, pOutput []byte, framecount uint32) {
		if !d.isCapturing {
			return
		}
		sampleCopy := make([]byte, len(pSample))
		copy(sampleCopy, pSample)
		select {
		case d.audioBuffer <- sampleCopy:
		default:
			// Drop frame if buffer is full to avoid blocking audio callback
			// logger.Warn("Audio buffer full, dropping frame.") // Use logger here if desired
		}
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	d.device, err = malgo.InitDevice(d.context.Context, deviceConfig, captureCallbacks)
	if err != nil {
		d.context.Uninit()
		d.context.Free()
		d.isCapturing = false
		return fmt.Errorf("failed to initialize malgo device: %w", err)
	}

	err = d.device.Start()
	if err != nil {
		d.device.Uninit()
		d.context.Uninit()
		d.context.Free()
		d.isCapturing = false
		return fmt.Errorf("failed to start malgo device: %w", err)
	}

	logger.Info("Audio capture started on device: %s (ID: %v)", d.Name(), d.ID()) // Changed
	return nil
}

// Read reads a chunk of audio data from the device.
// It blocks until data is available or the device is closed.
func (d *MalgoDevice) Read() ([]byte, error) {
	select {
	case <-d.quitRead:
		return nil, fmt.Errorf("malgo device closed")
	case frame := <-d.audioBuffer:
		return frame, nil
	case <-time.After(100 * time.Millisecond): // Timeout for reading
		// If no data comes in a while, check if capturing is still true
		d.bufferMutex.Lock()
		capturing := d.isCapturing
		d.bufferMutex.Unlock()
		if !capturing {
			return nil, fmt.Errorf("malgo device not capturing")
		}
		return nil, nil // Return nil, nil to indicate no data yet, but still capturing
	}
}

// Close stops the malgo audio capture device and cleans up resources.
func (d *MalgoDevice) Close() error {
	d.bufferMutex.Lock()
	if !d.isCapturing {
		d.bufferMutex.Unlock()
		return fmt.Errorf("malgo device not capturing or already closed")
	}
	d.isCapturing = false
	close(d.quitRead) // Signal Read() to stop blocking
	d.bufferMutex.Unlock()

	if d.device != nil {
		d.device.Uninit()
		d.device = nil
	}
	if d.context != nil {
		d.context.Uninit()
		d.context.Free()
		d.context = nil
	}
	logger.Info("Audio capture stopped on device: %s (ID: %v)", d.Name(), d.ID()) // Changed
	return nil
}

// CalculateRMS computes the RMS value for a slice of raw int16 LE bytes.
func CalculateRMS(audioData []byte) float64 {
	var sumSquares float64
	numSamples := len(audioData) / 2
	for i := 0; i < numSamples; i++ {
		val := int16(uint16(audioData[2*i]) | uint16(audioData[2*i+1])<<8)
		normalized := float64(val) / 32768.0
		sumSquares += normalized * normalized
	}
	if numSamples > 0 {
		return math.Sqrt(sumSquares / float64(numSamples))
	}
	return 0.0
}

// MalgoProvider implements types.AudioProvider using the malgo library
type MalgoProvider struct{}

func (p MalgoProvider) GetDevices() ([]types.AudioDevice, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	devices, err := ctx.Devices(malgo.Capture)
	if err != nil {
		return nil, err
	}

	wrapped := make([]types.AudioDevice, len(devices))
	for i, d := range devices {
		// Return MalgoDevice pointers as types.AudioDevice
		wrapped[i] = &MalgoDevice{Info: d} 
	}
	return wrapped, nil
}