package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"bta-wiki-import/importer"

	"github.com/spf13/cobra"
)

const (
	USERNAME_ENV = "WIKI_USER"
	PASSWORD_ENV = "WIKI_PASS"
)

var (
	flagDryRun       bool
	flagWikiUsername string
	flagWikiPassFile string
	flagWikiURL      string
)

var ImportCmd = &cobra.Command{
	Use:   "import <wikidata>",
	Short: "import mod data to wiki",
	RunE: func(cmd *cobra.Command, args []string) error {
		// first, check the flags for a username. Prefer this over the
		// environment variable.
		username := flagWikiUsername
		url := flagWikiURL

		// if there is no username flag set, then check the environment
		if flagWikiUsername == "" {
			if user := os.Getenv(USERNAME_ENV); user != "" {
				username = user
			} else if !flagDryRun {
				// if there is no username provided, we can still do a dry run
				// on the public wiki
				return fmt.Errorf("no wiki username provided")
			}
		}

		password := ""
		// again, first check to see if the password is in a file provided by
		// flags.
		passFile := flagWikiPassFile
		// if the passfile is not empty, open and read that file, trimming off
		// spaces
		if flagWikiPassFile != "" {
			fileContents, err := ioutil.ReadFile(passFile)
			if err != nil {
				return err
			}
			password = strings.TrimSpace(string(fileContents))
		} else {
			password = os.Getenv(PASSWORD_ENV)
		}

		if password == "" && !flagDryRun {
			return fmt.Errorf("no wiki password provided")
		}

		return importer.Import(args[0], flagDryRun, username, password, url)
	},
}

func init() {
	ImportCmd.Flags().BoolVarP(
		&flagDryRun, "dry-run", "d", false,
		"do a dry run, checking data but making no changes to the wiki",
	)
	ImportCmd.Flags().StringVarP(
		&flagWikiUsername, "username", "u", "",
		"the username to use when logging into the wiki",
	)
	ImportCmd.Flags().StringVarP(
		&flagWikiURL, "url", "l", "",
		"the wiki URL to log in against. Expects https://WEBSITE/api.php",
	)
	ImportCmd.Flags().StringVar(
		&flagWikiPassFile, "passfile", "",
		"a file to read the wiki password from",
	)
	// NEVER accept the password as a flag, which would leave the password in
	// the user's shell history.
}
