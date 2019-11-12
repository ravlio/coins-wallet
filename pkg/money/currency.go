package money

type Currency int

const (
	CurrencyUSD Currency = iota + 1
	CurrencyEUR
	CurrencyRUB
)

func IsCurrency(i Currency) bool {
	if i >= CurrencyUSD && i <= CurrencyRUB {
		return true
	}

	return false
}
