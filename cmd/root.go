package cmd

import (
	"os"

	"github.com/I1820/link/cmd/server"
	"github.com/I1820/link/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ExitFailure status code
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	var root = &cobra.Command{
		Use:   "link",
		Short: "Who receives ingress data from LoRa server and many more",
	}
	root.Println("13 Feb 2020, Best Day Ever")

	server.Register(root, cfg)

	if err := root.Execute(); err != nil {
		logrus.Errorf("failed to execute root command: %s", err.Error())
		os.Exit(ExitFailure)
	}
}
