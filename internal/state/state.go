package state

import "sync"

// State holds the application's global state.
type State struct {
	mu                 sync.Mutex
	isCapturing        bool
	currentTargetLanguage string
}

// NewState creates a new State object.
func NewState() *State {
	return &State{
		isCapturing:        false,
		currentTargetLanguage: "en", // Default language
	}
}

// SetCapturing sets the capturing state.
func (s *State) SetCapturing(isCapturing bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isCapturing = isCapturing
}

// IsCapturing returns the current capturing state.
func (s *State) IsCapturing() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isCapturing
}

// SetTargetLanguage sets the target language for captions.
func (s *State) SetTargetLanguage(lang string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentTargetLanguage = lang
}

// TargetLanguage returns the current target language.
func (s *State) TargetLanguage() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.currentTargetLanguage
}
