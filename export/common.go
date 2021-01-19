package export

import (
	"fmt"
)

const InventoryEquipmentTemplate = "MechInventory"

// Description is the golang construction of the json Description object found
// on most mechs and equipment.
type Description struct {
	Cost         int
	Rarity       int
	Purchasable  bool
	Manufacturer string
	Model        string
	UIName       string
	Id           string
	Name         string
	Details      string
	Icon         string
}

func (d Description) WikiArgs(wt *WikiTemplate) {
	wt.AddArg("Cost", fmt.Sprint(d.Cost))
	wt.AddArg("Rarity", fmt.Sprint(d.Rarity))
	wt.AddArg("Purchasable", fmt.Sprint(d.Purchasable))
	wt.AddArg("Manufacturer", d.Manufacturer)
	wt.AddArg("Model", d.Model)
	wt.AddArg("UIName", d.UIName)
	wt.AddArg("Id", d.Id)
	wt.AddArg("Name", d.Name)
	wt.AddArg("Details", d.Details)
	wt.AddArg("Icon", d.Icon)
}

// Tags is the golang construction of the Tags type found in many different
// json definitions.
type Tags struct {
	Items            []string `json:"items"`
	TagSetSourceFile string   `json:"tagSetSourceFile"`
}

// InventoryEquipment represents the a piece of equipment mounted on a mech or
// chassis.
type InventoryEquipment struct {
	MountedLocation  string
	ComponentDefID   string
	ComponentDefType string
	HardpointSlot    int
	IsFixed          bool
}

func (i InventoryEquipment) ToWiki(mechID string) string {
	wt := NewWikiTemplate(InventoryEquipmentTemplate)

	wt.AddArg("MechID", mechID)

	wt.AddArg("MountedLocation", i.MountedLocation)
	wt.AddArg("ComponentDefID", i.ComponentDefID)
	wt.AddArg("ComponentDefType", i.ComponentDefType)
	wt.AddArg("HardpointSlot", fmt.Sprint(i.HardpointSlot))
	wt.AddArg("FixedEquipment", fmt.Sprint(i.IsFixed))

	return wt.String()
}
