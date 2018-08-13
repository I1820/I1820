package grifts

import (
  "github.com/gobuffalo/buffalo"
	"github.com/I1820/dm/actions"
)

func init() {
  buffalo.Grifts(actions.App())
}
