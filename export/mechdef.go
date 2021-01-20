package export

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const MechDefsTemplate = "MechDef"
const MechLocationTemplate = "MechLocation"

// MechDef is the golang construction of the mechdef json type.
type MechDef struct {
	MechTags            Tags
	ChassisID           string
	Description         Description
	SimGameMechPartCost int `json:"simGameMechPartCost"`
	Version             int
	Locations           []MechLocation
	Inventory           []InventoryEquipment `json:"inventory"`
}

type MechLocation struct {
	Location                 string
	CurrentArmor             int
	CurrentRearArmor         int
	CurrentInternalStructure int
	AssignedArmor            int
	AssignedRearArmor        int
}

func ParseMechDef(data io.Reader) (MechDef, error) {
	var mech MechDef
	d := json.NewDecoder(data)
	err := d.Decode(&mech)

	if err == nil && mech.Description.Id == "" {
		return mech, fmt.Errorf("missing Id")
	}
	return mech, err
}

func (l MechLocation) ToWiki(mechID string) string {
	wt := NewWikiTemplate(MechLocationTemplate)

	wt.AddArg("MechID", mechID)
	wt.AddArg("Location", l.Location)
	wt.AddArg("CurrentArmor", fmt.Sprint(l.CurrentArmor))
	wt.AddArg("CurrentRearArmor", fmt.Sprint(l.CurrentRearArmor))
	wt.AddArg("CurrentInternalStructure", fmt.Sprint(l.CurrentInternalStructure))
	wt.AddArg("AssignedArmor", fmt.Sprint(l.AssignedArmor))
	wt.AddArg("AssignedRearArmor", fmt.Sprint(l.AssignedRearArmor))

	return wt.String()
}

func (md MechDef) ToWiki() string {
	wt := NewWikiTemplate(MechDefsTemplate)

	md.Description.WikiArgs(wt)

	wt.AddArg("ChassisID", md.ChassisID)
	wt.AddArg("simGameMechPartCost", fmt.Sprint(md.SimGameMechPartCost))
	wt.AddArg("Version", fmt.Sprint(md.Version))

	locations := make([]string, len(md.Locations))
	for i, location := range md.Locations {
		locations[i] = location.ToWiki(md.Description.Id)
	}

	equipmentDupes := map[InventoryEquipment]int{}
	inventory := make([]string, len(md.Inventory))
	for i, equipment := range md.Inventory {
		count := equipmentDupes[equipment]
		inventory[i] = equipment.ToWiki(md.Description.Id, false, count)
		equipmentDupes[equipment] = count + 1
	}

	if len(md.MechTags.Items) > 0 {
		wt.AddArg("MechTags", strings.Join(md.MechTags.Items, ","))
	}

	return wt.String() + strings.Join(locations, "") + strings.Join(inventory, "")
}
