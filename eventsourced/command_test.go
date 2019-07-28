package eventsourced_test

import (
	"testing"

	"github.com/fromz/givethemcookies/eventsourced"
)

func TestCommandBus_HandleCommand(t *testing.T) {
	aid := eventsourced.AggregateID("100")
	es := eventsourced.NewInMemoryEventStore()
	cb := eventsourced.NewCommandBus(eventsourced.AggregateMap{
		"product": func(ID eventsourced.AggregateID) eventsourced.AggregateRoot {
			p := NewProduct(ID)
			return &p
		},
	}, es)
	cb.HandleCommand(eventsourced.CommandEnvelope{
		AggregateType: "product",
		AggregateID:   aid,
		Data:          EnterProduct{},
	})

	evs, err := es.Load(aid)
	if err != nil {
		t.Error(err)
	}
	if len(evs) != 1 {
		t.Errorf("Expected 1 event got %d", len(evs))
	}
	t.Error("Testing failed tests")
}
