package grifts

import (
	"github.com/jani123/book-keeping/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
