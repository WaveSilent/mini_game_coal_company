package miners

import (
	"context"

	"github.com/google/uuid"
)

type (
	MinerType string
	Coal      int64
)

const (
	LittleMinerType   MinerType = "little"
	NormalMinerType   MinerType = "normal"
	PowerfulMinerType MinerType = "powerful"
)

var MinersSalaries = map[MinerType]Coal{
	LittleMinerType:   5,
	NormalMinerType:   50,
	PowerfulMinerType: 450,
}

type Miner interface {
	Run(ctx context.Context) <-chan Coal
	Info() MinerInfo
}

type MinerInfo struct {
	ID          uuid.UUID
	MinerType   MinerType
	MinerEnergy int64
	MinerPower  int64
}
