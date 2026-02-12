package main

import (
	"flag"
	"fmt"
	"livelylivecaptions/internal/audio"
	"livelylivecaptions/internal/hardware"
	"livelylivecaptions/internal/transcriber"
	"livelylivecaptions/internal/types"
	"livelylivecaptions/internal/ui"
	"os"
	"time" // Added for time.Sleep
)

func main() {
	// Parse CLI flags
	useCPU := flag.Bool("cpu", false, "Force CPU usage")
	useGPU := flag.Bool("gpu", false, "Force GPU (CUDA) usage")
	flag.Parse()

	// Initialize the application state - appState is no longer directly used for capture control
	// appState := state.NewState() // Removed

	// Determine compute provider
	var provider hardware.Provider
	requestedGPU := *useGPU // Store the initial request
	if requestedGPU {
		provider = hardware.ProviderCUDA
	} else if *useCPU {
		provider = hardware.ProviderCPU
	} else {
		provider = hardware.DetectBestProvider()
	}
	fmt.Printf("Compute provider: %s\n", provider)

	// Log a warning if GPU was requested but CPU is being used
	if requestedGPU && provider == hardware.ProviderCPU {
		fmt.Println("Warning: GPU (CUDA) was requested but not detected/available. Falling back to CPU provider.")
	}

	// Get audio devices using the new Provider interface
	audioProvider := audio.MalgoProvider{} // Use the MalgoProvider
	devices, err := audioProvider.GetDevices()
	if err != nil {
		fmt.Println("Failed to get audio devices:", err)
		return
	}

	fmt.Println("Available audio devices:")
	for i, device := range devices {
		fmt.Printf("%d: %s\n", i, device.Name())
	}

	fmt.Print("Select a device: ")
	var selectedDeviceIndex int
	_, err = fmt.Scanln(&selectedDeviceIndex)
	if err != nil || selectedDeviceIndex < 0 || selectedDeviceIndex >= len(devices) {
		fmt.Println("Invalid selection.")
		return
	}

	selectedDevice := devices[selectedDeviceIndex]
    
    // Start the selected audio device
    err = selectedDevice.Start()
    if err != nil {
        fmt.Printf("Failed to start audio device: %v\n", err)
        return
    }
    defer selectedDevice.Close() // Ensure device is closed on exit

	// Initialize Transcriber
	tr, err := transcriber.NewTranscriber(provider)
	if err != nil {
		fmt.Printf("Failed to initialize transcriber: %v\n", err)
		return
	}
	defer tr.Close()
	fmt.Println("Transcriber initialized (Sherpa-ONNX)")

	// Create channels
	micAudioChan := tr.InputChan
	uiUpdateChan := tr.OutputChan
	levelChan := make(chan types.AudioLevelMsg, 60) // Buffer for 60fps
	quitChan := make(chan struct{})

	fmt.Println("Channels created")

	// Start audio capture goroutine using the new AudioDevice interface
	go func() {
		defer close(micAudioChan) // Close the input channel when audio capture stops
        defer close(levelChan)    // Close level channel as well

		for {
			select {
			case <-quitChan:
				return
			default:
				audioData, err := selectedDevice.Read()
				if err != nil {
					fmt.Printf("Error reading from audio device: %v\n", err)
					return // Exit goroutine on error
				}
				
				if audioData == nil {
					// No data yet, wait a bit to prevent busy-looping
					time.Sleep(10 * time.Millisecond)
					continue
				}

				// Send audio data to transcriber
				micAudioChan <- audioData

				// Calculate and send RMS level
				rms := audio.CalculateRMS(audioData)
				select {
				case levelChan <- types.AudioLevelMsg(rms):
				default: // Non-blocking send to levelChan
				}
			}
		}
	}()

	// Start transcriber
	tr.Start()

    // Initialize and run Bubble Tea program
    if err := ui.RunProgram(uiUpdateChan, levelChan, quitChan); err != nil {
        fmt.Printf("Error running UI: %v\n", err)
        os.Exit(1)
    }

	// Cleanup after UI exits
	fmt.Println("\nShutting down gracefully...")
}