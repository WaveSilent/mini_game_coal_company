package equipment

import "company_coal_mine/company/miners"

type Equipment string

var Equipments = map[string]miners.Coal{
	"Pickaxes":             3000,
	"VentilationToTheMine": 15_000,
	"Trollyes":             50000,
}

type EquipmentInfo struct {
	Equipment string
	Cost      miners.Coal
	Bought    bool
}
