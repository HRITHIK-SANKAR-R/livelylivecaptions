package main

import (
	"bufio"
	"fmt"
	"livelylivecaptions/internal/audio"
	"livelylivecaptions/internal/state"
	"livelylivecaptions/internal/transcriber"
	"livelylivecaptions/internal/ui"
	"os"
	"os/signal"
)

func main() {
	// Initialize the application state
	appState := state.NewState()
	fmt.Println("Application state initialized")

	// Get audio devices
	devices, err := audio.GetAudioDevices()
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
	tr, err := transcriber.NewTranscriber()
	if err != nil {
		fmt.Printf("Failed to initialize transcriber: %v\n", err)
		return
	}
	defer tr.Close()
	fmt.Println("Transcriber initialized (Sherpa-ONNX)")

	// Create channels
	micAudioChan := tr.InputChan
	uiUpdateChan := tr.OutputChan
	quitChan := make(chan struct{})

	fmt.Println("Channels created")

	// Start audio capture
	go audio.CaptureAudio(appState, micAudioChan, quitChan, selectedDevice)

	// Start transcriber
	tr.Start()

	// Start UI updater
	go ui.UpdateUI(uiUpdateChan, quitChan)

	// Handle graceful shutdown
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		<-sigs
		fmt.Println("Shutdown signal received")
		close(quitChan)
		close(tr.QuitChan)
	}()

	// Interactive toggle for capturing
	fmt.Println("Press Enter to start/stop captioning.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		appState.SetCapturing(!appState.IsCapturing())
		if appState.IsCapturing() {
			fmt.Println("Starting captioning...")
		} else {
			fmt.Println("Stopping captioning...")
		}
	}

	fmt.Println("Application shutting down")
}