package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/I1820/I1820/runner"
)

type MockedManager struct {
	servers map[string]*httptest.Server
}

func (m *MockedManager) New(_ context.Context, name string, _ []runner.Env) (runner.Runner, error) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	if m.servers == nil {
		m.servers = make(map[string]*httptest.Server)
	}

	m.servers[fmt.Sprintf("dstn_%s", name)] = s

	return runner.Runner{
		ID:      fmt.Sprintf("dstn_%s", name),
		Port:    strings.Split(s.URL, ":")[1],
		RedisID: fmt.Sprintf("rd_%s", name),
	}, nil
}

func (m *MockedManager) Restart(context.Context, runner.Runner) error {
	return nil
}

func (m *MockedManager) Remove(_ context.Context, r runner.Runner) error {
	m.servers[r.ID].Close()

	delete(m.servers, r.ID)

	return nil
}

func (m *MockedManager) Pull(context.Context) ([2]string, error) {
	return [2]string{"runner", "redis"}, nil
}
