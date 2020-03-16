package link

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/I1820/I1820/config"
	"github.com/I1820/I1820/core"
	"github.com/I1820/I1820/db"
	"github.com/I1820/I1820/mqtt"
	"github.com/I1820/I1820/pkg/model/aolab"
	"github.com/I1820/I1820/pkg/protocol/lan"
	"github.com/I1820/I1820/pkg/protocol/lora"
	"github.com/I1820/I1820/rabbitmq"
	"github.com/I1820/I1820/store"
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

	// setup RabbitMQ producer
	rpr := rabbitmq.NewProducer(cfg.Rabbitmq, "raw")
	ppr := rabbitmq.NewProducer(cfg.Rabbitmq, "parsed")

	// creates the core application and registers the defaults
	core := core.New(tm, st, rpr, ppr)
	core.RegisterModel(aolab.Model{})
	if err := core.Run(); err != nil {
		logrus.Fatalf("Core Service failed with %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	core.Exit()
}

// Register link command
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