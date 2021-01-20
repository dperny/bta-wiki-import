package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dperny/bta-wiki-import/export"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func makeFilename(desc export.Description) string {
	return fmt.Sprintf("%s.wiki", desc.Id)
}

var ExportCmd = &cobra.Command{
	Use:   "export <mod directory> <destination>",
	Short: "export all mod data to wikitext",
	RunE: func(cmd *cobra.Command, args []string) error {
		modDirectory := args[0]
		destination := args[1]
		mods, _ := export.WalkModsDirectory(modDirectory)

		for _, mod := range mods {
			for variant, mech := range mod.Mechs {
				blacklisted := false
				for _, tag := range mech.Mech.MechTags.Items {
					if tag == "BLACKLISTED" {
						blacklisted = true
					}
				}
				if blacklisted {
					logrus.Infof("Skipping BLACKLISTED mech %s", mech.Chassis.Description.Name)
					continue
				}

				filename := fmt.Sprintf("MechDef_%s_%s.wiki", mech.Chassis.Description.Name, variant)

				logrus.Infof("Writing mech wiki %s", filename)

				path := filepath.Join(destination, filename)

				file, err := os.Create(path)
				if err != nil {
					logrus.Errorf("Error opening %s: %s", path, err)
					continue
				}

				wiki := mech.Chassis.ToWiki() + mech.Mech.ToWiki()
				_, err = file.WriteString(wiki)
				if err != nil {
					logrus.Errorf("Error opening %s: %s", path, err)
				}
				file.Close()
			}

			for _, gear := range mod.Gear {
				blacklisted := false
				for _, tag := range gear.ComponentTags.Items {
					if tag == "BLACKLISTED" {
						blacklisted = true
						break
					}
				}

				if strings.HasPrefix(gear.Description.Id, "Gear_Quirk_") {
					// quirks are never blacklisted.
					blacklisted = false
				}
				if strings.Contains(mod.Mod, "BT Advanced") {
					// nothing in the core mods is blacklisted
					blacklisted = false
				}
				if strings.Contains(mod.Mod, "MechEngineer") {
					// nothing in MechEngineer is blacklisted
					blacklisted = false
				}

				if blacklisted {
					logrus.Infof("Skipping BLACKLISTED gear %s", gear.Description.Id)
					continue
				}
				filename := makeFilename(gear.Description)

				logrus.Infof("Writing gear wiki %s", filename)

				path := filepath.Join(destination, filename)

				file, err := os.Create(path)
				if err != nil {
					logrus.Errorf("Error opening %s: %s", path, err)
					continue
				}

				wiki := gear.ToWiki()
				_, err = file.WriteString(wiki)
				if err != nil {
					logrus.Errorf("Error writing %s: %s", path, err)
				}
				file.Close()
			}

			for _, weapon := range mod.Weapons {
				blacklisted := false
				for _, tag := range weapon.ComponentTags.Items {
					if tag == "BLACKLISTED" {
						blacklisted = true
					}
				}
				if blacklisted {
					logrus.Infof("Skipping BLACKLISTED weapon %s", weapon.Description.Id)
					continue
				}
				filename := makeFilename(weapon.Description)

				logrus.Infof("Writing weapon wiki %s", filename)

				path := filepath.Join(destination, filename)

				file, err := os.Create(path)
				if err != nil {
					logrus.Errorf("Error opening %s: %s", path, err)
					continue
				}

				wiki := weapon.ToWiki()
				_, err = file.WriteString(wiki)
				if err != nil {
					logrus.Errorf("Error writing %s: %s\n", path, err)
				}
				file.Close()
			}

			for _, jumpjet := range mod.JumpJets {
				blacklisted := false
				for _, tag := range jumpjet.ComponentTags.Items {
					if tag == "BLACKLISTED" {
						blacklisted = true
					}
				}
				if blacklisted {
					logrus.Infof("Skipping BLACKLISTED jumpjet %s", jumpjet.Description.Id)
					continue
				}
				filename := makeFilename(jumpjet.Description)

				logrus.Infof("Writing jumpjet  wiki %s", filename)

				path := filepath.Join(destination, filename)

				file, err := os.Create(path)
				if err != nil {
					logrus.Infof("Error opening %s: %s", path, err)
					continue
				}

				wiki := jumpjet.ToWiki()
				_, err = file.WriteString(wiki)
				if err != nil {
					logrus.Errorf("Error writing %s: %s", path, err)
				}
				file.Close()
			}

			for _, ammo := range mod.Ammo {
				blacklisted := false
				for _, tag := range ammo.AmmunitionBox.ComponentTags.Items {
					if tag == "BLACKLISTED" {
						logrus.Infof("Skipping BLACKLISTED ammo %s", ammo.AmmunitionBox.Description.Id)
						blacklisted = true
					}
				}
				if blacklisted {
					continue
				}
				filename := makeFilename(ammo.AmmunitionBox.Description)

				logrus.Infof("Writing ammo  wiki %s\n", filename)

				path := filepath.Join(destination, filename)

				file, err := os.Create(path)
				if err != nil {
					logrus.Errorf("Error opening %s: %s", path, err)
					continue
				}

				wiki := ammo.ToWiki()
				_, err = file.WriteString(wiki)
				if err != nil {
					logrus.Errorf("Error writing %s: %s", path, err)
				}
				file.Close()
			}
		}

		return nil
	},
}

var ExportMechCmd = &cobra.Command{
	Use:   "exportmech <mod directory> <mech variant>",
	Short: "export mod data to wikitext",
	RunE: func(cmd *cobra.Command, args []string) error {
		modDirectory := args[0]
		mechVariant := args[1]

		mods, _ := export.WalkModsDirectory(modDirectory)

		var (
			variant export.CompleteMechDef
			ok      bool
		)
		for _, mod := range mods {
			variant, ok = mod.Mechs[mechVariant]
			if ok {
				break
			}
		}

		if !ok {
			return fmt.Errorf("mech variant %s not found", mechVariant)
		}

		chassisWiki := variant.Chassis.ToWiki()
		mechWiki := variant.Mech.ToWiki()

		fmt.Print(chassisWiki + mechWiki)

		return nil
	},
}

var ExportGearCmd = &cobra.Command{
	Use:   "exportgear <file>",
	Short: "export gear data from the given file",
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		gear, err := export.ParseGear(file)
		if err != nil {
			return err
		}

		fmt.Println(gear.ToWiki())

		return nil
	},
}

/*
func init() {
	ExportCmd.PersistentFlags().String(&mechfile, "mechfile", "a file containing a list of mechs to export")
}
*/
