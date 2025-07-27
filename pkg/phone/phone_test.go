package phone

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMauritanianPhone(t *testing.T) {
	tests := []struct {
		name        string
		number      string
		expectedOp  Operator
		expectError bool
	}{
		{
			name:        "valid Mauritel number with +222",
			number:      "+22222334455",
			expectedOp:  Mauritel,
			expectError: false,
		},
		{
			name:        "valid Mattel number with 00222",
			number:      "0022266778899",
			expectedOp:  Mattel,
			expectError: false,
		},
		{
			name:        "valid Chinguitel number local format",
			number:      "88776655",
			expectedOp:  Chinguitel,
			expectError: false,
		},
		{
			name:        "invalid number - too short",
			number:      "1234567",
			expectError: true,
		},
		{
			name:        "invalid number - wrong prefix",
			number:      "12345678",
			expectError: true,
		},
		{
			name:        "empty number",
			number:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhone(tt.number)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, phone)
			} else {
				require.NoError(t, err)
				require.NotNil(t, phone)
				assert.Equal(t, tt.expectedOp, phone.Operator())
				assert.Len(t, phone.Number(), 8)
			}
		})
	}
}

func TestPhoneNumberFormats(t *testing.T) {
	phone, err := NewPhone("+22222334455")
	require.NoError(t, err)

	assert.Equal(t, "22334455", phone.LocalFormat())
	assert.Equal(t, "+22222334455", phone.String())
	assert.Equal(t, "+222 22 334 455", phone.InternationalFormat())
	assert.Equal(t, "22334455", phone.ForProvider(false))
	assert.Equal(t, "22222334455", phone.ForProvider(true))
}

func TestIsValidMauritanianNumber(t *testing.T) {
	tests := []struct {
		number string
		valid  bool
	}{
		{"+22222334455", true},
		{"0022266778899", true},
		{"88776655", true},
		{"222 22 33 44 55", true}, // with spaces
		{"12345678", false},       // invalid prefix
		{"1234567", false},        // too short
		{"+221123456789", false},  // wrong country code
		{"", false},               // empty
	}

	for _, tt := range tests {
		t.Run(tt.number, func(t *testing.T) {
			assert.Equal(t, tt.valid, IsValidMauritanianNumber(tt.number))
		})
	}
}

func TestDetermineOperator(t *testing.T) {
	tests := []struct {
		localNumber string
		expectedOp  Operator
	}{
		{"22334455", Mauritel},
		{"33445566", Mauritel},
		{"44556677", Mauritel},
		{"55667788", Mauritel},
		{"66778899", Mattel},
		{"77889900", Mattel},
		{"88990011", Chinguitel},
		{"99001122", Chinguitel},
		{"1234567", ""},  // invalid length
		{"12345678", ""}, // invalid prefix
	}

	for _, tt := range tests {
		t.Run(tt.localNumber, func(t *testing.T) {
			assert.Equal(t, tt.expectedOp, determineOperator(tt.localNumber))
		})
	}
}
