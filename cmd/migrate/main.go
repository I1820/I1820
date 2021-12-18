package migrate

import (
	"context"
	"log"

	"github.com/I1820/I1820/config"
	"github.com/I1820/I1820/db"
	"github.com/I1820/I1820/store"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const enable = 1

func main(cfg config.Config) {
	db, err := db.New(cfg.Database)
	if err != nil {
		panic(err)
	}

	{
		idx, err := db.Collection(store.ThingCollection).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bson.M{"name": enable},
				Options: options.Index().SetUnique(true),
			})
		if err != nil {
			panic(err)
		}

		log.Println(idx)
	}

	{
		idx, err := db.Collection(store.ThingCollection).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Options: nil,
				Keys:    bson.M{"project": enable},
			})
		if err != nil {
			panic(err)
		}

		log.Println(idx)
	}

	{
		idx, err := db.Collection(store.DataCollection).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Options: nil,
				Keys: bson.M{
					"timestamp": -enable,
					"thing_id":  enable,
					"project":   enable,
				},
			})
		if err != nil {
			panic(err)
		}

		log.Println(idx)
	}
}

// Register migrate command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "migrate",
			Short: "Setup database indices",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
