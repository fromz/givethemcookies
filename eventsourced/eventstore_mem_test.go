package eventsourced_test

import (
	"testing"

	"github.com/fromz/givethemcookies/eventsourced"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// all tests event store tests should be merged, using same tests against all implementations of EventStore

type InMemEventStoreEventTest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TestInMmemoryEventStore_Save tests save method works without errors
func TestInMmemoryEventStore_Save(t *testing.T) {
	e := eventsourced.Event{
		AggregateType: "test",
		Type:          "created",
		Version:       1,
		AggregateID:   eventsourced.AggregateID(uuid.New().String()),
		Data: InMemEventStoreEventTest{
			Name: "Roger",
			Age:  20,
		},
	}

	es := eventsourced.NewInMemoryEventStore()
	if err := es.Save(e); err != nil {
		t.Error(err)
	}

}

// TestInMmemoryEventStore_Load tests save method works without errors
func TestInMmemoryEventStore_Load(t *testing.T) {
	aid := eventsourced.AggregateID(uuid.New().String())
	e := eventsourced.Event{
		AggregateType: "test",
		Type:          "created",
		Version:       1,
		AggregateID:   aid,
		Data: InMemEventStoreEventTest{
			Name: "Roger",
			Age:  20,
		},
	}

	es := eventsourced.NewInMemoryEventStore()
	if err := es.Save(e); err != nil {
		t.Error(err)
		return
	}

	evs, err := es.Load(aid)
	if err != nil {
		t.Error(err)
		return
	}

	if len(evs) != 1 {
		t.Errorf("Expected 1 event, got %d", len(evs))
		return
	}

	for _, ev := range evs {
		if ev.AggregateID != aid {
			t.Errorf("Expected %s aggregate ID, got %s", aid, ev.AggregateID)
		}
		switch ty := ev.Data.(type) {
		case InMemEventStoreEventTest:
			if ty.Name != "Roger" {
				t.Errorf("Expected %s got %s", "Roger", ty.Name)
			}
			if ty.Age != 20 {
				t.Errorf("Expected %d got %d", 20, ty.Age)
			}
		default:
			t.Errorf("Expected EventStoreEventTest and it was: %s", ty)
		}
	}
}
