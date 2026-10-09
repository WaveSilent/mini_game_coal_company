package company

import "errors"

var ErrorNotFound error = errors.New("miner not found")

var ErrorNotFoundEquipment error = errors.New("equipment not found")

var ErrorEquipBought error = errors.New("equipment already bought")

var ErrorClass error = errors.New("class miner not found")

var ErrorNotEnoughCoal = errors.New("not enough coal")

var ErrorUnknownEquipment = errors.New("unknown equipment")
