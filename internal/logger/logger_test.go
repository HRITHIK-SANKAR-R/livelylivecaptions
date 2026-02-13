package logger

import (
	"livelylivecaptions/internal/types"
	"strings"
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	testCases := []struct {
		input    string
		expected LogLevel
	}{
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"WARN", LevelWarn},
		{"Error", LevelError},
		{"FATAL", LevelFatal},
		{"UNKNOWN", LevelInfo}, // Test default case
		{"", LevelInfo},        // Test empty string
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if got := parseLogLevel(tc.input); got != tc.expected {
				t.Errorf("parseLogLevel(%q) = %v; want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestRingBuffer(t *testing.T) {
	capacity := 3
	rb := NewRingBuffer(types.LogConfig{Level: "debug"}, capacity)

	// Test basic add and get
	rb.Add(LevelInfo, "message 1")
	rb.Add(LevelWarn, "message 2")

	entries := rb.GetEntries()
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}
	if entries[0].Message != "message 1" || entries[1].Message != "message 2" {
		t.Errorf("Unexpected entries: %+v", entries)
	}

	// Test wrapping around
	rb.Add(LevelDebug, "message 3")
	rb.Add(LevelError, "message 4") // This should overwrite "message 1"

	entries = rb.GetEntries()
	if len(entries) != capacity {
		t.Fatalf("Expected %d entries after wrapping, got %d", capacity, len(entries))
	}

	// The order should be message 2, message 3, message 4
	expectedMessages := []string{"message 2", "message 3", "message 4"}
	for i, entry := range entries {
		if entry.Message != expectedMessages[i] {
			t.Errorf("Expected message %q at index %d, got %q", expectedMessages[i], i, entry.Message)
		}
	}
}

func TestRingBufferLogLevelFiltering(t *testing.T) {
	rb := NewRingBuffer(types.LogConfig{Level: "info"}, 5)

	rb.Add(LevelDebug, "should be filtered")
	rb.Add(LevelInfo, "should be kept")
	rb.Add(LevelWarn, "should also be kept")

	entries := rb.GetEntries()
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries after filtering, got %d", len(entries))
	}

	if !strings.Contains(entries[0].Message, "should be kept") {
		t.Errorf("Expected first message to be kept, but it was not. Got: %s", entries[0].Message)
	}
	if !strings.Contains(entries[1].Message, "should also be kept") {
		t.Errorf("Expected second message to be kept, but it was not. Got: %s", entries[1].Message)
	}
}

func TestGlobalLogger(t *testing.T) {
	// Initialize global logger for this test
	InitGlobalLogger(types.LogConfig{Level: "info"})

	// It's hard to inspect the defaultLogger directly without exporting it,
	// so we'll test by observing behavior. Here we test if Fatal exits.
	// We can't truly test the logging without more complex setups (e.g., redirecting stdout),
	// but this covers the basic mechanism.

	// A simple test to ensure the functions can be called without panic
	Info("global info")
	Warn("global warn")

	// To test Fatal's os.Exit, you'd need a more advanced setup like
	// exec'ing the test with an environment variable. We'll skip that here.
	if defaultLogger == nil {
		t.Fatal("defaultLogger was not initialized")
	}

	entries := defaultLogger.GetEntries()
	if len(entries) != 2 {
		t.Errorf("Expected 2 log entries in global logger, got %d", len(entries))
	}
}
