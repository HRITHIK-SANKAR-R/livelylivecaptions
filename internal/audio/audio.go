package audio

import (
	"fmt"
	"livelylivecaptions/internal/state"
	"livelylivecaptions/internal/types"
	"math"
	"time"

	"github.com/gen2brain/malgo"
)

// MalgoDevice wraps malgo.DeviceInfo to implement types.AudioDevice
type MalgoDevice struct {
	Info malgo.DeviceInfo
}

func (d MalgoDevice) Name() string {
	return d.Info.Name()
}

func (d MalgoDevice) ID() interface{} {
	return d.Info.ID
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

// Capture implements types.AudioDevice for MalgoDevice
func (d MalgoDevice) Capture(audioChan chan<- []byte, levelChan chan<- types.AudioLevelMsg, quitChan <-chan struct{}) error {
	malgoCtx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = malgoCtx.Uninit()
		malgoCtx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 16000
	deviceConfig.Capture.DeviceID = d.Info.ID.Pointer()

	onRecvFrames := func(pSample, pOutput []byte, framecount uint32) {
		sampleCopy := make([]byte, len(pSample))
		copy(sampleCopy, pSample)
		audioChan <- sampleCopy

		rms := CalculateRMS(pSample)
		if levelChan != nil {
			select {
			case levelChan <- types.AudioLevelMsg(rms):
			default:
			}
		}
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	device, err := malgo.InitDevice(malgoCtx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return err
	}
	defer device.Uninit()

	err = device.Start()
	if err != nil {
		return err
	}

	fmt.Printf("Audio capture started on device: %s\n", d.Name())
	<-quitChan
	fmt.Printf("Audio capture stopped on device: %s\n", d.Name())
	return nil
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
		wrapped[i] = MalgoDevice{Info: d}
	}
	return wrapped, nil
}

// GetAudioDevices is a legacy helper (optional, can be removed if all callers updated)
func GetAudioDevices() ([]malgo.DeviceInfo, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	return ctx.Devices(malgo.Capture)
}

// CaptureAudio bridges the application state and the selected audio device.
// It handles the start/stop logic based on appState.
func CaptureAudio(appState *state.State, audioChan chan<- []byte, levelChan chan<- types.AudioLevelMsg, quitChan <-chan struct{}, selectedDevice types.AudioDevice) {
	deviceQuitChan := make(chan struct{})
	var deviceRunning bool

	for {
		select {
		case <-quitChan:
			if deviceRunning {
				close(deviceQuitChan)
			}
			return
		default:
			if appState.IsCapturing() {
				if !deviceRunning {
					deviceRunning = true
					deviceQuitChan = make(chan struct{})
					go func() {
						err := selectedDevice.Capture(audioChan, levelChan, deviceQuitChan)
						if err != nil {
							fmt.Println("Capture error:", err)
							appState.SetCapturing(false)
						}
					}()
				}
			} else {
				if deviceRunning {
					close(deviceQuitChan)
					deviceRunning = false
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}