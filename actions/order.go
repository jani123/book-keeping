package actions

import (
	"fmt"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/jani123/book-keeping/models"
	"github.com/pkg/errors"
)

// NewOrderGetHandler handles get request for new order page
func NewOrderGetHandler(c buffalo.Context) error {
	c.Set("action", "/order/new")

	order := models.Order{}
	order.ID = uuid.Nil
	order.Customer = "Neighbour!"
	order.Description = "Big order!"
	order.Date = time.Now()
	order.Reference = "1234567"
	order.PaymentMessage = "You shall pay!"
	order.DueDate = time.Now().Add(time.Duration(3600))
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Rows = make(models.OrderRows, 2)
	order.Rows[0] = models.OrderRow{
		Name:     "Something small",
		Price:    0.35,
		VAT:      8.0,
		Quantity: 30,
	}
	order.Rows[1] = models.OrderRow{
		Name:     "Something else",
		Price:    35,
		VAT:      8.0,
		Quantity: 5,
	}
	c.Set("order", order)
	return c.Render(200, r.HTML("order.html"))
}

// NewOrderPostHandler handles post request for new order page
func NewOrderPostHandler(c buffalo.Context) error {
	order := models.Order{}
	c.Set("action", "/order")
	c.Set("order", &order)
	if err := c.Bind(&order); err != nil {
		return errors.Wrap(err, "unable to bind form parameters")
	}

	if err := order.MakeRows(); err != nil {
		return errors.Wrap(err, "unable to make order rows")
	}

	if verr, err := saveOrder(c, &order); err != nil {
		return errors.Wrap(err, "saveOrder failed")
	} else if verr != nil {
		return c.Render(400, r.HTML("order.html"))
	}

	c.Flash().Add("success", T.Translate(c, "order-was-created"))

	return c.Redirect(302, "orderPath()", render.Data{"id": order.ID})
}

// EditOrderGetHandler handles get request for edit order page
func EditOrderGetHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	id := c.Param("id")

	order := models.Order{}
	if err := tx.Find(&order, id); err != nil {
		return errors.Wrap(err, "unable to find order")
	}

	tx.Load(&order, "Rows")
	c.Set("order", order)
	c.Set("action", fmt.Sprintf("/order/%s", order.ID.String()))

	return c.Render(200, r.HTML("order.html"))
}

// EditOrderPostHandler handles post request for edit order page
func EditOrderPostHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	id := c.Param("id")
	c.Set("action", fmt.Sprintf("/order/%s", id))

	uuid, err := uuid.FromString(id)
	if err != nil {
		return errors.Wrap(err, "unable to parse id")
	}

	order := models.Order{}
	c.Set("order", &order)
	if err := tx.Find(&order, uuid); err != nil {
		return errors.Wrap(err, "unable to find order")
	}

	if err := c.Bind(&order); err != nil {
		return errors.Wrap(err, "unable to bind form parameters")
	}

	if order.ID != uuid {
		return errors.Wrap(err, "user is trying to hack the system")
	}

	if err := order.MakeRows(); err != nil {
		return errors.Wrap(err, "unable to make order rows")
	}

	if verr, err := saveOrder(c, &order); err != nil {
		return errors.Wrap(err, "saveOrder failed")
	} else if verr != nil {
		return c.Render(400, r.HTML("order.html"))
	}

	c.Flash().Add("success", T.Translate(c, "order-was-edited"))

	fmt.Printf("order: %v\n", order)
	return c.Render(200, r.HTML("order.html"))
}

func saveOrder(c buffalo.Context, order *models.Order) (error, error) {
	tx := c.Value("tx").(*pop.Connection)

	order.UpdatedAt = time.Now()

	var err error
	var verr *validate.Errors

	if order.ID == uuid.Nil {
		verr, err = tx.ValidateAndCreate(order)
	} else {
		verr, err = tx.ValidateAndSave(order)
	}

	if err != nil {
		return nil, errors.Wrap(err, "unable to validate and create order row")
	} else if verr.HasAny() {
		for key, es := range verr.Errors {
			for _, e := range es {
				c.Flash().Add(key, T.Translate(c, e))
			}
		}
		return errors.New("order"), nil
	}

	if err := tx.RawQuery("DELETE FROM order_rows WHERE order_id = ?", order.ID.String()).Exec(); err != nil {
		return nil, errors.Wrap(err, "could not delete order rows")
	}

	for _, row := range order.Rows {
		row.Order = order.ID
		row.UpdatedAt = time.Now()
		verr, err = tx.ValidateAndCreate(&row)

		if err != nil {
			return nil, errors.Wrap(err, "unable to validate and create order row")
		}

		if verr.HasAny() {
			for key, es := range verr.Errors {
				for _, e := range es {
					c.Flash().Add(key, T.Translate(c, e))
				}
			}
			return errors.New("order row"), nil
		}
	}

	return nil, nil
}
