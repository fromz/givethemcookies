package eventsourced

import (
	"database/sql"
	"encoding/json"
)

func NewPostgresEventStore(db *sql.DB, eventMap EventMap) PostgresEventStore {
	return PostgresEventStore{
		db: db,
		em: eventMap,
	}
}

type PostgresEventStore struct {
	db *sql.DB
	em EventMap
}

func (p PostgresEventStore) Save(event Event) error {
	d, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	sqlStatement := `
INSERT INTO events (aggregate_id, aggregate_type, version, type, data)
VALUES ($1, $2, $3, $4, $5)
`
	_, err = p.db.Exec(
		sqlStatement,
		event.AggregateID,
		event.AggregateType,
		event.Version,
		event.Type,
		string(d),
	)
	return err
}

func (p PostgresEventStore) Load(ID AggregateID) ([]Event, error) {
	evs := []Event{}
	rows, err := p.db.Query(`SELECT aggregate_id, aggregate_type, version, type, data FROM events where aggregate_id =$1`, ID)
	if err != nil {
		return []Event{}, err
	}
	defer rows.Close()
	for rows.Next() {
		ev := Event{}
		var d []byte
		err = rows.Scan(&ev.AggregateID, &ev.AggregateType, &ev.Version, &ev.Type, &d)
		if err != nil {
			return []Event{}, err
		}

		ed, err := p.em[ev.AggregateType](d)
		if err != nil {
			return []Event{}, err
		}
		ev.Data = ed
		evs = append(evs, ev)
	}
	err = rows.Err()
	if err != nil {
		return []Event{}, err
	}
	return evs, nil
}
