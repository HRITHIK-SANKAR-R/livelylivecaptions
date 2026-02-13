package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"livelylivecaptions/internal/types"
)

const (
	width            = 120 // Fixed width for now, or adapt to window size
	height           = 20
	silenceThreshold = 0.01
	silenceDuration  = 5 * time.Second
)

var (
	// Styles
	meterStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			Width(10).
			Height(height)

	logStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			Width(width - 14).
			Height(height)

	finalTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))    // White
	partialTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))   // Grey
	warningTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))   // Orange
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	transcription  []string // History of final transcriptions
	partial        string   // Current partial transcription
	audioLevel     float64
	viewport       viewport.Model
	lastSoundTime  time.Time
	silenceWarning bool

	// Channels for receiving updates
	transChan <-chan types.TranscriptionEvent
	levelChan <-chan types.AudioLevelMsg
	quitChan  chan<- struct{}
}

func InitialModel(transChan <-chan types.TranscriptionEvent, levelChan <-chan types.AudioLevelMsg, quitChan chan<- struct{}) model {
	vp := viewport.New(width-16, height-2)
	vp.SetContent("Waiting for speech...")

	return model{
		transcription:  make([]string, 0),
		transChan:      transChan,
		levelChan:      levelChan,
		quitChan:       quitChan,
		viewport:       vp,
		lastSoundTime:  time.Now(),
		silenceWarning: false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForTranscription(m.transChan),
		waitForAudioLevel(m.levelChan),
		tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			close(m.quitChan)
			return m, tea.Quit
		}

	case types.TranscriptionEvent:
		if msg.IsFinal {
			m.transcription = append(m.transcription, msg.Text)
			m.partial = ""
		} else {
			m.partial = msg.Text
		}
		// We always update the viewport on a transcription event
		cmds = append(cmds, waitForTranscription(m.transChan))

	case types.AudioLevelMsg:
		m.audioLevel = float64(msg)
		if m.audioLevel > silenceThreshold {
			if m.silenceWarning {
				// If we were showing a warning, force a redraw now that sound is back
				m.lastSoundTime = time.Now()
				m.silenceWarning = false
			} else {
				m.lastSoundTime = time.Now()
			}
		}
		cmds = append(cmds, waitForAudioLevel(m.levelChan))

	case tickMsg:
		if time.Since(m.lastSoundTime) > silenceDuration {
			m.silenceWarning = true
		}
		// Always re-tick
		cmds = append(cmds, tickCmd())
	}

	// === Viewport Update Logic ===
	// This logic now runs on every message to keep the view consistent.
	var sb strings.Builder
	if m.silenceWarning {
		sb.WriteString(warningTextStyle.Render("Warning: No audio detected. Check microphone.\n\n"))
	}
	for _, line := range m.transcription {
		sb.WriteString(finalTextStyle.Render(line) + "\n")
	}
	if m.partial != "" {
		sb.WriteString(partialTextStyle.Render(m.partial))
	}
	m.viewport.SetContent(sb.String())
	m.viewport.GotoBottom()
	// ===========================

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Render Audio Meter - scale RMS to visible range
	// RMS is typically 0.0-0.3 for normal speech, so we amplify it
	scaledLevel := m.audioLevel * 3.0 // Amplify for visibility
	if scaledLevel > 1.0 {
		scaledLevel = 1.0
	}

	meterHeight := int(scaledLevel * float64(height-2))
	if meterHeight > height-2 {
		meterHeight = height - 2
	}

	// Build meter display from bottom up
	var meterLines []string
	meterLines = append(meterLines, "Level")
	for i := height - 3; i >= 0; i-- {
		if i >= (height-2)-meterHeight {
			meterLines = append(meterLines, "█████")
		} else {
			meterLines = append(meterLines, "     ")
		}
	}
	meterContent := strings.Join(meterLines, "\n")

	// Render Layout
	return lipgloss.JoinHorizontal(lipgloss.Top,
		meterStyle.Render(meterContent),
		logStyle.Render(m.viewport.View()),
	)
}

// Commands
func waitForTranscription(sub <-chan types.TranscriptionEvent) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func waitForAudioLevel(sub <-chan types.AudioLevelMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

// RunProgram starts the Bubble Tea program
func RunProgram(transChan <-chan types.TranscriptionEvent, levelChan <-chan types.AudioLevelMsg, quitChan chan<- struct{}) error {
	p := tea.NewProgram(InitialModel(transChan, levelChan, quitChan))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
