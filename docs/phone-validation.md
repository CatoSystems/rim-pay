# Phone Number Validation

RimPay includes comprehensive phone number validation specifically designed for Mauritanian mobile numbers.

## Overview

The phone validation system ensures that only valid Mauritanian phone numbers are accepted for payment processing, preventing errors and improving success rates.

## Supported Formats

### International Format
```go
phone, err := phone.NewPhone("+22233445566")    // ✅ Valid
phone, err := phone.NewPhone("+222 33 44 55 66") // ✅ Valid (spaces ignored)
```

### National Format
```go
phone, err := phone.NewPhone("22233445566")     // ✅ Valid
phone, err := phone.NewPhone("222 33 44 55 66") // ✅ Valid (spaces ignored)
```

### Local Format
```go
phone, err := phone.NewPhone("33445566")        // ✅ Valid (assumes +222 prefix)
phone, err := phone.NewPhone("33 44 55 66")     // ✅ Valid (spaces ignored)
```

## Valid Prefixes

Mauritanian mobile operators use the following prefixes:

| Prefix | Operator | Description |
|--------|----------|-------------|
| **2** | Mauritel | National telecommunications operator |
| **3** | Chinguitel | Major mobile operator |
| **4** | Mattel | Mobile telecommunications operator |

## Validation Rules

### Length Requirements
- **With country code**: Exactly 11 digits (`+222` + 8 digits)
- **Without country code**: Exactly 8 digits
- **Local format**: Exactly 8 digits (prefix + 7 digits)

### Prefix Requirements
- Must start with 2, 3, or 4 (after country code)
- Invalid prefixes (5, 6, 7, 8, 9) are rejected

### Character Requirements
- Only numeric digits allowed (after normalization)
- Spaces, dashes, and parentheses are stripped during validation
- Country code +222 is automatically added if missing

## Usage Examples

### Basic Validation

```go
package main

import (
    "fmt"
    "github.com/CatoSystems/rim-pay/pkg/phone"
)

func main() {
    // Valid examples
    validNumbers := []string{
        "+22233445566",  // Chinguitel
        "+22222334455",  // Mauritel
        "+22244556677",  // Mattel
        "22233445566",   // Without + sign
        "33445566",      // Local format
    }

    for _, number := range validNumbers {
        phone, err := phone.NewPhone(number)
        if err != nil {
            fmt.Printf("❌ %s: %v\n", number, err)
        } else {
            fmt.Printf("✅ %s → %s\n", number, phone.String())
        }
    }
}
```

### Invalid Examples

```go
func demonstrateInvalidNumbers() {
    invalidNumbers := []string{
        "+22255667788",   // Invalid prefix 5
        "+22266778899",   // Invalid prefix 6
        "+2223344556",    // Too short
        "+222334455667",  // Too long
        "+221334455667",  // Wrong country code
        "1234567890",     // Wrong format
        "+222abcd5566",   // Contains letters
    }

    for _, number := range invalidNumbers {
        _, err := phone.NewPhone(number)
        if err != nil {
            fmt.Printf("❌ %s: %v\n", number, err)
        }
    }
}
```

## Phone Number Methods

### Basic Information

```go
phone, _ := phone.NewPhone("+22233445566")

// Get formatted string
fmt.Println(phone.String())           // "+22233445566"

// Get just the number part (without country code)
fmt.Println(phone.Number())           // "33445566"

// Check if valid
fmt.Println(phone.IsValid())          // true

// Get country code
fmt.Println(phone.CountryCode())      // "+222"
```

### Provider-Specific Formatting

```go
phone, _ := phone.NewPhone("+22233445566")

// Format for provider API (with country code)
fmt.Println(phone.ForProvider(true))  // "+22233445566"

// Format for provider API (without country code)
fmt.Println(phone.ForProvider(false)) // "33445566"

// Format for display to users
fmt.Println(phone.ForDisplay())       // "+222 33 44 55 66"
```

## Integration with Payment Requests

### B-PAY Integration

```go
// Create phone number
phone, err := phone.NewPhone("+22233445566")
if err != nil {
    return fmt.Errorf("invalid phone number: %w", err)
}

// Use in B-PAY request
request := &rimpay.BPayPaymentRequest{
    PhoneNumber: phone,
    Amount:      amount,
    Reference:   "ORDER-123",
    Passcode:    "1234",
}
```

### MASRVI Integration

```go
// Create phone number
phone, err := phone.NewPhone("+22244556677")
if err != nil {
    return fmt.Errorf("invalid phone number: %w", err)
}

// Use in MASRVI request
request := &rimpay.MasrviPaymentRequest{
    PhoneNumber: phone,
    Amount:      amount,
    Reference:   "ORDER-456",
    CallbackURL: "https://webhook.example.com",
}
```

## Error Handling

### Validation Errors

```go
phone, err := phone.NewPhone("invalid-number")
if err != nil {
    switch e := err.(type) {
    case *phone.ValidationError:
        fmt.Printf("Phone validation error: %s\n", e.Message)
        fmt.Printf("Invalid input: %s\n", e.Input)
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
}
```

### Common Error Messages

| Input | Error Message |
|-------|---------------|
| `"+22255667788"` | `invalid Mauritanian phone number: invalid prefix 5` |
| `"+2223344556"` | `invalid Mauritanian phone number: too short` |
| `"+222334455667"` | `invalid Mauritanian phone number: too long` |
| `"+221334455667"` | `invalid Mauritanian phone number: wrong country code` |
| `"abcd5566"` | `invalid phone number format: contains non-numeric characters` |

## Normalization Process

The phone validation system automatically normalizes input:

1. **Remove formatting**: Strips spaces, dashes, parentheses
2. **Add country code**: Adds +222 if missing
3. **Validate length**: Ensures correct digit count
4. **Validate prefix**: Checks operator prefix (2, 3, 4)
5. **Validate format**: Ensures only numeric characters

### Normalization Examples

```go
// All of these normalize to "+22233445566"
inputs := []string{
    "+222 33 44 55 66",
    "+222-33-44-55-66",
    "+222(33)445566",
    "222 33 44 55 66",
    "33445566",          // Assumes +222
}

for _, input := range inputs {
    phone, _ := phone.NewPhone(input)
    fmt.Printf("%s → %s\n", input, phone.String())
}
```

## Testing Phone Validation

### Unit Tests

```go
func TestPhoneValidation(t *testing.T) {
    tests := []struct {
        input    string
        expected string
        wantErr  bool
    }{
        {"+22233445566", "+22233445566", false},
        {"33445566", "+22233445566", false},
        {"+22255667788", "", true}, // Invalid prefix
        {"1234567", "", true},      // Too short
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            phone, err := phone.NewPhone(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, phone.String())
        })
    }
}
```

### Integration Tests

```go
func TestPhoneInPaymentRequest(t *testing.T) {
    phone, err := phone.NewPhone("+22233445566")
    require.NoError(t, err)

    request := &rimpay.BPayPaymentRequest{
        PhoneNumber: phone,
        Amount:      money.New(decimal.NewFromInt(1000), money.MRU),
        Reference:   "TEST-123",
        Passcode:    "1234",
    }

    err = request.Validate()
    assert.NoError(t, err)
}
```

## Best Practices

### 1. Input Validation
```go
// Always validate phone numbers from user input
func createPaymentFromUserInput(phoneInput string) error {
    phone, err := phone.NewPhone(phoneInput)
    if err != nil {
        return fmt.Errorf("please enter a valid Mauritanian phone number: %w", err)
    }
    
    // Use validated phone number
    request.PhoneNumber = phone
    return nil
}
```

### 2. Error Messages
```go
// Provide user-friendly error messages
func handlePhoneValidationError(err error) string {
    if validationErr, ok := err.(*phone.ValidationError); ok {
        switch {
        case strings.Contains(validationErr.Message, "invalid prefix"):
            return "Phone number must start with 2, 3, or 4 (after +222)"
        case strings.Contains(validationErr.Message, "too short"):
            return "Phone number is too short. Please enter 8 digits."
        case strings.Contains(validationErr.Message, "too long"):
            return "Phone number is too long. Please enter 8 digits."
        default:
            return "Please enter a valid Mauritanian phone number"
        }
    }
    return "Phone number validation failed"
}
```

### 3. Form Validation
```go
// Validate phone numbers in web forms
func validatePhoneField(phoneInput string) (string, error) {
    // Trim whitespace
    phoneInput = strings.TrimSpace(phoneInput)
    
    // Check if empty
    if phoneInput == "" {
        return "", errors.New("phone number is required")
    }
    
    // Validate format
    phone, err := phone.NewPhone(phoneInput)
    if err != nil {
        return "", handlePhoneValidationError(err)
    }
    
    // Return normalized format for display
    return phone.ForDisplay(), nil
}
```

## Performance Considerations

The phone validation is highly optimized:

- **O(1) complexity**: Constant time validation
- **No external dependencies**: Pure Go implementation
- **Memory efficient**: Minimal allocations
- **Thread safe**: Can be used concurrently

## Next Steps

- Learn about [money handling](money-handling.md)
- Explore [error handling](error-handling.md) patterns
- Check out [provider integration](providers/README.md) examples
