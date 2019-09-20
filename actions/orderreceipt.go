package actions

import "github.com/gobuffalo/buffalo"

// OrderReceiptGetHandler is handler for order receipt
func OrderReceiptGetHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("receipt.html", "yield.html"))
}
