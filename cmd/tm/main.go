package tm

import (
	"context"
	"log"
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

func main(cfg config.Config) {
	e := router.App()

	// routes
	db, err := db.New(cfg.Database)
	if err != nil {
		logrus.Fatal(err)
	}

	th := handler.Things{
		Store: store.Thing{
			DB: db,
		},
	}

	e.GET("/about", handler.AboutHandler)
	api := e.Group("/api")
	{
		th.Register(api)
	}

	go func() {
		if err := e.Start(":1995"); err != http.ErrServerClosed {
			logrus.Fatalf("API Service failed with %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("API Service failed on exit: %s", err)
	}
}

// Register tm command
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "tm",
			Short: "Who manages your things",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
