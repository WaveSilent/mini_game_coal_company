package miners

import (
	"context"

	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	LittleMinerSalary Coal = 5
)

type LittleMiner struct {
	id        uuid.UUID
	energy    *atomic.Int64
	power     int64
	timeBreak time.Duration
}

func NewLittleMiner() *LittleMiner {
	const (
		Energy    = 30
		Power     = 1
		TimeBreak = 3 * time.Second
	)

	energy := &atomic.Int64{}

	energy.Add(Energy)

	return &LittleMiner{
		id:        uuid.New(),
		energy:    energy,
		power:     Power,
		timeBreak: TimeBreak,
	}
}

func (lm *LittleMiner) Run(ctx context.Context) <-chan Coal {
	transferPoint := make(chan Coal)

	go func() {
		defer close(transferPoint)
		for {
			select {
			case <-ctx.Done():
				return

			case <-time.After(lm.timeBreak):
			}

			if lm.energy.Add(-1) < 0 {
				lm.energy.Add(1)
				return
			}

			select {
			case <-ctx.Done():
				return

			case transferPoint <- Coal(lm.power):
			}
		}
	}()

	return transferPoint
}

func (lm *LittleMiner) Info() MinerInfo {
	return MinerInfo{
		ID:          lm.id,
		MinerType:   LittleMinerType,
		MinerEnergy: lm.energy.Load(),
		MinerPower:  lm.power,
	}
}
