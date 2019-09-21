package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// OrderRow model struct
type OrderRow struct {
	ID        uuid.UUID `json:"id" db:"id" form:"order-row-id"`
	Order     uuid.UUID `json:"order" db:"order_id" form:"order-row-order-id"`
	Name      string    `json:"name" db:"name" form:"order-row-name"`
	Price     float64   `json:"price" db:"price" form:"order-row-price"`
	VAT       float64   `json:"vat" db:"vat" form:"order-row-vat"`
	Quantity  int       `json:"quantity" db:"quantity" form:"order-row-quantity"`
	CreatedAt time.Time `json:"created-at" db:"created_at"`
	UpdatedAt time.Time `json:"updated-at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (o OrderRow) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// OrderRows is not required by pop and may be deleted
type OrderRows []OrderRow

// String is not required by pop and may be deleted
func (o OrderRows) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *OrderRow) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: o.Name, Name: "Name", Message: "empty-order-row-name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *OrderRow) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *OrderRow) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// UnitPriceWithoutVAT returns value of one item without VAT
func (o *OrderRow) UnitPriceWithoutVAT() float64 {
	return o.Price
}

// UnitPriceWithVAT returns value of one item with VAT
func (o *OrderRow) UnitPriceWithVAT() float64 {
	return o.UnitPriceWithoutVAT() * (1.0 + o.VAT/100.0)
}

// PriceWithoutVAT returns value of all items without VAT
func (o *OrderRow) PriceWithoutVAT() float64 {
	return o.UnitPriceWithoutVAT() * float64(o.Quantity)
}

// PriceWithVAT returns value of all items with VAT
func (o *OrderRow) PriceWithVAT() float64 {
	return o.UnitPriceWithVAT() * float64(o.Quantity)
}

// UnitVAT returns VAT value for one item
func (o *OrderRow) UnitVAT() float64 {
	return o.UnitPriceWithVAT() - o.UnitPriceWithoutVAT()
}

// TotalVAT returns VAT value for all items
func (o *OrderRow) TotalVAT() float64 {
	return o.PriceWithVAT() - o.PriceWithoutVAT()
}
