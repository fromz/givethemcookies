package eventsourced

func PopulateAggregateFromEventStore(es EventStore, root AggregateRoot) error {
	evs, err := es.Load(root.GetAggregate().ID)
	if err != nil {
		return err
	}
	for _, ev := range evs {
		root.ApplyChange(ev)
		root.GetAggregate().Version++
	}
	return nil
}

func PersistAggregateChanges(es EventStore, r *Aggregate) error {
	for _, c := range r.Changes {
		if err := es.Save(c); err != nil {
			return err
		}
		r.Version++
	}
	r.Changes = []Event{}
	return nil
}
