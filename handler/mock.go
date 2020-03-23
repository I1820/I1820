package handler

import (
	"context"
	"fmt"

	"github.com/I1820/I1820/runner"
)

type MockedManager struct {
}

func (m MockedManager) New(_ context.Context, name string, _ []runner.Env) (runner.Runner, error) {
	return runner.Runner{
		ID:      fmt.Sprintf("dstn_%s", name),
		Port:    "1234",
		RedisID: fmt.Sprintf("rd_%s", name),
	}, nil
}

func (m MockedManager) Restart(context.Context, runner.Runner) error {
	return nil
}

func (m MockedManager) Remove(context.Context, runner.Runner) error {
	return nil
}

func (m MockedManager) Pull(context.Context) ([2]string, error) {
	return [2]string{"runner", "redis"}, nil
}
