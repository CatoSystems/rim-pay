/*
Package phone provides Mauritanian phone number validation and formatting.

This package validates phone numbers according to Mauritanian numbering standards,
supporting the three main mobile operators with prefixes 2, 3, and 4.

# Usage

	import "github.com/CatoSystems/rim-pay/pkg/phone"

	// Create and validate phone number
	phone, err := phone.NewPhone("+22233445566")
	if err != nil {
		// Handle invalid phone number
	}

	fmt.Println(phone.String())           // +22233445566
	fmt.Println(phone.Number())           // 33445566
	fmt.Println(phone.IsValid())          // true

# Supported Formats

The following phone number formats are accepted:
  - +22233445566 (international format)
  - 22233445566 (national format)
  - 33445566 (local format - prefix 2, 3, or 4 assumed)

# Valid Prefixes

Mauritanian mobile numbers use the following prefixes:
  - 2: Mauritel (legacy)
  - 3: Chinguitel
  - 4: Mattel

# Validation Rules

Phone numbers must:
  - Start with country code +222 (optional in input)
  - Have a valid prefix (2, 3, or 4)
  - Be exactly 8 digits after the country code
  - Contain only numeric characters

Invalid examples:
  - +22255667788 (prefix 5 not supported)
  - +222334455 (too short)
  - +2223344556677 (too long)
  - +22233445abc (contains letters)
*/
package phone
