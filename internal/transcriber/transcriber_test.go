package transcriber

import (
	"reflect"
	"testing"
)

func TestBytesToSamples(t *testing.T) {
	tests := []struct {
		name      string
		audioData []byte
		expected  []float32
	}{
		{
			name:      "Empty input",
			audioData: []byte{},
			expected:  nil,
		},
		{
			name:      "Single byte input (too small)",
			audioData: []byte{0x01},
			expected:  nil,
		},
		{
			name:      "Zero value sample",
			audioData: []byte{0x00, 0x00},
			expected:  []float32{0.0},
		},
		{
			name:      "Maximum positive int16",
			audioData: []byte{0xFF, 0x7F}, // 32767
			expected:  []float32{32767.0 / 32768.0},
		},
		{
			name:      "Minimum negative int16",
			audioData: []byte{0x00, 0x80}, // -32768
			expected:  []float32{-1.0},
		},
		{
			name:      "Two samples",
			audioData: []byte{0x00, 0x00, 0xFF, 0x7F},
			expected:  []float32{0.0, 32767.0 / 32768.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesToSamples(tt.audioData)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("BytesToSamples() = %v, want %v", got, tt.expected)
			}
		})
	}
}
