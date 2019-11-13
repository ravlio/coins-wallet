package event_bus

import (
	"sync"
	"sync/atomic"
)

type LocalBroker struct {
	mx   sync.RWMutex
	subs map[string][]Subscription
	lid  int32
}

var _ Broker = &LocalBroker{}

func NewLocalBroker() *LocalBroker {
	return &LocalBroker{mx: sync.RWMutex{}, subs: make(map[string][]Subscription)}
}

func (b *LocalBroker) Subscribe(name string, cb Cb) (Subscription, error) {
	sub := &LocalSubscription{
		id:    atomic.AddInt32(&b.lid, 1),
		event: name,
		Cb:    cb,
		b:     b,
	}

	b.mx.Lock()
	b.subs[name] = append(b.subs[name], sub)
	b.mx.Unlock()

	return sub, nil
}

func (b *LocalBroker) Publish(name string, payload interface{}, metadata map[string]interface{}) error {
	b.mx.RLock()
	defer b.mx.RUnlock()

	subs, ok := b.subs[name]
	if ok {
		for _, s := range subs {
			msg := &Message{
				Event:    name,
				Payload:  payload,
				Metadata: metadata,
			}
			s.Call(msg)
		}
	}

	return nil
}

func (b *LocalBroker) unsubscribe(s *LocalSubscription) {
	b.mx.Lock()
	defer b.mx.Unlock()

	subs, ok := b.subs[s.event]
	if ok {
		for k, sub := range subs {
			if sub.ID() == s.id {
				b.subs[s.event] = append(b.subs[s.event][:k], b.subs[s.event][k+1:]...)
				break
			}
		}
	}
}

type LocalSubscription struct {
	id    int32
	event string
	Cb    Cb
	b     *LocalBroker
}

func (l *LocalSubscription) Call(msg *Message) {
	l.Cb(msg)
}

func (l *LocalSubscription) ID() int32 {
	return l.id
}

func (l *LocalSubscription) Unsubscribe() error {
	l.b.unsubscribe(l)
	return nil
}
