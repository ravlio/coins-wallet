package event_bus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocal(t *testing.T) {
	broker := NewLocalBroker()

	/*
		Create two variables
		Create three subscribers to same channels with callbacks incrementing vars
		Remove some subscribers
		Publish pointers to variables
		Check final values
		??
		Profit
	*/
	var a int
	var b int

	cb := func(msg *Message) {
		*(msg.Payload.(*int))++
	}

	_, err := broker.Subscribe("a", cb)
	require.NoError(t, err)

	_, err = broker.Subscribe("a", cb)
	require.NoError(t, err)

	s3, err := broker.Subscribe("a", cb)
	require.NoError(t, err)

	require.NoError(t, s3.Unsubscribe())

	_, err = broker.Subscribe("b", cb)
	require.NoError(t, err)

	s5, err := broker.Subscribe("b", cb)
	require.NoError(t, err)

	s6, err := broker.Subscribe("b", cb)
	require.NoError(t, err)

	require.NoError(t, s5.Unsubscribe())
	require.NoError(t, s6.Unsubscribe())

	require.NoError(t, broker.Publish("a", &a, nil))
	require.NoError(t, broker.Publish("b", &b, nil))

	require.Equal(t, 2, a)
	require.Equal(t, 1, b)

}
