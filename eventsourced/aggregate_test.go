package eventsourced_test

import (
	"github.com/fromz/givethemcookies/eventsourced"
)

type Product struct {
	Aggregate   *eventsourced.Aggregate
	Code        *string
	Title       *string
	Description *string
	OnCatalogue bool
}

func (p *Product) GetAggregate() *eventsourced.Aggregate {
	return p.Aggregate
}

func NewProduct(ID eventsourced.AggregateID) Product {
	return Product{
		Aggregate: &eventsourced.Aggregate{
			ID:      ID,
			Type:    "product",
			Version: 0,
		},
	}
}

type EnterProduct struct {
	Code string
}

type SetDetails struct {
	Title       string
	Description string
}

type ProductEntered struct {
	Code string
}

type DetailsSet struct {
	Title       string
	Description string
}

func (p *Product) HandleCommand(command interface{}) error {
	event := eventsourced.Event{
		AggregateID:   p.Aggregate.ID,
		AggregateType: "Product",
	}

	switch c := command.(type) {
	case EnterProduct:
		event.Data = &ProductEntered{c.Code}

	case SetDetails:
		event.Data = &DetailsSet{
			c.Title,
			c.Description,
		}
	}
	p.Aggregate.Changes = append(p.Aggregate.Changes, event)
	return nil
}

func (p *Product) ApplyChange(event eventsourced.Event) {
	switch e := event.Data.(type) {
	case ProductEntered:
		p.Code = &e.Code
	case DetailsSet:
		p.Title = &e.Title
		p.Description = &e.Description
	default:
		panic("Don't know how")
	}
}
