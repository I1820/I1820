package link

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/I1820/I1820/internal/config"
	"github.com/I1820/I1820/internal/core"
	"github.com/I1820/I1820/internal/db"
	"github.com/I1820/I1820/internal/mqtt"
	"github.com/I1820/I1820/internal/nats"
	"github.com/I1820/I1820/internal/store"
	"github.com/I1820/I1820/pkg/client/tm"
	"github.com/I1820/I1820/pkg/model/aolab"
	"github.com/I1820/I1820/pkg/protocol/lan"
	"github.com/I1820/I1820/pkg/protocol/lora"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	// create MQTT service for communicating with outer layers.
	msvc := mqtt.New(cfg.MQTT)

	// register protocols for gathering data.
	msvc.Register(lora.Protocol{})
	msvc.Register(lan.Protocol{})

	if err := msvc.Run(); err != nil {
		logrus.Fatalf("MQTT service failed with %s", err.Error())
	}

	// create a data store
	db, err := db.New(cfg.Database)
	if err != nil {
		logrus.Fatalf("Database failed with %s", err.Error())
	}

	st := store.New(db)

	// create a tm client for communicating with tm service.
	tm := tm.New(cfg.TM.URL)

	// setup NATS producer
	ns := nats.NewClient(cfg.NATS)

	// creates the core application and registers the defaults
	core := core.New(tm, st, ns)
	core.RegisterModel(aolab.Model{})

	if err := core.Run(); err != nil {
		logrus.Fatalf("Core Service failed with %s", err.Error())
	}

	go func() {
		for d := range msvc.Channel() {
			core.Handle(d)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	core.Exit()
}

// Register link command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "link",
			Short: "Who receives ingress data from LoRa server and many more",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
