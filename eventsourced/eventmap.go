package eventsourced

// EventMap maps event names to a function which returns the events struct
type EventMap map[string]func(b []byte) (interface{}, error)
