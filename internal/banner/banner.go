package banner

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
)

// RGB creates an ANSI 24-bit color code
func RGB(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// interpolateColor calculates the color at position t between start and end colors
func interpolateColor(startR, startG, startB, endR, endG, endB int, t float64) (int, int, int) {
	r := int(float64(startR) + (float64(endR-startR) * t))
	g := int(float64(startG) + (float64(endG-startG) * t))
	b := int(float64(startB) + (float64(endB-startB) * t))
	return r, g, b
}

// ASCII art for "LIVELYLIVECAPTIONS" (Standard)
var asciiArt = []string{
	" _      _           _         _     _           ____             _  _                ",
	"| |    (_)_   _____| |_   _  | |   (_)_   _____|  _ \\ __ _ _ __ | |_(_) ___  _ __   ___",
	"| |    | \\ \\ / / _ \\ | | | | | |   | \\ \\ / / _ \\ |_) / _` | '_ \\| __| |/ _ \\| '_ \\/ __|",
	"| |___| |\\ V /  __/ | |_| | | |___| |\\ V /  __/  __/ (_| | |_) | |_| | (_) | | | \\__ \\",
	"|_____|_|_\\_/ \\___|_|\\__, | |_____|_|_\\_/ \\___|_|   \\__,_| .__/ \\__|_|\\___/|_| |_|___/",
	"                    |___/                             |_|                          ",
}

// Block style ASCII art for "LIVELY LIVE CAPTIONS"
// Block style ASCII art for "LIVELY LIVE" (Line 1) and "CAPTIONS" (Line 2)
var blockArt = []string{
	// Line 1: LIVELY LIVE
	"██╗     ██╗██╗   ██╗███████╗██╗  ██╗   ██╗    ██╗     ██╗██╗   ██╗███████╗",
	"██║     ██║██║   ██║██╔════╝██║  ╚██╗ ██╔╝    ██║     ██║██║   ██║██╔════╝",
	"██║     ██║██║   ██║█████╗  ██║   ╚████╔╝     ██║     ██║██║   ██║█████╗  ",
	"██║     ██║╚██╗ ██╔╝██╔══╝  ██║    ╚██╔╝      ██║     ██║╚██╗ ██╔╝██╔══╝  ",
	"███████╗██║ ╚████╔╝ ███████╗███████╗██║       ███████╗██║ ╚████╔╝ ███████╗",
	"╚══════╝╚═╝  ╚═══╝  ╚══════╝╚══════╝╚═╝       ╚══════╝╚═╝  ╚═══╝  ╚══════╝",
	"",
	// Line 2: CAPTIONS (Centered relative to the top line)
	"      ██████╗ █████╗ ██████╗ ████████╗██╗ ██████╗ ███╗   ██╗███████╗",
	"     ██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██║██╔═══██╗████╗  ██║██╔════╝",
	"     ██║     ███████║██████╔╝   ██║   ██║██║   ██║██╔██╗ ██║███████╗",
	"     ██║     ██╔══██║██╔═══╝    ██║   ██║██║   ██║██║╚██╗██║╚════██║",
	"     ╚██████╗██║  ██║██║        ██║   ██║╚██████╔╝██║ ╚████║███████║",
	"      ╚═════╝╚═╝  ╚═╝╚╝        ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝",
}

// Print displays the Gemini-style banner with blue-to-pink gradient
func Print() {
	// Gemini gradient: Blue → Pink
	startR, startG, startB := 20, 90, 200    // Deep blue
	endR, endG, endB := 255, 130, 140        // Pink

	fmt.Println()

	for _, line := range blockArt {
		lineLen := len(line)
		for charIdx, char := range line {
			if char != ' ' {
				// Calculate gradient position (0.0 to 1.0)
				t := float64(charIdx) / float64(lineLen)
				if lineLen == 0 { t = 0 } // Prevent divide by zero for empty lines

				// Get interpolated color
				r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)

				// Print character with color
				fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
	printWave(startR, startG, startB, endR, endG, endB)
	fmt.Println()

	// Status message in cyan/magenta
	fmt.Println(RGB(150, 100, 200) + "    > Listening for audio input..." + Reset)
}

// PrintOrangePeach displays an Orange → Peach gradient banner
func PrintOrangePeach() {
	startR, startG, startB := 255, 140, 0    // Orange
	endR, endG, endB := 255, 218, 185        // Peach

	fmt.Println()

	for _, line := range blockArt {
		lineLen := len(line)
		for charIdx, char := range line {
			if char != ' ' {
				t := float64(charIdx) / float64(lineLen)
				if lineLen == 0 { t = 0 }
				r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
				fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
	printWave(startR, startG, startB, endR, endG, endB)
	fmt.Println()

	fmt.Println(RGB(255, 180, 150) + "    > Listening for audio input..." + Reset)
}

// PrintFireSunset displays a Fire → Sunset gradient banner
func PrintFireSunset() {
	startR, startG, startB := 255, 69, 0     // Fire Orange-Red
	endR, endG, endB := 255, 165, 0          // Sunset Orange

	fmt.Println()

	for _, line := range blockArt {
		lineLen := len(line)
		for charIdx, char := range line {
			if char != ' ' {
				t := float64(charIdx) / float64(lineLen)
				if lineLen == 0 { t = 0 }
				r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
				fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
	printWave(startR, startG, startB, endR, endG, endB)
	fmt.Println()

	fmt.Println(RGB(255, 100, 50) + "    > Listening for audio input..." + Reset)
}

// PrintGoldenHour displays a Golden Hour (Orange → Yellow) gradient banner
func PrintGoldenHour() {
	startR, startG, startB := 255, 165, 0    // Orange
	endR, endG, endB := 255, 255, 0          // Yellow

	fmt.Println()

	for _, line := range blockArt {
		lineLen := len(line)
		for charIdx, char := range line {
			if char != ' ' {
				t := float64(charIdx) / float64(lineLen)
				if lineLen == 0 { t = 0 }
				r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
				fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
	printWave(startR, startG, startB, endR, endG, endB)
	fmt.Println()

	fmt.Println(RGB(255, 200, 0) + "    > Listening for audio input..." + Reset)
}

// PrintAmberGlow displays an Amber Glow (Warm & Elegant) gradient banner
func PrintAmberGlow() {
	startR, startG, startB := 255, 191, 0    // Amber
	endR, endG, endB := 255, 215, 0          // Gold

	fmt.Println()

	for _, line := range blockArt {
		lineLen := len(line)
		for charIdx, char := range line {
			if char != ' ' {
				t := float64(charIdx) / float64(lineLen)
				if lineLen == 0 { t = 0 }
				r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
				fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
	printWave(startR, startG, startB, endR, endG, endB)
	fmt.Println()

	fmt.Println(RGB(255, 200, 50) + "    > Listening for audio input..." + Reset)
}

// PrintStandard uses the standard ASCII art
func PrintStandard() {
	startR, startG, startB := 20, 90, 200
	endR, endG, endB := 255, 130, 140

	fmt.Println()

	for _, line := range asciiArt {
		lineLen := len(line)
		for charIdx, char := range line {
			if char != ' ' {
				t := float64(charIdx) / float64(lineLen)
				r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
				fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
	printWave(startR, startG, startB, endR, endG, endB)
	fmt.Println()
	fmt.Println(RGB(150, 100, 200) + "    > Listening for audio input..." + Reset)
}

// PrintCompact displays a compact version
func PrintCompact() {
	compactArt := []string{
		" _      _           _          ____             _  _                ",
		"| |    (_)_   _____| |_   _   / ___|__ _ _ __ | |_(_) ___  _ __   ___",
		"| |    | \\ \\ / / _ \\ | | | | | |   / _` | '_ \\| __| |/ _ \\| '_ \\/ __|",
		"| |___| |\\ V /  __/ | |_| |  | |__| (_| | |_) | |_| | (_) | | | \\__ \\",
		"|_____|_| \\_/ \\___|_|\\__, |   \\____\\__,_| .__/ \\__|_|\\___/|_| |_|___/",
		"                    |___/            |_|                          ",
	}

	startR, startG, startB := 20, 90, 200
	endR, endG, endB := 255, 130, 140

	fmt.Println()

	for _, line := range compactArt {
		lineLen := len(line)
		for charIdx, char := range line {
			if char != ' ' {
				t := float64(charIdx) / float64(lineLen)
				r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
				fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println(RGB(150, 100, 200) + "    Real-time AI-powered transcription" + Reset)
}

// PrintWithVersion displays banner with version
func PrintWithVersion(version string) {
	Print()
	versionText := "v" + version + " • Real-time Transcription"
	fmt.Println(RGB(150, 100, 200) + "    " + versionText + Reset)
}

// PrintWithStats displays banner with audio stats
func PrintWithStats(sampleRate int, channels int, device string) {
	Print()
	fmt.Printf(RGB(150, 100, 200)+"    📊 Sample Rate: %d Hz | Channels: %d\n"+Reset, sampleRate, channels)
	fmt.Printf(RGB(150, 100, 200)+"    🎤 Device: %s\n"+Reset, device)
}

// PrintCustomText prints any text with Gemini gradient
func PrintCustomText(text string) {
	startR, startG, startB := 20, 90, 200
	endR, endG, endB := 255, 130, 140

	textLen := len(text)
	for i, char := range text {
		t := float64(i) / float64(textLen)
		if textLen == 0 { t = 0 }
		r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
		fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
	}
	fmt.Println()
}

// PrintMinimal - Minimal clean version
func PrintMinimal() {
	startR, startG, startB := 20, 90, 200
	endR, endG, endB := 255, 130, 140

	fmt.Println()

	// Render "Lively Live Captions" with gradient
	text := "Lively Live Captions"
	textLen := len(text)
	for i, char := range text {
		t := float64(i) / float64(textLen)
		r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
		fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
	}
	fmt.Println()

	fmt.Println()
	fmt.Println(RGB(150, 100, 200) + "  ► Ready" + Reset)
}

// printWave displays a decorative wave with gradient
func printWave(startR, startG, startB, endR, endG, endB int) {
	wave := "~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~"
	waveLen := len(wave)

	fmt.Print("    ")
	for i, char := range wave {
		t := float64(i) / float64(waveLen)
		r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
		fmt.Print(RGB(r, g, b) + string(char) + Reset)
	}
	fmt.Println()
}

// printGradientWave - Wave matching the banner gradient
func printGradientWave() {
	wave := "~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~≈~"
	waveLen := len(wave)

	startR, startG, startB := 20, 90, 200
	endR, endG, endB := 255, 130, 140

	fmt.Print("    ")
	for i, char := range wave {
		t := float64(i) / float64(waveLen)
		r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
		fmt.Print(RGB(r, g, b) + string(char) + Reset)
	}
	fmt.Println()
}

// PrintOceanBlue displays ocean blue gradient (blue → cyan)
func PrintOceanBlue() {
	startR, startG, startB := 15, 45, 85     // Deep ocean blue
	endR, endG, endB := 130, 220, 255        // Bright cyan

	fmt.Println()

	for _, line := range blockArt {
		lineLen := len(line)
		for charIdx, char := range line {
			if char != ' ' {
				t := float64(charIdx) / float64(lineLen)
				if lineLen == 0 { t = 0 }
				r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
				fmt.Print(RGB(r, g, b) + Bold + string(char) + Reset)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
	printWave(startR, startG, startB, endR, endG, endB)
	fmt.Println()
	fmt.Println(RGB(100, 180, 255) + "    > Listening for audio input..." + Reset)
}

// GetGradientText returns text with gradient applied (for custom use)
func GetGradientText(text string) string {
	startR, startG, startB := 20, 90, 200
	endR, endG, endB := 255, 130, 140

	var result strings.Builder
	textLen := len(text)

	for i, char := range text {
		t := float64(i) / float64(textLen)
		r, g, b := interpolateColor(startR, startG, startB, endR, endG, endB, t)
		result.WriteString(RGB(r, g, b) + Bold + string(char) + Reset)
	}

	return result.String()
}
