package export

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const ChassisDefWikiTemplate = "ChassisDef"

// ChassisDef is the golang construction of a chassisdef json object.
type ChassisDef struct {
	Custom struct {
		ArmActuatorSupport struct {
			LeftLimit  string
			RightLimit string
		}
	}
	Description        Description
	MovementCapDefID   string
	PathingCapDefID    string
	HardpointDataDefID string
	PrefabIdentifier   string
	PrefabBase         string

	Tonnage           float64
	InitialTonnage    float64
	WeightClass       string `json:"weightClass"`
	BattleValue       int
	Heatsinks         int
	TopSpeed          int
	TurnRadius        int
	MaxJumpjets       int
	Stability         int
	StabilityDefenses []int

	SpotterDistanceMultiplier float64
	VisibilityMultiplier      float64
	SensorRangeMultiplier     float64
	// TODO(rust dev): is Signature an int or a float?
	Signature          float64
	Radius             int
	PunchesWithLeftArm bool
	MeleeDamage        int
	MeleeInstability   int
	MeleeToHitModifier int
	DFADamage          int
	DFAToHitModifier   int
	DFASelfDamage      int
	DFAInstability     int
	Locations          []ChassisLocation

	VariantName    string
	ChassisTags    Tags
	StockRole      string
	YangsThoughts  string
	FixedEquipment []InventoryEquipment
}

type ChassisLocation struct {
	Location   string
	Hardpoints []struct {
		WeaponMount string
		Omni        bool
	}
	Tonnage           float64
	InventorySlots    int
	MaxArmor          int
	MaxRearArmor      int
	InternalStructure int
}

func (cl ChassisLocation) ToWiki(chassisID string) string {
	wt := NewWikiTemplate("ChassisLocation")

	wt.AddArg("ChassisID", chassisID)
	wt.AddArg("Location", cl.Location)
	wt.AddArg("Tonnage", fmt.Sprint(cl.Tonnage))
	wt.AddArg("InventorySlots", fmt.Sprint(cl.InventorySlots))
	wt.AddArg("InternalStructure", fmt.Sprint(cl.InternalStructure))
	wt.AddArg("MaxArmor", fmt.Sprint(cl.MaxArmor))
	wt.AddArg("MaxRearArmor", fmt.Sprint(cl.MaxRearArmor))

	hardpoints := []string{}
	omniHardpoints := []string{}
	for _, hardpoint := range cl.Hardpoints {
		if !hardpoint.Omni {
			hardpoints = append(hardpoints, hardpoint.WeaponMount)
		} else {
			omniHardpoints = append(omniHardpoints, hardpoint.WeaponMount)
		}
	}

	wt.AddArg("Hardpoints", strings.Join(hardpoints, ","))
	wt.AddArg("OmniHardpoints", strings.Join(omniHardpoints, ","))

	return wt.String()
}

func ParseChassisDef(data io.Reader) (ChassisDef, error) {
	var chassis ChassisDef

	d := json.NewDecoder(data)
	err := d.Decode(&chassis)

	if err == nil && chassis.Description.Id == "" {
		return chassis, fmt.Errorf("missing Id")
	}

	return chassis, err
}

func (cd ChassisDef) ToWiki() string {
	wt := NewWikiTemplate(ChassisDefWikiTemplate)

	cd.Description.WikiArgs(wt)

	wt.AddArg("Tonnage", fmt.Sprint(cd.Tonnage))
	wt.AddArg("InitialTonnage", fmt.Sprint(cd.InitialTonnage))
	wt.AddArg("weightClass", cd.WeightClass)
	wt.AddArg("VariantName", cd.VariantName)
	wt.AddArg("StockRole", cd.StockRole)
	wt.AddArg("YangsThoughts", cd.YangsThoughts)
	wt.AddArg("RightArmActuatorLimit", cd.Custom.ArmActuatorSupport.RightLimit)
	wt.AddArg("LeftArmActuatorLimit", cd.Custom.ArmActuatorSupport.LeftLimit)
	wt.AddArg("MeleeDamage", cd.MeleeDamage)
	wt.AddArg("MeleeInstability", cd.MeleeInstability)
	wt.AddArg("MeleeToHitModifier", cd.MeleeToHitModifier)

	locations := make([]string, len(cd.Locations))
	for i, location := range cd.Locations {
		locations[i] = location.ToWiki(cd.Description.Id)
	}

	inventory := make([]string, len(cd.FixedEquipment))
	equipmentDupes := map[InventoryEquipment]int{}
	for i, equipment := range cd.FixedEquipment {
		count := equipmentDupes[equipment]
		inventory[i] = equipment.ToWiki(cd.Description.Id, true, count)
		equipmentDupes[equipment] = count + 1
	}

	return wt.String() + strings.Join(locations, "") + strings.Join(inventory, "")
}
