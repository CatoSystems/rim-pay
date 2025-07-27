package money

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	amount := decimal.NewFromFloat(10.50)
	money := New(amount, MRU)

	assert.Equal(t, amount.Round(2), money.Amount())
	assert.Equal(t, MRU, money.Currency())
}

func TestFromFloat64(t *testing.T) {
	money := FromFloat64(10.50, MRU)

	assert.Equal(t, "10.5", money.Amount().String())
	assert.Equal(t, MRU, money.Currency())
}

func TestFromString(t *testing.T) {
	money, err := FromString("10.50", MRU)
	require.NoError(t, err)

	assert.Equal(t, "10.5", money.Amount().String())
	assert.Equal(t, MRU, money.Currency())
}

func TestFromStringInvalid(t *testing.T) {
	_, err := FromString("invalid", MRU)
	assert.Error(t, err)
}

func TestFromCents(t *testing.T) {
	money := FromCents(1050, MRU)

	assert.Equal(t, "10.5", money.Amount().String())
	assert.Equal(t, int64(1050), money.Cents())
}

func TestToProviderAmount(t *testing.T) {
	money := FromFloat64(10.50, MRU)

	assert.Equal(t, "10.50", money.ToProviderAmount(false))
	assert.Equal(t, "1050", money.ToProviderAmount(true))
}

func TestGetCurrencyCode(t *testing.T) {
	mruMoney := FromFloat64(10, MRU)
	assert.Equal(t, "929", mruMoney.GetCurrencyCode())

	mroMoney := FromFloat64(10, MRO)
	assert.Equal(t, "478", mroMoney.GetCurrencyCode())
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		money     Money
		wantError bool
	}{
		{
			name:      "valid positive amount",
			money:     FromFloat64(10.50, MRU),
			wantError: false,
		},
		{
			name:      "zero amount",
			money:     FromFloat64(0, MRU),
			wantError: false,
		},
		{
			name:      "negative amount",
			money:     FromFloat64(-10.50, MRU),
			wantError: true,
		},
		{
			name:      "empty currency",
			money:     Money{amount: decimal.NewFromFloat(10), currency: ""},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.money.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	money1 := FromFloat64(10.50, MRU)
	money2 := FromFloat64(5.25, MRU)

	result, err := money1.Add(money2)
	require.NoError(t, err)
	assert.Equal(t, "15.75", result.Amount().String())
}

func TestAddDifferentCurrencies(t *testing.T) {
	money1 := FromFloat64(10.50, MRU)
	money2 := FromFloat64(5.25, MRO)

	_, err := money1.Add(money2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currency mismatch")
}

func TestJSONMarshaling(t *testing.T) {
	money := FromFloat64(10.50, MRU)

	data, err := money.MarshalJSON()
	require.NoError(t, err)

	var result Money
	err = result.UnmarshalJSON(data)
	require.NoError(t, err)

	assert.True(t, money.Amount().Equal(result.Amount()))
	assert.Equal(t, money.Currency(), result.Currency())
}
