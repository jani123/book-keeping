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
	Price     float32   `json:"price" db:"price" form:"order-row-price"`
	Rate      float32   `json:"rate" db:"rate" form:"order-row-rate"`
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
