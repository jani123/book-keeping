package actions

import (
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/jani123/book-keeping/models"
	"github.com/pkg/errors"
)

// OrderReceiptGetHandler is handler for order receipt
func OrderReceiptGetHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	id := c.Param("id")

	uuid, err := uuid.FromString(id)
	if err != nil {
		return errors.Wrap(err, "unable to parse id")
	}

	order := models.Order{}
	if err := tx.Find(&order, uuid); err != nil {
		return errors.Wrap(err, "unable to find order")
	}

	if err := tx.Load(&order, "Rows"); err != nil {
		return errors.Wrap(err, "unable to load order rows")
	}

	order.Date = time.Now()
	order.Reference = "12345"

	c.Set("order", order)

	return c.Render(200, r.HTML("receipt.html", "yield.html"))
}
