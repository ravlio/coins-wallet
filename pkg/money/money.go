package money

const MoneyMultiplier = 1000000

type Money int

func (m *Money) SetFloat64(amount float64) {
	*m = Money(amount * float64(MoneyMultiplier))
}

func (m *Money) GetFloat64() float64 {
	return float64(*m / MoneyMultiplier)
}
