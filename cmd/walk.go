package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dperny/bta-wiki-import/export"
)

var WalkCommand = &cobra.Command{
	Use:   "walk <directory>",
	Short: "walks a given mod directory and explains what it finds",
	RunE: func(cmd *cobra.Command, args []string) error {
		modDirectory := args[0]

		mods, errors := export.WalkModsDirectory(modDirectory)

		var (
			mechCount    int
			gearCount    int
			weaponCount  int
			ammoCount    int
			jumpjetCount int
		)

		for _, mod := range mods {
			mechCount = mechCount + len(mod.Mechs)
			gearCount = gearCount + len(mod.Gear)
			weaponCount = weaponCount + len(mod.Weapons)
			ammoCount = ammoCount + len(mod.Ammo)
			jumpjetCount = jumpjetCount + len(mod.JumpJets)
		}

		fmt.Printf(
			"handled %d mechs, %d gear, %d jumpjets, %d weapons, %d ammunition with %d errors\n",
			mechCount, gearCount, jumpjetCount, weaponCount, ammoCount, len(errors),
		)

		return nil
	},
}
