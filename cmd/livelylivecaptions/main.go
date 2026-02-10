package main

import (
	"bufio"
	"fmt"
	"livelylivecaptions/internal/audio"
	"livelylivecaptions/internal/multiplexer"
	"livelylivecaptions/internal/network"
	"livelylivecaptions/internal/state"
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

	// Create channels

	micAudioChan := make(chan []byte)
	systemAudioChan := make(chan []byte)
	webSocketOutgoingChan := make(chan []byte, 100) // Buffered channel
	webSocketIncomingChan := make(chan string)
	uiUpdateChan := make(chan string)
	quitChan := make(chan struct{})

	fmt.Println("Channels created")

	// Start audio capture
	go audio.CaptureAudio(appState, micAudioChan, quitChan, selectedDevice)

	// Start network communication
	go network.ManageWebSocket(webSocketOutgoingChan, webSocketIncomingChan, quitChan)

	// Start audio multiplexer
	go multiplexer.MultiplexAudio(micAudioChan, systemAudioChan, webSocketOutgoingChan, quitChan)

	// Start UI updater
	go ui.UpdateUI(uiUpdateChan, quitChan)

	// Forward captions from network to UI
	go func() {
		for caption := range webSocketIncomingChan {
			uiUpdateChan <- caption
		}
	}()

	// Handle graceful shutdown in a separate goroutine
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		<-sigs
		fmt.Println("Shutdown signal received")
		close(quitChan)
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