package main

import (
	"flag"
	"fmt"
	"livelylivecaptions/internal/audio"
	"livelylivecaptions/internal/hardware"
	"livelylivecaptions/internal/state"
	"livelylivecaptions/internal/transcriber"
	"livelylivecaptions/internal/types"
	"livelylivecaptions/internal/ui"
	"os"
)

func main() {
	// Parse CLI flags
	useCPU := flag.Bool("cpu", false, "Force CPU usage")
	useGPU := flag.Bool("gpu", false, "Force GPU (CUDA) usage")
	flag.Parse()

	// Initialize the application state
	appState := state.NewState()
	fmt.Println("Application state initialized")

	// Determine compute provider
	var provider hardware.Provider
	if *useGPU {
		provider = hardware.ProviderCUDA
	} else if *useCPU {
		provider = hardware.ProviderCPU
	} else {
		provider = hardware.DetectBestProvider()
	}
	fmt.Printf("Compute provider: %s\n", provider)

	// Get audio devices using the new Provider interface
	audioProvider := audio.MalgoProvider{}
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

	// Enable audio capturing
	appState.SetCapturing(true)

	// Start audio capture
	go audio.CaptureAudio(appState, micAudioChan, levelChan, quitChan, selectedDevice)

	// Start transcriber
	tr.Start()

	// Start UI updater
	// We'll replace this with the Bubble Tea program in the next step, 
	// but keeping the signature compatible for now or removing if we switch to full TUI immediately.
	// For this step, let's comment out the old UI updater as we are about to replace it.
	// go ui.UpdateUI(uiUpdateChan, quitChan)
    
    // Initialize and run Bubble Tea program
    if err := ui.RunProgram(uiUpdateChan, levelChan, quitChan); err != nil {
        fmt.Printf("Error running UI: %v\n", err)
        os.Exit(1)
    }

	// Cleanup after UI exits
	fmt.Println("\nShutting down gracefully...")
}