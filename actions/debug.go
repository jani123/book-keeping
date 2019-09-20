package actions

import "github.com/gobuffalo/buffalo"

// DebugHandler is a debug handler to serve up
// some debug info
func DebugHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("debug.html"))
}
