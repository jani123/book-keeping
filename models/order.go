package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Order model struct
type Order struct {
	ID          uuid.UUID `json:"id" db:"id" form:"order-id"`
	Description string    `json:"description" db:"description" form:"order-description"`
	Customer    string    `json:"customer" db:"customer" form:"order-customer"`
	CreatedAt   time.Time `json:"created-at" db:"created_at" form:"order-created-at"`
	UpdatedAt   time.Time `json:"updated-at" db:"updated_at" form:"order-updated-at"`

	Rows OrderRows `has_many:"order_rows" fk_id:"order_id" json:"-" db:"-" form:"-"`

	RowIDs    []uuid.UUID `json:"row-ids" db:"-" form:"order-row-ids"`
	RowNames  []string    `json:"row-names" db:"-" form:"order-row-names"`
	RowPrices []float32   `json:"row-prices" db:"-" form:"order-row-prices"`
	RowRates  []float32   `json:"row-rates" db:"-" form:"order-row-rates"`
}

// String is not required by pop and may be deleted
func (o Order) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Orders is not required by pop and may be deleted
type Orders []Order

// String is not required by pop and may be deleted
func (o Orders) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *Order) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: o.Customer, Name: "Customer", Message: "empty-order-customer"},
		&validators.StringIsPresent{Field: o.Description, Name: "Description", Message: "empty-order-description"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *Order) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *Order) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// MakeRows converts Rows* member variables to Rows array
func (o *Order) MakeRows() error {
	rows := OrderRows{}
	l := len(o.RowIDs)

	if l != len(o.RowNames) {
		return errors.New("row names length differs from row ids length")
	} else if l != len(o.RowPrices) {
		return errors.New("row prices length differs from row ids length")
	} else if l != len(o.RowRates) {
		return errors.New("row rates length differs from row ids length")
	}

	for i := 0; i < l; i++ {
		row := OrderRow{
			ID:    o.RowIDs[i],
			Order: o.ID,
			Name:  o.RowNames[i],
			Price: o.RowPrices[i],
			Rate:  o.RowRates[i],
		}
		rows = append(rows, row)
	}

	o.Rows = rows

	return nil
}
