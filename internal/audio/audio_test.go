package audio

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestCalculateRMS(t *testing.T) {
	testCases := []struct {
		name     string
		samples  []int16
		expected float64
		epsilon  float64
	}{
		{
			name:     "Silence",
			samples:  []int16{0, 0, 0, 0},
			expected: 0.0,
			epsilon:  0.0001,
		},
		{
			name:     "Constant Signal Half Max",
			samples:  []int16{16384, 16384, 16384, 16384}, // 0.5 of max int16
			expected: 0.5,
			epsilon:  0.0001,
		},
		{
			name:     "Constant Signal Full Max",
			samples:  []int16{32767, 32767, 32767, 32767}, // Max int16
			expected: 1.0,
			epsilon:  0.0001,
		},
		{
			name:     "Alternating Signal",
			samples:  []int16{-16384, 16384, -16384, 16384},
			expected: 0.5,
			epsilon:  0.0001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert int16 samples to a byte slice (little-endian)
			buf := make([]byte, len(tc.samples)*2)
			for i, s := range tc.samples {
				binary.LittleEndian.PutUint16(buf[i*2:], uint16(s))
			}

			rms := CalculateRMS(buf)

			if math.Abs(rms-tc.expected) > tc.epsilon {
				t.Errorf("Expected RMS to be ~%f, but got %f", tc.expected, rms)
			}
		})
	}
}
