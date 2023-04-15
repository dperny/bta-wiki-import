package export

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	HeatsinkWikiTemplate        = "Heatsink"
	CoolingWikiTemplate         = "Cooling"
	EngineCoreWikiTemplate      = "EngineCore"
	EngineHeatBlockWikiTemplate = "EngineHeatBlock"
	EngineShieldWikiTemplate    = "EngineShield"
	JumpJetWikiTemplate         = "JumpJet"
	UpgradeWikiTemplate         = "Upgrade"
	WeaponWikiTemplate          = "Weapon"
	AmmoWikiTemplate            = "Ammunition"
)

type GearCustom struct {
	Category interface{}

	Weights struct {
		ReservedSlots int     `json:",omitempty"`
		EngineFactor  float64 `json:",omitempty"`
	} `json:",omitempty"`

	BonusDescriptions []string `json:",omitempty"`

	EngineHeatBlock struct {
		HeatSinkCount int `json:",omitempty"`
	} `json:",omitempty"`

	Cooling struct {
		HeatSinkDefId string `json:",omitempty"`
	} `json:",omitempty"`

	EngineCore struct {
		Rating int `json:",omitempty"`
	} `json:",omitempty"`

	ArmActuator struct {
		AccuracyBonus int    `json:",omitempty"`
		Type          string `json:",omitempty"`
	} `json:",omitempty"`

	DynamicSlots DynamicSlots `json:",omitempty"`

	AmmoCost struct {
		PerUnitCost int `json:",omitempty"`
	} `json:",omitempty"`
}

type DynamicSlots struct {
	ReservedSlots int
	ShowIcon      bool
	NameText      string
	BonusAText    string
}

type Gear struct {
	Custom              GearCustom
	Description         Description
	BonusValueA         string
	BonusValueB         string
	ComponentType       string
	ComponentSubType    string
	PrefabIdentifier    string
	BattleValue         int
	InventorySize       int
	Tonnage             float64
	AllowedLocations    string
	DisallowedLocations string
	// DissipationCapacity is only found on heatsinks
	DissipationCapacity int `json:",omitempty"`
	ComponentTags       struct {
		Items []string `json:"items"`
	}
}

type JumpJet struct {
	Gear
	JumpCapacity float64
	MinTonnage   int
	MaxTonnage   int
}

type Weapon struct {
	Gear
	Category                   string
	Type                       string
	MinRange                   int
	MaxRange                   int
	RangeSplit                 []int
	AmmoCategory               string
	StartingAmmoCapacity       int
	HeatGenerated              int
	Damage                     int
	OverheatedDamageMultiplier float64
	EvasiveDamageMultiplier    float64
	EvasivePipsIgnored         float64
	DamageVariance             float64
	HeatDamage                 int
	AccuracyModifier           float64
	CriticalChanceMultiplier   float64
	AOECapable                 bool
	IndirectFireCapable        bool
	RefireModifier             int
	ShotsWhenFired             int
	ProjectilesPerShot         int
	AttackRecoil               int
	Instability                int
	WeaponEffectID             string
}

type AmmunitionBox struct {
	Gear
	AmmoID   string
	Capacity int
}

type Ammunition struct {
	Description Description
	// Ammunition can store Category in two different places.
	Category       string
	AmmoCategoryID string
}

// CompleteAmmunition bundles up the AmmunitionBox, which contains most of the
// info we want, with the Category, which is found in the AmmunitionDef.
type CompleteAmmunition struct {
	AmmunitionBox AmmunitionBox
	Category      string
}

type Category struct {
	CategoryID string
}

func ParseWeapon(data io.Reader) (Weapon, error) {
	var weapon Weapon

	d := json.NewDecoder(data)
	err := d.Decode(&weapon)

	if err == nil && weapon.Description.Id == "" {
		return weapon, fmt.Errorf("missing Id")
	}

	return weapon, err
}

func ParseGear(data io.Reader) (Gear, error) {
	var gear Gear

	d := json.NewDecoder(data)
	err := d.Decode(&gear)

	if err == nil && gear.Description.Id == "" {
		return gear, fmt.Errorf("missing Id")
	}

	return gear, err
}

func ParseJumpJet(data io.Reader) (JumpJet, error) {
	var jumpjet JumpJet

	d := json.NewDecoder(data)
	err := d.Decode(&jumpjet)

	if err == nil && jumpjet.Description.Id == "" {
		return jumpjet, fmt.Errorf("missing Id")
	}

	return jumpjet, err
}

func ParseAmmunition(data io.Reader) (Ammunition, error) {
	var ammo Ammunition
	d := json.NewDecoder(data)
	err := d.Decode(&ammo)

	if err == nil && ammo.Description.Id == "" {
		return ammo, fmt.Errorf("missing Id")
	}

	return ammo, err
}

func ParseAmmunitionBox(data io.Reader) (AmmunitionBox, error) {
	var ammo AmmunitionBox
	d := json.NewDecoder(data)
	err := d.Decode(&ammo)

	if err == nil && ammo.Description.Id == "" {
		return ammo, fmt.Errorf("missing Id")
	}

	return ammo, err
}

func (w Weapon) ToWiki() string {
	wt := NewWikiTemplate(WeaponWikiTemplate)

	w.Description.WikiArgs(wt)

	wt.AddArg("Tonnage", w.Tonnage)
	wt.AddArg("InventorySize", w.InventorySize)
	wt.AddArg("BonusValueA", w.BonusValueA)
	wt.AddArg("ComponentType", w.ComponentType)
	wt.AddArg("ComponentSubType", w.ComponentSubType)
	wt.AddArg("BattleValue", w.BattleValue)
	wt.AddArg("AllowedLocations", w.AllowedLocations)
	wt.AddArg("DisallowedLocations", w.DisallowedLocations)
	wt.AddArg("Bonuses", strings.Join(w.Custom.BonusDescriptions, ","))

	wt.AddArg("Category", w.Category)
	wt.AddArg("Type", w.Type)
	wt.AddArg("MinRange", w.MinRange)
	wt.AddArg("MaxRange", w.MaxRange)
	rangeSplit := []string{}
	for _, bracket := range w.RangeSplit {
		rangeSplit = append(rangeSplit, fmt.Sprint(bracket))
	}
	wt.AddArg("RangeSplit", strings.Join(rangeSplit, ","))
	wt.AddArg("AmmoCategory", w.AmmoCategory)
	wt.AddArg("StartingAmmoCapacity", w.StartingAmmoCapacity)
	wt.AddArg("HeatGenerated", w.HeatGenerated)
	wt.AddArg("Damage", w.Damage)
	wt.AddArg("OverheatedDamageMultiplier", w.OverheatedDamageMultiplier)
	wt.AddArg("EvasiveDamageMultiplier", w.EvasiveDamageMultiplier)
	wt.AddArg("EvasivePipsIgnored", w.EvasivePipsIgnored)
	wt.AddArg("DamageVariance", w.DamageVariance)
	wt.AddArg("HeatDamage", w.HeatDamage)
	wt.AddArg("AccuracyModifier", w.AccuracyModifier)
	wt.AddArg("CriticalChanceMultiplier", w.CriticalChanceMultiplier)
	wt.AddArg("AOECapable", w.AOECapable)
	wt.AddArg("IndirectFireCapable", w.IndirectFireCapable)
	wt.AddArg("RefireModifier", w.RefireModifier)
	wt.AddArg("ShotsWhenFired", w.ShotsWhenFired)
	wt.AddArg("ProjectilesPerShot", w.ProjectilesPerShot)
	wt.AddArg("AttackRecoil", w.AttackRecoil)
	wt.AddArg("Instability", w.Instability)
	wt.AddArg("WeaponEffectID", w.WeaponEffectID)

	return wt.String()
}

func (a CompleteAmmunition) ToWiki() string {
	wt := NewWikiTemplate(AmmoWikiTemplate)

	a.AmmunitionBox.Description.WikiArgs(wt)

	wt.AddArg("Tonnage", a.AmmunitionBox.Tonnage)
	wt.AddArg("InventorySize", a.AmmunitionBox.InventorySize)
	wt.AddArg("BonusValueA", a.AmmunitionBox.BonusValueA)
	wt.AddArg("BonusValueB", a.AmmunitionBox.BonusValueA)
	wt.AddArg("ComponentType", a.AmmunitionBox.ComponentType)
	wt.AddArg("ComponentSubType", a.AmmunitionBox.ComponentSubType)
	wt.AddArg("BattleValue", a.AmmunitionBox.BattleValue)
	wt.AddArg("AllowedLocations", a.AmmunitionBox.AllowedLocations)
	wt.AddArg("DisallowedLocations", a.AmmunitionBox.DisallowedLocations)
	wt.AddArg("Bonuses", strings.Join(a.AmmunitionBox.Custom.BonusDescriptions, ","))

	wt.AddArg("AmmoID", a.AmmunitionBox.AmmoID)
	wt.AddArg("Capacity", a.AmmunitionBox.Capacity)
	wt.AddArg("Category", a.Category)
	wt.AddArg("PerUnitCost", a.AmmunitionBox.Custom.AmmoCost.PerUnitCost)

	return wt.String()
}

func (j JumpJet) ToWiki() string {
	wt := NewWikiTemplate(JumpJetWikiTemplate)

	j.Description.WikiArgs(wt)

	wt.AddArg("Tonnage", j.Tonnage)
	wt.AddArg("InventorySize", j.InventorySize)
	wt.AddArg("BonusValueA", j.BonusValueA)
	wt.AddArg("BonusValueB", j.BonusValueB)
	wt.AddArg("ComponentType", j.ComponentType)
	wt.AddArg("ComponentSubType", j.ComponentSubType)
	wt.AddArg("BattleValue", j.BattleValue)
	wt.AddArg("AllowedLocations", j.AllowedLocations)
	wt.AddArg("DisallowedLocations", j.DisallowedLocations)
	wt.AddArg("Bonuses", strings.Join(j.Custom.BonusDescriptions, ","))

	wt.AddArg("JumpCapacity", j.JumpCapacity)
	wt.AddArg("MinTonnage", j.MinTonnage)
	wt.AddArg("MaxTonnage", j.MaxTonnage)

	return wt.String()
}

func (g Gear) ToWiki() string {
	var templateType string
	// first, figure out what kind of wiki template this is. We can do that
	// by looking at Category.
	//
	// If the category is nil, that means
	categories := []string{}
	if g.Custom.Category == nil {
		switch g.ComponentType {
		case "HeatSink":
			templateType = HeatsinkWikiTemplate
		case "Upgrade":
			templateType = UpgradeWikiTemplate
		default:
			// TODO(rust dev): handle error
			return ""
		}
	} else {
		switch c := g.Custom.Category.(type) {
		case map[string]interface{}:
			category, ok := c["CategoryID"]
			if !ok {
				// TODO(rust dev): handle this case
				return ""
			}
			categoryString, ok := category.(string)
			if !ok {
				return ""
			}
			categories = append(categories, categoryString)
		case []interface{}:
			for _, iface := range c {
				cat, ok := iface.(map[string]interface{})
				if !ok {
					// TODO(rust dev): handle this case.
					return ""
				}
				c, ok := cat["CategoryID"]
				categoryString, ok := c.(string)
				if !ok {
					return ""
				}

				categories = append(categories, categoryString)
			}
		}

		for _, category := range categories {
			switch category {
			case "Cooling":
				templateType = CoolingWikiTemplate
			case "EngineCore":
				templateType = EngineCoreWikiTemplate
			case "EngineShield":
				templateType = EngineShieldWikiTemplate
			case "EngineHeatBlock":
				templateType = EngineHeatBlockWikiTemplate
			case "Heatsink":
				templateType = HeatsinkWikiTemplate
			}
		}

		if templateType == "" {
			templateType = UpgradeWikiTemplate
		}
	}

	// now that we have the templateType, let's make the wikitext
	wt := NewWikiTemplate(templateType)
	g.Description.WikiArgs(wt)

	wt.AddArg("Tonnage", g.Tonnage)
	wt.AddArg("InventorySize", g.InventorySize)
	wt.AddArg("BonusValueA", g.BonusValueA)
	wt.AddArg("BonusValueB", g.BonusValueB)
	wt.AddArg("ComponentType", g.ComponentType)
	wt.AddArg("ComponentSubType", g.ComponentSubType)
	wt.AddArg("BattleValue", g.BattleValue)
	wt.AddArg("Bonuses", strings.Join(g.Custom.BonusDescriptions, ","))
	wt.AddArg("CustomCategories", strings.Join(categories, ","))

	if g.ComponentType == "HeatSink" {
		wt.AddArg("Dissipation", g.DissipationCapacity)
	}

	switch templateType {
	case CoolingWikiTemplate:
		wt.AddArg("HeatsinkDefID", g.Custom.Cooling.HeatSinkDefId)
	case EngineHeatBlockWikiTemplate:
		wt.AddArg("HeatsinkCount", g.Custom.EngineHeatBlock.HeatSinkCount)
	case EngineCoreWikiTemplate:
		wt.AddArg("Rating", g.Custom.EngineCore.Rating)
	case EngineShieldWikiTemplate:
		wt.AddArg("ReservedSlots", g.Custom.Weights.ReservedSlots)
		wt.AddArg("EngineFactor", g.Custom.Weights.EngineFactor)
	}

	return wt.String()
}
