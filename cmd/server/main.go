package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/I1820/link/config"
	"github.com/I1820/link/core"
	"github.com/I1820/link/db"
	"github.com/I1820/link/mqtt"
	"github.com/I1820/link/pkg/model/aolab"
	"github.com/I1820/link/pkg/protocol/lan"
	"github.com/I1820/link/pkg/protocol/lora"
	"github.com/I1820/link/store"
	"github.com/I1820/tm/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	// create MQTT service
	msvc := mqtt.New(cfg.MQTT)
	msvc.Register(lora.Protocol{})
	msvc.Register(lan.Protocol{})

	if err := msvc.Run(); err != nil {
		logrus.Fatalf("MQTT service failed with %s", err)
	}

	// create a data store
	db, err := db.New(cfg.Database)
	if err != nil {
		logrus.Fatalf("Database failed with %s", err)
	}
	st := store.New(db)

	// create a tm service
	tm := client.New(cfg.TM.URL)

	// creates the core application and registers the defaults
	core := core.New(tm, st)
	core.RegisterModel(aolab.Model{})
	if err := core.Run(); err != nil {
		logrus.Fatalf("Core Service failed with %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	core.Exit()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logrus.Printf("API Service failed on exit: %s", err)
	}
}

// Register server command
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run server to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
