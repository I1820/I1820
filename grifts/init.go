package grifts

import (
	"github.com/aiotrc/pm/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
