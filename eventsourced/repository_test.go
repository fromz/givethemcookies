package eventsourced_test

import (
	"testing"

	"github.com/fromz/givethemcookies/eventsourced"
)

func TestPopulateAggregateFromEventStore(t *testing.T) {
	aid := eventsourced.AggregateID("100")
	es := eventsourced.NewInMemoryEventStore()
	es.Save(eventsourced.Event{
		Version:       1,
		Type:          "entered",
		AggregateID:   aid,
		AggregateType: "product",
		Data: ProductEntered{
			Code: "sku",
		},
	})
	es.Save(eventsourced.Event{
		Version:       2,
		Type:          "product",
		AggregateID:   aid,
		AggregateType: "product",
		Data: DetailsSet{
			Description: "Fancy products long description",
			Title:       "Fancy product",
		},
	})
	product := NewProduct(aid)
	eventsourced.PopulateAggregateFromEventStore(es, &product)
	if product.Aggregate.Version != 2 {
		t.Errorf("Expected product version to be 2, got %d", product.Aggregate.Version)
		return
	}
	if product.Code == nil {
		t.Errorf("Expected product code to be set")
		return
	}
	if *product.Code != "sku" {
		t.Errorf("Expected code to be 'sku` but got %s", *product.Code)
		return
	}
	if product.Title == nil {
		t.Errorf("Expected product title to be set")
		return
	}
	if *product.Title != "Fancy product" {
		t.Errorf("Expected title to be 'Fancy product` but got %s", *product.Title)
		return
	}
	if product.Description == nil {
		t.Errorf("Expected product description to be set")
		return
	}
	if *product.Title != "Fancy product" {
		t.Errorf("Expected description to be 'Fancy products long description` but got %s", *product.Description)
		return
	}

}

func TestPersistAggregateChanges(t *testing.T) {
	aid := eventsourced.AggregateID("100")
	es := eventsourced.NewInMemoryEventStore()
	product := NewProduct(aid)
	product.HandleCommand(EnterProduct{
		Code: "sku",
	})
	product.HandleCommand(SetDetails{
		Title:       "Facny product",
		Description: "Fancy products long description",
	})
	if len(product.Aggregate.Changes) != 2 {
		t.Errorf("Expected 2 events to require persisting, got %d", len(product.Aggregate.Changes))
	}
	eventsourced.PersistAggregateChanges(es, product.Aggregate)
	evs, err := es.Load(aid)
	if err != nil {
		t.Error(err)
	}
	if len(evs) != 2 {
		t.Errorf("Expected 2 events, got %d", len(evs))
	}
	if len(product.Aggregate.Changes) != 0 {
		t.Errorf("Expected 0 events remaining on aggregate, got %d events", len(product.Aggregate.Changes))
	}
	if product.Aggregate.Version != 2 {
		t.Errorf("Expected aggregate to be at version 2, but its at %d", product.Aggregate.Version)
	}
}
