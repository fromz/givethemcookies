package eventsourced

// NewInMemoryEventStore returns a new in memory event store
func NewInMemoryEventStore() InMmemoryEventStore {
	return InMmemoryEventStore{
		store: make(map[AggregateID][]Event),
	}
}

// InMmemoryEventStore is an in-memory based event store
// todo needs to be concurrency safe
type InMmemoryEventStore struct {
	store map[AggregateID][]Event
}

// Save persists an event store to memory
func (p InMmemoryEventStore) Save(event Event) error {
	p.store[event.AggregateID] = append(p.store[event.AggregateID], event)
	return nil
}

// Load returns events from memory based on aggregate id
func (p InMmemoryEventStore) Load(id AggregateID) ([]Event, error) {
	evs, ok := p.store[id]
	if !ok {
		return []Event{}, nil
	}
	return evs, nil
}
