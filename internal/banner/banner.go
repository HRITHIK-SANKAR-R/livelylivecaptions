package banner

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

// Print displays the Gemini-style blue-to-pink gradient banner
func Print() {
	// Define gradient colors (Blue → Pink like Gemini)
	startColor := pterm.NewRGB(20, 90, 200)    // Deep blue
	endColor := pterm.NewRGB(255, 130, 140)    // Pink
	
	// Create the gradient
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	// Create letters
	letters := putils.LettersFromString("LivelyLiveCaptions")
	
	// Render the big text
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	// Print with gradient applied
	gradient.Println(bigText)
	
	// Wave decoration
	printWave()
	
	// Status message
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    > Listening for audio input...")
}

// PrintGeminiStyle - Alternative vibrant gradient
func PrintGeminiStyle() {
	startColor := pterm.NewRGB(15, 80, 180)    // Electric blue
	endColor := pterm.NewRGB(255, 120, 160)    // Hot pink
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	letters := putils.LettersFromString("LivelyLiveCaptions")
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	gradient.Println(bigText)
	printWave()
	
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    > Listening for audio input...")
}

// PrintGeminiExact - Closest to the Gemini screenshot
func PrintGeminiExact() {
	startColor := pterm.NewRGB(20, 90, 200)    // Blue
	endColor := pterm.NewRGB(255, 130, 140)    // Pink
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	letters := putils.LettersFromString("LivelyLiveCaptions")
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	gradient.Println(bigText)
	printGradientWave()
	
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    > Listening for audio input...")
}

// PrintTwoWords - Display "Lively" and "LiveCaptions" as separate words
func PrintTwoWords() {
	startColor := pterm.NewRGB(20, 90, 200)    // Blue
	endColor := pterm.NewRGB(255, 130, 140)    // Pink
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	lively := putils.LettersFromString("Lively")
	captions := putils.LettersFromString("LiveCaptions")
	
	bigText, _ := pterm.DefaultBigText.WithLetters(lively, captions).Srender()
	
	gradient.Println(bigText)
	printGradientWave()
	
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    > Listening for audio input...")
}

// PrintCompact - Compact single-line version
func PrintCompact() {
	startColor := pterm.NewRGB(20, 90, 200)    // Blue
	endColor := pterm.NewRGB(255, 130, 140)    // Pink
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	letters := putils.LettersFromString("LivelyCaptions")
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	gradient.Println(bigText)
	
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    Real-time AI-powered transcription")
}

// PrintWithCustomText - Display custom text with Gemini gradient
func PrintWithCustomText(text string) {
	startColor := pterm.NewRGB(20, 90, 200)
	endColor := pterm.NewRGB(255, 130, 140)
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	letters := putils.LettersFromString(text)
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	gradient.Println(bigText)
}

// PrintWithSubtitle - Banner with subtitle
func PrintWithSubtitle() {
	startColor := pterm.NewRGB(20, 90, 200)
	endColor := pterm.NewRGB(255, 130, 140)
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	letters := putils.LettersFromString("LivelyLiveCaptions")
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	gradient.Println(bigText)
	printGradientWave()
	
	// Subtitle with gradient
	subtitle := "Real-time AI-Powered Speech Recognition"
	subtitleGradient := pterm.NewRGB(100, 120, 200).Fade(0, 1, 0, pterm.NewRGB(255, 150, 170))
	subtitleGradient.Println("    " + subtitle)
	
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    > Listening for audio input...")
}

// PrintWithVersion - Banner with version info
func PrintWithVersion(version string) {
	startColor := pterm.NewRGB(20, 90, 200)
	endColor := pterm.NewRGB(255, 130, 140)
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	letters := putils.LettersFromString("LivelyLiveCaptions")
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	gradient.Println(bigText)
	printGradientWave()
	
	// Version info
	versionText := "v" + version + " • Real-time Transcription"
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    " + versionText)
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    > Listening for audio input...")
}

// PrintWithStats - Banner with system stats
func PrintWithStats(sampleRate int, channels int, device string) {
	startColor := pterm.NewRGB(20, 90, 200)
	endColor := pterm.NewRGB(255, 130, 140)
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	letters := putils.LettersFromString("LivelyLiveCaptions")
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	gradient.Println(bigText)
	printGradientWave()
	
	// System info
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Printfln("    📊 Sample Rate: %d Hz | Channels: %d", sampleRate, channels)
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Printfln("    🎤 Device: %s", device)
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("    > Listening for audio input...")
}

// PrintMinimal - Minimal clean version
func PrintMinimal() {
	startColor := pterm.NewRGB(20, 90, 200)
	endColor := pterm.NewRGB(255, 130, 140)
	
	gradient := startColor.Fade(0, 1, 0, endColor)
	
	letters := putils.LettersFromString("LivelyCaptions")
	bigText, _ := pterm.DefaultBigText.WithLetters(letters).Srender()
	
	gradient.Println(bigText)
	pterm.Println()
	
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightMagenta)).
		Println("  ► Ready")
}

// printWave - Decorative wave with gradient
func printWave() {
	wave := "~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~"
	
	waveStart := pterm.NewRGB(30, 100, 200)
	waveEnd := pterm.NewRGB(255, 100, 150)
	waveGradient := waveStart.Fade(0, 1, 0, waveEnd)
	
	waveGradient.Println("    " + wave)
}

// printGradientWave - Wave matching the banner gradient
func printGradientWave() {
	wave := "~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~"
	
	waveStart := pterm.NewRGB(20, 90, 200)
	waveEnd := pterm.NewRGB(255, 130, 140)
	waveGradient := waveStart.Fade(0, 1, 0, waveEnd)
	
	waveGradient.Println("    " + wave)
}
