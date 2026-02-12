package audio

import (
	"fmt"
	"livelylivecaptions/internal/state"
	"livelylivecaptions/internal/types"
	"math"
	"time"

	"github.com/gen2brain/malgo"
)

// GetAudioDevices enumerates the available audio capture devices.
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

// CaptureAudio captures audio from a specified device and sends it to a channel.
func CaptureAudio(appState *state.State, audioChan chan<- []byte, levelChan chan<- types.AudioLevelMsg, quitChan <-chan struct{}, selectedDevice malgo.DeviceInfo) {
	// Initialize audio context
	malgoCtx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = malgoCtx.Uninit()
		malgoCtx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 16000
	deviceConfig.Capture.DeviceID = selectedDevice.ID.Pointer()

	onRecvFrames := func(pSample, pOutput []byte, framecount uint32) {
		// It's important to copy the sample data because the buffer will be reused by the audio driver.
		sampleCopy := make([]byte, len(pSample))
		copy(sampleCopy, pSample)
		audioChan <- sampleCopy

		// Calculate RMS
		var sumSquares float64
		numSamples := len(pSample) / 2
		for i := 0; i < numSamples; i++ {
			val := int16(uint16(pSample[2*i]) | uint16(pSample[2*i+1])<<8)
			normalized := float64(val) / 32768.0
			sumSquares += normalized * normalized
		}
		rms := 0.0
		if numSamples > 0 {
			rms = float64(math.Sqrt(sumSquares / float64(numSamples)))
		}

		if levelChan != nil {
			select {
			case levelChan <- types.AudioLevelMsg(rms):
			default:
				// Drop level update if channel is full to prevent blocking audio
			}
		}
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	var device *malgo.Device

	for {
		select {
		case <-quitChan:
			fmt.Println("Stopping audio capture.")
			if device != nil {
				device.Uninit()
			}
			return
		default:
			if appState.IsCapturing() {
				if device == nil {
					device, err = malgo.InitDevice(malgoCtx.Context, deviceConfig, captureCallbacks)
					if err != nil {
						fmt.Println(err)
						// Stop capturing on error
						appState.SetCapturing(false)
						continue
					}

					err = device.Start()
					if err != nil {
						fmt.Println(err)
						// Stop capturing on error
						appState.SetCapturing(false)
						continue
					}
					fmt.Println("Audio capture started.")
				}
			} else {
				if device != nil {
					device.Uninit()
					device = nil
					fmt.Println("Audio capture stopped.")
				}
			}
		}
		// Add a small sleep to prevent busy-waiting
		time.Sleep(100 * time.Millisecond)
	}
}