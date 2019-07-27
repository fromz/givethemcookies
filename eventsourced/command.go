package eventsourced

import (
	"fmt"
)

// CommandEnvelope contains the basic info
// that all commands should have
type CommandEnvelope struct {
	AggregateID   AggregateID
	AggregateType string
	Data          interface{}
}

type CommandBus struct {
	am AggregateMap
	es EventStore
}

func NewCommandBus(am AggregateMap, es EventStore) CommandBus {
	return CommandBus{
		am: am,
		es: es,
	}
}

// HandleCommand needs to retry when an event versioning conflict occurs
func (c *CommandBus) HandleCommand(ce CommandEnvelope) error {
	f, ok := c.am[ce.AggregateType]
	if !ok {
		return fmt.Errorf("no mapped aggregate for type %s", ce.AggregateType)
	}
	a := f(ce.AggregateID)
	err := PopulateAggregateFromEventStore(c.es, a)
	if err != nil {
		return err
	}
	if err := a.HandleCommand(ce.Data); err != nil {
		return err
	}
	if err := PersistAggregateChanges(c.es, a.GetAggregate()); err != nil {
		return err
	}
	return nil
}
