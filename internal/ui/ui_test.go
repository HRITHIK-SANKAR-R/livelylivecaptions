package ui_test

import (
	"io"
	"livelylivecaptions/internal/types"
	"livelylivecaptions/internal/ui"
	"testing"

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
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 25))

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
	tm.WaitFinished(t)

	// Check the final output against a golden snapshot
	finalOutput, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Fatal(err)
	}
	teatest.RequireEqualOutput(t, finalOutput)
}
