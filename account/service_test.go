package account_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/commandhandler/bus"
	eventbus "github.com/looplab/eventhorizon/eventbus/local"
	eventstore "github.com/looplab/eventhorizon/eventstore/memory"
	repo "github.com/looplab/eventhorizon/repo/memory"
	"github.com/ravlio/wallet/account"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	store := eventstore.NewEventStore()

	// Create the event bus that distributes events.
	eventBus := eventbus.NewEventBus(nil)
	go func() {
		for e := range eventBus.Errors() {
			log.Printf("eventbus: %s", e.Error())
		}
	}()

	// Add a logger as an observer.

	svc, err := account.NewService(store, eventBus, bus.NewCommandHandler(), repo.NewRepo())
	require.NoError(t, err)

	id1 := uuid.New()
	err = svc.CreateAccount(context.Background(), &account.CreateAccountRequest{ID: id1, Name: "account1", Email: "account1@gmail.com"})
	require.NoError(t, err)

	id2 := uuid.New()
	err = svc.CreateAccount(context.Background(), &account.CreateAccountRequest{ID: id2, Name: "account2", Email: "account2@gmail.com"})
	require.NoError(t, err)

	notdound := uuid.New()
	_, err = svc.GetAccount(context.Background(), notdound)
	require.EqualError(t, err, eh.ErrEntityNotFound.Error())

	time.Sleep(time.Millisecond)

	acc, err := svc.GetAccount(context.Background(), id1)
	require.NoError(t, err)

	require.Equal(t, &account.Account{ID: id1, Name: "account1", Email: "account1@gmail.com", Version: 1}, acc)

	acc, err = svc.GetAccount(context.Background(), id2)
	require.NoError(t, err)

	require.Equal(t, &account.Account{ID: id2, Name: "account2", Email: "account2@gmail.com", Version: 1}, acc)
}
