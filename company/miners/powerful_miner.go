package miners

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	PowerfulMinerSalary = 450
)

type PowerfulMiner struct {
	id        uuid.UUID
	energy    *atomic.Int64
	power     *atomic.Int64
	timeBreak time.Duration
}

func NewPowerfulMiner() *PowerfulMiner {
	const (
		Energy    = 60
		Power     = 10
		TimeBreak = 1 * time.Second
	)

	energy := &atomic.Int64{}
	power := &atomic.Int64{}

	energy.Add(Energy)
	power.Add(Power)

	return &PowerfulMiner{
		id:        uuid.New(),
		energy:    energy,
		power:     power,
		timeBreak: TimeBreak,
	}
}

func (pm *PowerfulMiner) Run(ctx context.Context) <-chan Coal {
	transferPoint := make(chan Coal)

	go func() {
		defer close(transferPoint)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(pm.timeBreak):
			}

			if pm.energy.Add(-1) < 0 {
				pm.energy.Add(1)
				return
			}

			select {
			case <-ctx.Done():
				return
			case transferPoint <- Coal(pm.power.Load()):
				pm.power.Add(3)
			}
		}
	}()

	return transferPoint
}

func (pm *PowerfulMiner) Info() MinerInfo {
	return MinerInfo{
		ID:          pm.id,
		MinerType:   PowerfulMinerType,
		MinerEnergy: pm.energy.Load(),
		MinerPower:  pm.power.Load(),
	}
}
