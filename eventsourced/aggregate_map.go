package eventsourced

// AggregateMap maps aggregate names to a function which returns an empty aggregate of that type
type AggregateMap map[string]func(id AggregateID) AggregateRoot
