package ui_test

import (
	"context"
	"livelylivecaptions/internal/types"
	"livelylivecaptions/internal/ui"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestUI(t *testing.T) {
	// Create mock channels for the UI model
	transChan := make(chan types.TranscriptionEvent)
	levelChan := make(chan types.AudioLevelMsg)
	quitChan := make(chan struct{})

	// Initialize the UI model
	m := ui.InitialModel(transChan, levelChan, quitChan)

	// Create a test program
	ctx := context.Background()
	tm := teatest.NewTestProgram(t, m,
		teatest.WithInitialTermSize(120, 25), // Set initial terminal size
		teatest.WithSendTimeout(time.Second),
	)

	// Start the test program
	go tm.Run(ctx)
	
	// Wait for Init
	tm.Wait(m.Init())

	// Send some events
	tm.Send(types.TranscriptionEvent{Text: "Hello world", IsFinal: false})
	tm.Send(types.TranscriptionEvent{Text: "Hello world, how are you?", IsFinal: true})
	tm.Send(types.TranscriptionEvent{Text: "I am fine, thank you.", IsFinal: false})
	tm.Send(types.TranscriptionEvent{Text: "I am fine, thank you. And you?", IsFinal: true})
	tm.Send(types.TranscriptionEvent{Text: "This is a test.", IsFinal: false})
	tm.Send(types.TranscriptionEvent{Text: "This is a test transcription.", IsFinal: true})

	// Send some audio level messages
	tm.Send(types.AudioLevelMsg(0.1))
	tm.Send(types.AudioLevelMsg(0.5))
	tm.Send(types.AudioLevelMsg(0.9))

	// Send a quit message
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// Wait for the program to exit
	tm.Wait(tea.Quit)

	// Check the final output against a golden snapshot
	teatest.RequireEqual(t, tm.FinalOutput(), "testdata/final_output.golden")
}
