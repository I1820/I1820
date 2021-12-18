package cmd

import (
	"os"

	"github.com/I1820/I1820/internal/cmd/dm"
	"github.com/I1820/I1820/internal/cmd/link"
	"github.com/I1820/I1820/internal/cmd/migrate"
	"github.com/I1820/I1820/internal/cmd/pm"
	"github.com/I1820/I1820/internal/cmd/tm"
	"github.com/I1820/I1820/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "I1820",
		Short: "I1820, IoT Platform that can be used with LoRa, Zigbee and ...",
	}

	root.Println("13 Feb 2020, Best Day Ever")

	link.Register(root, cfg)
	tm.Register(root, cfg)
	dm.Register(root, cfg)
	pm.Register(root, cfg)
	migrate.Register(root, cfg)

	if err := root.Execute(); err != nil {
		logrus.Errorf("failed to execute root command: %s", err.Error())
		os.Exit(ExitFailure)
	}
}
