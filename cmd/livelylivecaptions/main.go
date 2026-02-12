package main

import (
	"fmt"
	"livelylivecaptions/internal/audio"
	"livelylivecaptions/internal/hardware"
	"livelylivecaptions/internal/transcriber"
	"livelylivecaptions/internal/types"
	"livelylivecaptions/internal/ui"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// Initialize Viper
	v := viper.New()
	v.SetConfigFile("config.yaml") // Look for config.yaml in the current directory
	v.SetConfigType("yaml")

	// Set default values (lowest priority)
	v.SetDefault("model.provider", "") // Auto-detect
	v.SetDefault("model.path", "")
	v.SetDefault("model.encoder", "")
	v.SetDefault("model.decoder", "")
	v.SetDefault("model.joiner", "")
	v.SetDefault("model.tokens", "")
	v.SetDefault("audio.sample_rate", 16000)
	v.SetDefault("audio.device_id", "") // Auto-select/prompt
	v.SetDefault("audio.monitor_mode", false)
	v.SetDefault("log.to_memory", true)
	v.SetDefault("log.file_path", "")
	v.SetDefault("log.level", "info")
	v.SetDefault("debug.enabled", false)

	// Read from config file (middle priority)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config.yaml found, using defaults and environment variables.")
		} else {
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Read from environment variables (higher priority)
	v.SetEnvPrefix("LIVELY") // e.g., LIVELY_MODEL_PROVIDER
	v.AutomaticEnv()         // Automatically bind environment variables

	// Define CLI arguments using pflag (highest priority)
	pflag.String("model.provider", "", "Force specific model provider (cpu, cuda, mock)")
	pflag.String("model.path", "", "Base path for Sherpa-ONNX models")
	pflag.String("audio.device_id", "", "ID or name of the audio device to use")
	pflag.Bool("audio.monitor_mode", false, "Enable monitor mode (capture output audio)")
	pflag.Bool("debug.enabled", false, "Enable general debug features")
	pflag.Int("audio.sample_rate", 16000, "Sample rate for audio capture (Hz)")
	pflag.String("log.file_path", "", "Path to a file for persistent logging")
	pflag.String("log.level", "info", "Minimum log level to capture")
	pflag.Bool("log.to_memory", true, "Log to in-memory ring buffer for UI display")

	// Parse pflags and bind to Viper
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// Unmarshal config into AppConfig struct
	var cfg types.AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Printf("Error unmarshalling config: %v\n", err)
		os.Exit(1)
	}

	// Determine compute provider based on resolved config
	var provider hardware.Provider
	requestedGPU := (cfg.Model.Provider == hardware.ProviderCUDA) // Check if GPU was explicitly requested

	if cfg.Model.Provider != "" {
		provider = cfg.Model.Provider
	} else {
		// Auto-detect if no provider is specified in config/env/flags
		provider = hardware.DetectBestProvider()
	}

	fmt.Printf("Compute provider: %s\n", provider)

	// Log a warning if GPU was requested but CPU is being used
	if requestedGPU && provider == hardware.ProviderCPU {
		fmt.Println("Warning: GPU (CUDA) was requested but not detected/available. Falling back to CPU provider.")
	}

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

	// Use device_id from config or prompt if not set
	var selectedDevice types.AudioDevice
	if cfg.Audio.DeviceID != "" {
		found := false
		for _, device := range devices {
			// Compare by Name or ID
			if device.ID() == cfg.Audio.DeviceID || device.Name() == cfg.Audio.DeviceID {
				selectedDevice = device
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Configured audio device '%s' not found. Please select from available devices.\n", cfg.Audio.DeviceID)
			// Fall through to interactive selection if not found
		} else {
			fmt.Printf("Using configured audio device: %s (ID: %v)\n", selectedDevice.Name(), selectedDevice.ID())
		}
	}

	if selectedDevice == nil { // If no device selected by config or not found
		fmt.Print("Select a device: ")
		var selectedDeviceIndex int
		_, err = fmt.Scanln(&selectedDeviceIndex)
		if err != nil || selectedDeviceIndex < 0 || selectedDeviceIndex >= len(devices) {
			fmt.Println("Invalid selection.")
			return
		}
		selectedDevice = devices[selectedDeviceIndex]
	}
    
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