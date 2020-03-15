package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/I1820/link/actions"
	"github.com/I1820/link/config"
	"github.com/I1820/link/core"
	"github.com/I1820/link/mqtt"
	"github.com/I1820/link/pkg/model/aolab"
	"github.com/I1820/link/pkg/protocol/lan"
	"github.com/I1820/link/pkg/protocol/lora"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	fmt.Println("13 Feb 2020")

	// create mqtt service
	msvc := mqtt.New(cfg.MQTT)
	msvc.Register(lora.Protocol{})
	msvc.Register(lan.Protocol{})

	if err := msvc.Run(); err != nil {
		logrus.Fatalf("MQTT service failed with %s", err)
	}

	// creates the core application and registers the defaults
	core, err := core.New(cfg.TM.URL, cfg.Database.URL, cfg.Core.Broker.Addr)
	if err != nil {
		logrus.Fatal(err)
	}
	core.RegisterModel(aolab.Model{})
	if err := core.Run(); err != nil {
		logrus.Fatalf("Core Service failed with %s", err)
	}

	e := actions.App(cfg.Debug)
	go func() {
		if err := e.Start(":1378"); err != http.ErrServerClosed {
			logrus.Fatalf("API Service failed with %s", err)

		}
	}()

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
