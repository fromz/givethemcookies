// +build integration

package eventsourced_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fromz/givethemcookies/eventsourced"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// all tests event store tests should be merged, using same tests against all implementations of EventStore

type EventStoreEventTest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func teardownDb(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, "test", "test", "eventstore")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Error(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		t.Error(err)
	}
	if err := m.Down(); err != nil {
		t.Error(err)
	}
}

func setupDb(t *testing.T) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, "test", "test", "eventstore")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Error(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		t.Error(err)
	}
	if err := m.Up(); err != nil {
		t.Error(err)
	}
	return db
}

// TestPostgresEventStore_Save tests save method works without errors, todo switch to mock db
func TestPostgresEventStore_Save_Integration(t *testing.T) {
	e := eventsourced.Event{
		AggregateType: "test",
		Type:          "created",
		Version:       1,
		AggregateID:   eventsourced.AggregateID(uuid.New().String()),
		Data: EventStoreEventTest{
			Name: "Roger",
			Age:  20,
		},
	}

	db := setupDb(t)

	es := eventsourced.NewPostgresEventStore(db, eventsourced.EventMap{
		"test": func(b []byte) (interface{}, error) {
			var ed EventStoreEventTest
			err := json.Unmarshal(b, &ed)
			return ed, err
		},
	})
	if err := es.Save(e); err != nil {
		t.Error(err)
	}

	teardownDb(t)
}

// TestPostgresEventStore_Save tests save method works without errors, todo switch to mock db
func TestPostgresEventStore_Load_Integration(t *testing.T) {
	aid := eventsourced.AggregateID(uuid.New().String())
	e := eventsourced.Event{
		AggregateType: "test",
		Type:          "created",
		Version:       1,
		AggregateID:   aid,
		Data: EventStoreEventTest{
			Name: "Roger",
			Age:  20,
		},
	}

	db := setupDb(t)

	es := eventsourced.NewPostgresEventStore(db, eventsourced.EventMap{
		"test": func(b []byte) (interface{}, error) {
			var ed EventStoreEventTest
			err := json.Unmarshal(b, &ed)
			return ed, err
		},
	})
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
		case EventStoreEventTest:
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

	teardownDb(t)
}
