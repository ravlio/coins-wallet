package event_bus

type Cb func(msg *Message)

type Message struct {
	Event    string
	Payload  interface{}
	Metadata map[string]interface{}
}

type Broker interface {
	Publish(name string, payload interface{}, metadata map[string]interface{}) error
	Subscribe(name string, cb Cb) (Subscription, error)
}

type Subscription interface {
	Unsubscribe() error
	Call(msg *Message)
	ID() int32
}
