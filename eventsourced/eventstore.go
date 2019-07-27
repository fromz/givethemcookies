package eventsourced

type EventStore interface {
	Save(event Event) error
	Load(ID AggregateID) ([]Event, error)
}
