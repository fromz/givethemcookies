package eventsourced

// AggregateID represents a unique ID for each aggregate
type AggregateID string

// AggregateType represents a type of aggregate, used in event store mapping etc
type AggregateType string

// Aggregate contains the basic info
// that all aggregates should have
type Aggregate struct {
	ID      AggregateID
	Type    AggregateType
	Version int
	Changes []Event
}

type AggregateRoot interface {
	GetAggregate() *Aggregate
	ApplyChange(event Event)
	HandleCommand(command interface{}) error
}
