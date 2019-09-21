package actions

import (
	"fmt"
	"math"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,

			"price": func(v float64) string {
				unit := 0.05
				v = math.Round(v/unit) * unit
				return fmt.Sprintf("%.2f", v)
			},
			"vat": func(v float64) string {
				return fmt.Sprintf("%.2f", v)
			},
		},
	})
}
