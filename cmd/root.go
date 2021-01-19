package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "btawiki",
	Short: "btawiki exports data from BTA3062 game files to mediawiki",
}
