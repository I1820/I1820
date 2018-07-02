package actions

import (
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/suite"
)

type ActionSuite struct {
	*suite.Action
	buffalo.Logger
}

func Test_ActionSuite(t *testing.T) {
	as := &ActionSuite{
		suite.NewAction(App()),
		buffalo.NewLogger("info"),
	}
	suite.Run(t, as)
}
