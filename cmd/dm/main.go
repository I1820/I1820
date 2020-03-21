package dm

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/I1820/I1820/config"
	"github.com/I1820/I1820/db"
	"github.com/I1820/I1820/handler"
	"github.com/I1820/I1820/router"
	"github.com/I1820/I1820/store"
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

	qh := handler.Queries{
		Store: store.Data{
			DB: db,
		},
	}

	api := e.Group("/api")
	{
		qh.Register(api)
	}

	go func() {
		if err := e.Start(":1373"); err != http.ErrServerClosed {
			logrus.Fatalf("API Service failed with %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), ExitTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logrus.Printf("API Service failed on exit: %s", err)
	}
}

// Register dm command
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "dm",
			Short: "Who manages the data with Mongo and love",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
