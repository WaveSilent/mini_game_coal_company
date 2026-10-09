package miners

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	NormalMinerSalary Coal = 50
)

type NormalMiner struct {
	id        uuid.UUID
	energy    *atomic.Int64
	power     int64
	timeBreak time.Duration
}

func NewNormalMiner() *NormalMiner {
	const (
		Energy    = 45
		Power     = 3
		TimeBreak = 2 * time.Second
	)

	energy := &atomic.Int64{}

	energy.Add(Energy)

	return &NormalMiner{
		id:        uuid.New(),
		energy:    energy,
		power:     Power,
		timeBreak: TimeBreak,
	}
}

func (nm *NormalMiner) Run(ctx context.Context) <-chan Coal {
	transferPoint := make(chan Coal)

	go func() {
		defer close(transferPoint)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(nm.timeBreak):
			}

			if nm.energy.Add(-1) < 0 {
				nm.energy.Add(1)
				return
			}

			select {
			case <-ctx.Done():
				return
			case transferPoint <- Coal(nm.power):
			}
		}
	}()

	return transferPoint
}

func (nm *NormalMiner) Info() MinerInfo {
	return MinerInfo{
		ID:          nm.id,
		MinerType:   NormalMinerType,
		MinerEnergy: nm.energy.Load(),
		MinerPower:  nm.power,
	}
}
