package money

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
)

type Currency string

const (
	MRU Currency = "MRU" // Mauritanian Ouguiya
	MRO Currency = "MRO" // Old Ouguiya (deprecated)
)

type Money struct {
	amount   decimal.Decimal
	currency Currency
}

func New(amount decimal.Decimal, currency Currency) Money {
	return Money{
		amount:   amount.Round(2),
		currency: currency,
	}
}

func FromFloat64(amount float64, currency Currency) Money {
	return New(decimal.NewFromFloat(amount), currency)
}

func FromString(amount string, currency Currency) (Money, error) {
	dec, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, fmt.Errorf("invalid amount: %w", err)
	}
	return New(dec, currency), nil
}

func FromCents(cents int64, currency Currency) Money {
	amount := decimal.NewFromInt(cents).Div(decimal.NewFromInt(100))
	return New(amount, currency)
}

// NewMRU creates a new MRU (Mauritanian Ouguiya) amount from cents
func NewMRU(cents int64) Money {
	return FromCents(cents, MRU)
}

func (m Money) Amount() decimal.Decimal { return m.amount }
func (m Money) Currency() Currency      { return m.currency }
func (m Money) String() string          { return fmt.Sprintf("%s %s", m.amount.StringFixed(2), m.currency) }
func (m Money) Cents() int64            { return m.amount.Mul(decimal.NewFromInt(100)).IntPart() }
func (m Money) Float64() float64        { f, _ := m.amount.Float64(); return f }
func (m Money) IsZero() bool            { return m.amount.IsZero() }
func (m Money) IsPositive() bool        { return m.amount.IsPositive() }
func (m Money) IsNegative() bool        { return m.amount.IsNegative() }

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch")
	}
	return New(m.amount.Add(other.amount), m.currency), nil
}

func (m Money) ToProviderAmount(inCents bool) string {
	if inCents {
		return fmt.Sprintf("%d", m.Cents())
	}
	return m.amount.StringFixed(2)
}

func (m Money) GetCurrencyCode() string {
	switch m.currency {
	case MRU:
		return "929"
	case MRO:
		return "478"
	default:
		return "929"
	}
}

func (m Money) Validate() error {
	if m.amount.IsNegative() {
		return fmt.Errorf("amount cannot be negative")
	}
	if m.currency == "" {
		return fmt.Errorf("currency required")
	}
	return nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"amount":   m.amount.String(),
		"currency": string(m.currency),
	})
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	amountStr, ok := obj["amount"].(string)
	if !ok {
		return fmt.Errorf("invalid amount")
	}

	currencyStr, ok := obj["currency"].(string)
	if !ok {
		return fmt.Errorf("invalid currency")
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return err
	}

	*m = New(amount, Currency(currencyStr))
	return nil
}
