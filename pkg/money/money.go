package money

const Multiplier = 1000000

type Money int

func (m *Money) SetFloat64(amount float64) {
	*m = Money(amount * float64(Multiplier))
}

func (m *Money) GetFloat64() float64 {
	return float64(*m / Multiplier)
}

func Float64(amount float64) Money {
	m := new(Money)
	m.SetFloat64(amount)

	return *m
}
