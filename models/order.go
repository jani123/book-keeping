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
	ID             uuid.UUID `json:"id" db:"id" form:"order-id"`
	Description    string    `json:"description" db:"description" form:"order-description"`
	Customer       string    `json:"customer" db:"customer" form:"order-customer"`
	Date           time.Time `json:"date" db:"date" form:"order-date"`
	Reference      string    `json:"reference" db:"reference" form:"order-reference"`
	PaymentMessage string    `json:"message" db:"payment_message" form:"order-payment-message"`
	DueDate        time.Time `json:"due-date" db:"due_date" form:"order-due-date"`
	CreatedAt      time.Time `json:"created-at" db:"created_at" form:"order-created-at"`
	UpdatedAt      time.Time `json:"updated-at" db:"updated_at" form:"order-updated-at"`

	Rows OrderRows `has_many:"order_rows" fk_id:"order_id" json:"-" db:"-" form:"-"`

	RowIDs       []uuid.UUID `json:"-" db:"-" form:"order-row-ids"`
	RowNames     []string    `json:"-" db:"-" form:"order-row-names"`
	RowPrices    []float64   `json:"-" db:"-" form:"order-row-prices"`
	RowVATs      []float64   `json:"-" db:"-" form:"order-row-vats"`
	RowQuantitys []int       `json:"-" db:"-" form:"order-row-quantitys"`
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
		&validators.StringIsPresent{Field: o.Reference, Name: "Reference", Message: "empty-order-reference"},
		&validators.StringIsPresent{Field: o.PaymentMessage, Name: "PaymentMessage", Message: "empty-order-payment-message"},
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
	} else if l != len(o.RowVATs) {
		return errors.New("row rates length differs from row ids length")
	} else if l != len(o.RowQuantitys) {
		return errors.New("row quantitys length differs from row ids length")
	}

	for i := 0; i < l; i++ {
		row := OrderRow{
			ID:       o.RowIDs[i],
			Order:    o.ID,
			Name:     o.RowNames[i],
			Price:    o.RowPrices[i],
			VAT:      o.RowVATs[i],
			Quantity: o.RowQuantitys[i],
		}
		rows = append(rows, row)
	}

	o.Rows = rows

	return nil
}

// TotalWithVAT sums up rows with VAT and returns total value of order
func (o *Order) TotalWithVAT() float64 {
	var sum float64 = 0.0
	for _, row := range o.Rows {
		sum += row.PriceWithVAT()
	}
	return sum
}

// TotalWithoutVAT sums up rows without VAT and returns total value of order
func (o *Order) TotalWithoutVAT() float64 {
	var sum float64 = 0.0
	for _, row := range o.Rows {
		sum += row.PriceWithoutVAT()
	}
	return sum
}

// TotalVAT returns total of VAT
func (o *Order) TotalVAT() float64 {
	var sum float64 = 0.0

	for _, row := range o.Rows {
		sum += row.PriceWithVAT() - row.PriceWithoutVAT()
	}

	return sum
}

// TotalVATPercent returns average VAT percent
func (o *Order) TotalVATPercent() float64 {
	var withVAT float64 = 0.0
	var withoutVAT float64 = 0.0

	for _, row := range o.Rows {
		withVAT += row.PriceWithVAT()
		withoutVAT += row.PriceWithoutVAT()
	}

	return (withVAT/withoutVAT - 1.0) * 100.0
}
