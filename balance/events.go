package balance

import "github.com/ravlio/wallet/pkg/money"

type Event struct {
	AccountID uint32
	Currency  money.Currency
	Amount    money.Money
}

type DebitEvent struct {
	Event
}

type DebitedEvent struct {
	Event
}

type DebitFailEvent struct {
	Event
	Error error
}

type CreditEvent struct {
	Event
}

type CreditedEvent struct {
	Event
}

type CreditFailedEvent struct {
	Event
	Error error
}
