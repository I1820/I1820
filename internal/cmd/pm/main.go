package pm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/I1820/I1820/internal/config"
	"github.com/I1820/I1820/internal/db"
	"github.com/I1820/I1820/internal/handler"
	"github.com/I1820/I1820/internal/router"
	"github.com/I1820/I1820/internal/runner"
	"github.com/I1820/I1820/internal/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	// ExitTimeout is a time that application waits for API service to exit
	ExitTimeout = 5 * time.Second
)

func main(cfg config.Config) {
	e := router.App()

	db, err := db.New(cfg.Database)
	if err != nil {
		logrus.Fatal(err)
	}

	m, err := runner.New()
	if err != nil {
		logrus.Fatal(err)
	}

	rh := handler.Runner{
		Store: store.Project{
			DB: db,
		},
		Manager:    m,
		DockerHost: cfg.Docker.Host,
	}

	api := e.Group("/api")
	{
		rh.Register(api)
	}

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", config.PMPort)); err != http.ErrServerClosed {
			logrus.Fatalf("API Service failed with %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), ExitTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("API Service failed on exit: %s", err)
	}
}

// Register tm command
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "pm",
			Short: "Who manages your project and their dockers",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
