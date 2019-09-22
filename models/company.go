package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Company model struct
type Company struct {
	ID          uuid.UUID `json:"id" db:"id" form:"order-id"`
	Name        string    `json:"name" db:"name" form:"company-name"`
	Address     string    `json:"address" db:"address" form:"company-address"`
	City        string    `json:"city" db:"city" form:"company-city"`
	BusinessID  string    `json:"business-id" db:"business_id" form:"company-business-id"`
	SMS         string    `json:"sms" db:"sms" form:"company-sms"`
	Email       string    `json:"email" db:"email" form:"company-email"`
	BankAccount string    `json:"bank-account" db:"bank_account" form:"company-bank-account"`
	CreatedAt   time.Time `json:"created-at" db:"created_at" form:"order-created-at"`
	UpdatedAt   time.Time `json:"updated-at" db:"updated_at" form:"order-updated-at"`
}

// String is not required by pop and may be deleted
func (o Company) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *Company) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: o.Name, Name: "Name", Message: "empty-company-name"},
		&validators.StringIsPresent{Field: o.Address, Name: "Address", Message: "empty-company-address"},
		&validators.StringIsPresent{Field: o.City, Name: "City", Message: "empty-company-city"},
		&validators.StringIsPresent{Field: o.BusinessID, Name: "BusinessID", Message: "empty-company-business-id"},
		&validators.StringIsPresent{Field: o.SMS, Name: "SMS", Message: "empty-sms"},
		&validators.StringIsPresent{Field: o.Email, Name: "Email", Message: "empty-email"},
		&validators.StringIsPresent{Field: o.BankAccount, Name: "BankAccount", Message: "empty-bank-account"},
	), nil
}
