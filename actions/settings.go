package actions

import (
	"database/sql"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/jani123/book-keeping/models"
	"github.com/pkg/errors"
)

// SettingsGetHandler is handler for settings
func SettingsGetHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("settings.html"))
}

// SettingsPostHandler is handler for settings
func SettingsPostHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	company := c.Value("company").(*models.Company)

	if err := c.Bind(company); err != nil {
		return errors.Wrap(err, "could not bind parameters")
	}

	var err error
	var verr *validate.Errors

	if company.ID == uuid.Nil {
		id, err := uuid.FromString("2077c996-c9af-42b7-96c7-45db7274f41b")
		if err != nil {
			return errors.Wrap(err, "could not create uuid")
		}
		company.ID = id

		verr, err = tx.ValidateAndCreate(company)
	} else {
		verr, err = tx.ValidateAndSave(company)
	}

	if err != nil {
		c.Flash().Add("error", T.Translate(c, "company-could-not-be-created"))
		return errors.Wrap(err, "unable to validate and create company")
	} else if verr.HasAny() {
		for key, es := range verr.Errors {
			for _, e := range es {
				c.Flash().Add(key, T.Translate(c, e))
			}
		}
	}

	c.Flash().Add("success", T.Translate(c, "company-was-saved"))
	return c.Render(200, r.HTML("settings.html"))
}

// SettingsMiddleware loads various settings to context
func SettingsMiddleware() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			tx := c.Value("tx").(*pop.Connection)
			id, err := uuid.FromString("2077c996-c9af-42b7-96c7-45db7274f41b")
			if err != nil {
				return errors.Wrap(err, "could not create uuid")
			}

			company := models.Company{}

			err = tx.Find(&company, id)
			if errors.Cause(err) == sql.ErrNoRows {
				company.ID = uuid.Nil
				company.Name = "Example Company"
				company.Address = "Testers street 123"
				company.PostalCode = "12345"
				company.City = "Townville"
				company.BusinessID = "12345-1"
				company.SMS = "+3540 123 4567"
				company.Email = "email@example.org"
				company.BankAccount = "FI00 0000 0000 0000 00"
			} else if err != nil {
				return errors.Wrap(err, "could not find company")
			}

			c.Set("company", &company)

			return next(c)
		}
	}
}
