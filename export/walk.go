package export

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	ManifestTypeChassisDef    = "ChassisDef"
	ManifestTypeMechDef       = "MechDef"
	ManifestTypeHeatsink      = "HeatSinkDef"
	ManifestTypeUpgrade       = "UpgradeDef"
	ManifestTypeJumpJet       = "JumpJetDef"
	ManifestTypeWeapon        = "WeaponDef"
	ManifestTypeAmmunition    = "AmmunitionDef"
	ManifestTypeAmmunitionBox = "AmmunitionBoxDef"
)

type ModDef struct {
	Name        string
	Enabled     bool
	Hidden      bool
	Version     string
	Description string
	Manifest    []ModManifest
}

type ModManifest struct {
	Type string
	Path string
}

type CompleteMechDef struct {
	Chassis ChassisDef
	Mech    MechDef
}

type ModData struct {
	Mod      string
	Mechs    map[string]CompleteMechDef
	Gear     []Gear
	Weapons  []Weapon
	JumpJets []JumpJet
	Ammo     []CompleteAmmunition
}

func WalkMechs(modpath string, chassisdefPaths, mechdefPaths []string) (map[string]CompleteMechDef, []error) {
	mechs := map[string]CompleteMechDef{}

	chassisDefs := map[string]ChassisDef{}
	errors := []error{}

	for _, chassisdefPath := range chassisdefPaths {
		p := filepath.Join(modpath, chassisdefPath)
		chassisFiles, err := ioutil.ReadDir(p)
		if err != nil {
			logrus.Errorf("error reading chassis directory: %s", err)
			errors = append(errors, err)
		} else {
			for _, fileinfo := range chassisFiles {
				chassisPath := filepath.Join(p, fileinfo.Name())
				file, err := os.Open(chassisPath)
				if err != nil {
					logrus.Errorf("error opening %s: %s", fileinfo.Name(), err)
					errors = append(errors, err)
					continue
				}

				cd, err := ParseChassisDef(file)
				if err != nil {
					logrus.Errorf("error parsing %s: %s", fileinfo.Name(), err)
					errors = append(errors, err)
					continue
				}

				chassisDefs[cd.Description.Id] = cd
			}
		}
	}
	logrus.Infof("parsed %d chassisdefs", len(chassisDefs))

	for _, mechdefPath := range mechdefPaths {
		p := filepath.Join(modpath, mechdefPath)
		mechFiles, err := ioutil.ReadDir(p)
		if err != nil {
			logrus.Errorf("error reading mechdef directory")
		} else {
			for _, fileinfo := range mechFiles {
				mechPath := filepath.Join(p, fileinfo.Name())
				file, err := os.Open(mechPath)
				if err != nil {
					logrus.Errorf("error opening %s: %s", fileinfo.Name(), err)
					errors = append(errors, err)
					continue
				}

				md, err := ParseMechDef(file)
				if err != nil {
					logrus.Errorf("error parsing %s: %s", fileinfo.Name(), err)
					errors = append(errors, err)
					continue
				}

				// if the chassis did not parse correct or is not present,
				// storing the mechdef
				chassis, ok := chassisDefs[md.ChassisID]
				if ok {
					mechs[chassis.VariantName] = CompleteMechDef{
						Chassis: chassis,
						Mech:    md,
					}
				}
			}
		}
	}

	logrus.Infof("parsed %d mechdefs", len(mechs))

	return mechs, errors
}

func WalkGear(modpath string, gearPaths []string) ([]Gear, []error) {
	errors := []error{}
	var allGear []Gear
	for _, gearPath := range gearPaths {
		p := filepath.Join(modpath, gearPath)
		gearFiles, err := ioutil.ReadDir(p)
		if err != nil {
			logrus.Errorf("error reading gear directory %s", p)
			errors = append(errors, err)
			continue
		}

		for _, fileinfo := range gearFiles {
			f := filepath.Join(p, fileinfo.Name())
			file, err := os.Open(f)
			if err != nil {
				logrus.Errorf("error opening %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}

			gd, err := ParseGear(file)
			if err != nil {
				logrus.Errorf("error parsing %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}
			allGear = append(allGear, gd)
		}
	}

	logrus.Infof("parsed %d gear", len(allGear))
	return allGear, errors
}

func WalkJumpJets(modpath string, jumpjetPaths []string) ([]JumpJet, []error) {
	var (
		jumpjets []JumpJet
		errors   []error
	)

	for _, path := range jumpjetPaths {
		p := filepath.Join(modpath, path)
		jumpjetFiles, err := ioutil.ReadDir(p)
		if err != nil {
			logrus.Errorf("error reading jumpjet  directory %s", p)
			errors = append(errors, err)
			continue
		}

		for _, fileinfo := range jumpjetFiles {
			f := filepath.Join(p, fileinfo.Name())
			file, err := os.Open(f)
			if err != nil {
				logrus.Errorf("error opening %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}

			jd, err := ParseJumpJet(file)
			if err != nil {
				logrus.Errorf("error parsing %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}
			jumpjets = append(jumpjets, jd)
		}
	}

	logrus.Infof("parsed %d jumpjets", len(jumpjets))
	return jumpjets, errors
}

func WalkWeapons(modpath string, weaponPaths []string) ([]Weapon, []error) {
	var (
		errors  []error
		weapons []Weapon
	)

	for _, weaponPath := range weaponPaths {
		p := filepath.Join(modpath, weaponPath)
		weaponFiles, err := ioutil.ReadDir(p)
		if err != nil {
			logrus.Errorf("error reading weapon directory %s", p)
			errors = append(errors, err)
			continue
		}

		for _, fileinfo := range weaponFiles {
			f := filepath.Join(p, fileinfo.Name())
			file, err := os.Open(f)
			if err != nil {
				logrus.Errorf("error opening %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}

			w, err := ParseWeapon(file)
			if err != nil {
				logrus.Errorf("error parsing %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}
			weapons = append(weapons, w)
		}
	}

	logrus.Infof("parsed %d weapons", len(weapons))
	return weapons, errors
}

func WalkAmmunition(modpath string, ammunitionPaths, ammunitionBoxPaths []string) ([]CompleteAmmunition, []error) {
	var (
		errors       []error
		completeAmmo []CompleteAmmunition

		ammunitionCategories = map[string]string{}
	)

	for _, ammunitionPath := range ammunitionPaths {
		p := filepath.Join(modpath, ammunitionPath)
		ammoFiles, err := ioutil.ReadDir(p)
		if err != nil {
			logrus.Errorf("error reading Ammunition directory %s", p)
			errors = append(errors, err)
			continue
		}

		for _, fileinfo := range ammoFiles {
			f := filepath.Join(p, fileinfo.Name())

			file, err := os.Open(f)
			if err != nil {
				logrus.Errorf("error opening %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}

			ammo, err := ParseAmmunition(file)
			if err != nil {
				logrus.Errorf("error parsing %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}
			if ammo.Category == "" {
				if ammo.AmmoCategoryID != "" {
					ammo.Category = ammo.AmmoCategoryID
				} else {
					catErr := fmt.Errorf("ammo %s missing category", ammo.Description.Id)
					logrus.Errorf("%s", catErr)
					errors = append(errors, catErr)
				}
			}

			ammunitionCategories[ammo.Description.Id] = ammo.Category
		}
	}

	for _, ammoBoxPath := range ammunitionBoxPaths {
		p := filepath.Join(modpath, ammoBoxPath)
		ammoFiles, err := ioutil.ReadDir(p)
		if err != nil {
			logrus.Errorf("error reading AmmunitionBox directory %s", p)
			errors = append(errors, err)
			continue
		}

		for _, fileinfo := range ammoFiles {
			f := filepath.Join(p, fileinfo.Name())

			file, err := os.Open(f)
			if err != nil {
				logrus.Errorf("error opening %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}

			ammo, err := ParseAmmunitionBox(file)
			if err != nil {
				logrus.Errorf("error parsing %s: %s", fileinfo.Name(), err)
				errors = append(errors, err)
				continue
			}

			category, ok := ammunitionCategories[ammo.AmmoID]
			if !ok {
				catErr := fmt.Errorf("Cannot find category for ammo ID %s", ammo.AmmoID)
				logrus.Errorf("%s", catErr)
				errors = append(errors, err)
				continue
			}

			completeAmmo = append(completeAmmo, CompleteAmmunition{
				AmmunitionBox: ammo,
				Category:      category,
			})
		}
	}

	return completeAmmo, errors
}

func WalkMod(modpath string) (ModData, []error) {
	modfilePath := filepath.Join(modpath, "mod.json")

	modData := ModData{}

	errors := []error{}

	modfile, err := os.Open(modfilePath)
	if err != nil {
		logrus.Warnf("directory %s has no mod.json: %s", modpath, err)
		return modData, []error{err}
	}
	defer modfile.Close()

	var mod ModDef
	d := json.NewDecoder(modfile)
	err = d.Decode(&mod)
	if err != nil {
		logrus.Errorf("error parsing %s: %s", modfilePath, err)
		return modData, []error{err}
	}

	modData.Mod = mod.Name

	logrus.Infof("checking mod %q", mod.Name)

	var (
		chassisdefPaths, mechdefPaths                                 []string
		gearPaths, jumpjetPaths, weaponPaths, ammoPaths, ammoBoxPaths []string
	)

	for _, manifest := range mod.Manifest {
		switch manifest.Type {
		case ManifestTypeChassisDef:
			chassisdefPaths = append(chassisdefPaths, manifest.Path)
			logrus.Infof("mod defines chassisdefs at %s", manifest.Path)
		case ManifestTypeMechDef:
			mechdefPaths = append(mechdefPaths, manifest.Path)
			logrus.Infof("mod defines mechdefs at %s", manifest.Path)
		case ManifestTypeHeatsink, ManifestTypeUpgrade:
			gearPaths = append(gearPaths, manifest.Path)
			logrus.Infof("mod defines %s at %s", manifest.Type, manifest.Path)
		case ManifestTypeJumpJet:
			jumpjetPaths = append(jumpjetPaths, manifest.Path)
			logrus.Infof("mod define jump jets at %s", manifest.Type, manifest.Path)
		case ManifestTypeWeapon:
			weaponPaths = append(weaponPaths, manifest.Path)
			logrus.Infof("mod defines weapons at %s", manifest.Path)
		case ManifestTypeAmmunitionBox:
			ammoBoxPaths = append(ammoBoxPaths, manifest.Path)
			logrus.Infof("mod defines ammo box at %s", manifest.Path)
		case ManifestTypeAmmunition:
			ammoPaths = append(ammoPaths, manifest.Path)
			logrus.Infof("mod defines ammo at %s", manifest.Path)
		default:
			logrus.Warnf("ignoring unknown manifest type %s", manifest.Type)
		}
	}

	mechs, mechErrs := WalkMechs(modpath, chassisdefPaths, mechdefPaths)
	errors = append(errors, mechErrs...)

	gear, gearErrs := WalkGear(modpath, gearPaths)
	errors = append(errors, gearErrs...)

	jumpjets, jumpjetErrs := WalkJumpJets(modpath, jumpjetPaths)
	errors = append(errors, jumpjetErrs...)

	weapons, weaponErrs := WalkWeapons(modpath, weaponPaths)
	errors = append(errors, weaponErrs...)

	ammo, ammoErrs := WalkAmmunition(modpath, ammoPaths, ammoBoxPaths)
	errors = append(errors, ammoErrs...)

	modData.Mechs = mechs
	modData.Gear = gear
	modData.JumpJets = jumpjets
	modData.Weapons = weapons
	modData.Ammo = ammo

	return modData, errors
}

func WalkModsDirectory(path string) ([]ModData, []error) {
	// List the directory contents
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logrus.Errorf("error reading mods directory: %s", err)
		return nil, []error{err}
	}

	mods := []ModData{}
	allErrors := []error{}
	for _, file := range files {
		// skip any raw files. They don't have anything we care about.
		if !file.IsDir() {
			continue
		}

		// skip the .git directories, if present.
		if strings.Contains(file.Name(), ".git") {
			continue
		}

		modData, errors := WalkMod(filepath.Join(path, file.Name()))
		allErrors = append(allErrors, errors...)
		mods = append(mods, modData)
	}

	return mods, allErrors
}
