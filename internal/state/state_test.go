package state

import (
	"sync"
	"testing"
)

func TestNewState(t *testing.T) {
	s := NewState()
	if s.IsCapturing() != false {
		t.Error("Expected isCapturing to be false on new state")
	}
	if s.TargetLanguage() != "en" {
		t.Errorf("Expected target language to be 'en' on new state, got '%s'", s.TargetLanguage())
	}
}

func TestSetAndGetCapturing(t *testing.T) {
	s := NewState()
	s.SetCapturing(true)
	if s.IsCapturing() != true {
		t.Error("Expected isCapturing to be true after setting")
	}
	s.SetCapturing(false)
	if s.IsCapturing() != false {
		t.Error("Expected isCapturing to be false after setting again")
	}
}

func TestSetAndGetTargetLanguage(t *testing.T) {
	s := NewState()
	lang := "fr"
	s.SetTargetLanguage(lang)
	if s.TargetLanguage() != lang {
		t.Errorf("Expected target language to be '%s', got '%s'", lang, s.TargetLanguage())
	}
}

func TestStateConcurrency(t *testing.T) {
	s := NewState()
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: sets capturing to true
	go func() {
		defer wg.Done()
		s.SetCapturing(true)
	}()

	// Goroutine 2: sets language to "de"
	go func() {
		defer wg.Done()
		s.SetTargetLanguage("de")
	}()

	wg.Wait()

	if !s.IsCapturing() {
		t.Error("Expected isCapturing to be true after concurrent set")
	}
	if s.TargetLanguage() != "de" {
		t.Errorf("Expected target language to be 'de' after concurrent set, got '%s'", s.TargetLanguage())
	}
}
