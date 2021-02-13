package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	flagColor bool
	flagDebug bool
)

var RootCmd = &cobra.Command{
	Use:   "btawiki",
	Short: "btawiki exports data from BTA3062 game files to mediawiki",
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if flagColor {
			logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		}
		if flagDebug {
			logrus.SetLevel(logrus.DebugLevel)
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&flagColor, "color", true, "enable color logging")
	RootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "enable debug-level logging")
}
