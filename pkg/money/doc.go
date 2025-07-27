/*
Package money provides precise decimal-based money handling for financial applications.

This package uses the shopspring/decimal library to avoid floating-point precision
issues common in financial calculations. It supports the Mauritanian Ouguiya (MRU)
currency and provides convenient methods for money operations.

# Usage

	import (
		"github.com/CatoSystems/rim-pay/pkg/money"
		"github.com/shopspring/decimal"
	)

	// Create money amounts
	amount1 := money.New(decimal.NewFromInt(10050), "MRU")  // 100.50 MRU
	amount2 := money.FromFloat64(75.25, "MRU")              // 75.25 MRU

	// Arithmetic operations
	sum := amount1.Add(amount2)         // 175.75 MRU
	diff := amount1.Subtract(amount2)   // 25.25 MRU

	// Formatting and conversion
	fmt.Println(amount1.String())       // "100.50 MRU"
	fmt.Println(amount1.Cents())        // 10050 (amount in minor units)
	fmt.Println(amount1.Amount())       // 100.50 (decimal amount)

# Currency Support

Currently supports:
	- MRU: Mauritanian Ouguiya (current since January 1, 2018; 1 MRU = 100 cents)

Note: The old MRO currency was replaced by MRU on January 1, 2018, at a rate of 10 MRO = 1 MRU.

# Precision

All calculations use decimal arithmetic to maintain precision:
	- No floating-point rounding errors
	- Exact decimal representation
	- Safe for financial calculations

# Validation

Money amounts are validated to ensure:
	- Positive values only
	- Valid currency codes
	- Proper decimal precision (2 decimal places for MRU)

# Thread Safety

Money instances are immutable and thread-safe. All operations return new
Money instances rather than modifying existing ones.
*/
package money
