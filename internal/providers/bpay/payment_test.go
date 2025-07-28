package bpay

import (
	"strconv"
	"testing"
)

func TestGeneratePasscode(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Generate passcode 1"},
		{"Generate passcode 2"},
		{"Generate passcode 3"},
		{"Generate passcode 4"},
		{"Generate passcode 5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passcode, err := generatePasscode()
			if err != nil {
				t.Errorf("generatePasscode() error = %v", err)
				return
			}

			// Check that passcode is exactly 4 digits
			if len(passcode) != 4 {
				t.Errorf("generatePasscode() length = %d, want 4", len(passcode))
			}

			// Check that passcode is numeric
			if _, err := strconv.Atoi(passcode); err != nil {
				t.Errorf("generatePasscode() = %v, want numeric string", passcode)
			}

			// Check that passcode is in valid range (1000-9999)
			code, _ := strconv.Atoi(passcode)
			if code < 1000 || code > 9999 {
				t.Errorf("generatePasscode() = %d, want between 1000-9999", code)
			}

			t.Logf("Generated passcode: %s", passcode)
		})
	}
}

func TestGeneratePasscodeUniqueness(t *testing.T) {
	// Generate multiple passcodes and check for basic distribution
	passcodes := make(map[string]int)
	iterations := 100

	for i := 0; i < iterations; i++ {
		passcode, err := generatePasscode()
		if err != nil {
			t.Errorf("generatePasscode() error = %v", err)
			return
		}
		passcodes[passcode]++
	}

	// We should have generated some variety (not all the same)
	if len(passcodes) == 1 {
		t.Error("generatePasscode() generated the same passcode every time, randomness may be broken")
	}

	t.Logf("Generated %d unique passcodes out of %d iterations", len(passcodes), iterations)
}
