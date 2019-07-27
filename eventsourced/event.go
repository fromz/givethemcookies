package eventsourced

// Event stores the data for every event
type Event struct {
	AggregateID   AggregateID `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Version       int         `json:"version"`
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
}
