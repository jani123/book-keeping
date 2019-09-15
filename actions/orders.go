package actions

import (
	"fmt"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/jani123/book-keeping/models"
	"github.com/pkg/errors"
)

// OrdersFilters defines all the filters we need to filter orders.. o_O
type ordersFilters struct {
	StartDate time.Time `form:"orders-filter-start-date"`
	EndDate   time.Time `form:"orders-filter-end-date"`
}

// OrdersHandler lists the orders
func OrdersHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	filters := ordersFilters{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	}

	if err := c.Bind(&filters); err != nil {
		return errors.Wrap(err, "failed to bind arguments")
	}

	fmt.Printf("params: %+v\n", c.Params())
	fmt.Printf("filters(after bind): %+v\n", filters)

	// Reset times so we do not get results based on current time
	filters.StartDate = time.Date(
		filters.StartDate.Year(),
		filters.StartDate.Month(),
		filters.StartDate.Day(),
		0,
		0,
		0,
		0,
		filters.StartDate.Location(),
	)
	filters.EndDate = time.Date(
		filters.EndDate.Year(),
		filters.EndDate.Month(),
		filters.EndDate.Day(),
		0,
		0,
		0,
		0,
		filters.EndDate.Location(),
	)

	orders := models.Orders{}

	query := tx.Where("created_at >= ?", filters.StartDate)
	query = query.Where("created_at <= ?", filters.EndDate.AddDate(0, 0, 1)) // We want to include end of the day, meaning starting of next day
	err := query.All(&orders)

	if err != nil {
		return errors.Wrap(err, "failed to fetch all the orders")
	}

	fmt.Printf("filters: %+v\n", filters)
	c.Set("orders", orders)
	c.Set("filters", filters)
	return c.Render(200, r.HTML("orders.html"))
}
