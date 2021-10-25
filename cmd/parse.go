package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dperny/bta-wiki-import/export"

	"github.com/spf13/cobra"
)

var ParseCmd = &cobra.Command{
	Use:   "parse <type> <file>",
	Short: "parse a json file. for testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := args[0]
		filename := args[1]

		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		switch t {
		case "Weapon":
			weapon, err := export.ParseWeapon(file)
			if err != nil {
				return err
			}
			d, err := json.MarshalIndent(weapon, "", "\t")
			if err != nil {
				return err
			}
			fmt.Println(string(d))
		case "Gear":
			gear, err := export.ParseGear(file)
			if err != nil {
				return err
			}
			d, err := json.MarshalIndent(gear, "", "\t")
			if err != nil {
				return err
			}
			fmt.Println(string(d))
		case "Chassis":
			chassis, err := export.ParseChassisDef(file)
			if err != nil {
				return err
			}
			d, err := json.MarshalIndent(chassis, "", "\t")
			if err != nil {
				return err
			}
			fmt.Println(string(d))
		case "Mech":
			mech, err := export.ParseMechDef(file)
			if err != nil {
				return err
			}
			d, err := json.MarshalIndent(mech, "", "\t")
			if err != nil {
				return err
			}
			fmt.Println(string(d))
		}

		return nil
	},
}
