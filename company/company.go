package company

import (
	"company_coal_mine/company/equipment"
	"company_coal_mine/company/miners"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	balance         *atomic.Int64
	miners          map[uuid.UUID]miners.Miner // Активные
	hired           map[miners.MinerType]int   // Сколько всего нанято
	boughtEquipment map[string]bool

	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	mtx      sync.Mutex
	gameTime time.Time
}

func NewCompany(ctx context.Context) *Company {
	balance := &atomic.Int64{}

	ctxMain, contextCancel := context.WithCancel(ctx)

	return &Company{
		balance:         balance,
		miners:          make(map[uuid.UUID]miners.Miner),
		hired:           make(map[miners.MinerType]int),
		boughtEquipment: make(map[string]bool),

		wg:       sync.WaitGroup{},
		ctx:      ctxMain,
		cancel:   contextCancel,
		mtx:      sync.Mutex{},
		gameTime: time.Now(),
	}
}

func (c *Company) Hire(class miners.MinerType) (uuid.UUID, error) {
	var salary miners.Coal

	switch class {
	case miners.LittleMinerType:
		salary = miners.LittleMinerSalary
	case miners.NormalMinerType:
		salary = miners.NormalMinerSalary
	case miners.PowerfulMinerType:
		salary = miners.PowerfulMinerSalary
	default:
		return uuid.Nil, ErrorClass
	}

	if c.balance.Add(-int64(salary)) < 0 {
		c.balance.Add(int64(salary))
		return uuid.Nil, ErrorNotEnoughCoal
	}

	var miner miners.Miner
	switch class {
	case miners.LittleMinerType:
		miner = miners.NewLittleMiner()
	case miners.NormalMinerType:
		miner = miners.NewNormalMiner()
	case miners.PowerfulMinerType:
		miner = miners.NewPowerfulMiner()
	}

	info := miner.Info()

	c.mtx.Lock()
	c.miners[info.ID] = miner
	c.hired[info.MinerType]++
	c.mtx.Unlock()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		defer func() {
			c.mtx.Lock()
			delete(c.miners, info.ID)
			c.mtx.Unlock()
		}()

		for v := range miner.Run(c.ctx) {
			c.balance.Add(int64(v))
		}
	}()

	return info.ID, nil
}

func (c *Company) GetAllWorkedMiners() map[uuid.UUID]miners.MinerInfo {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	tmp := make(map[uuid.UUID]miners.MinerInfo, len(c.miners))

	for k, v := range c.miners {
		tmp[k] = v.Info()
	}

	return tmp
}

func (c *Company) GetInfoAllMiners() map[miners.MinerType]int {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	tmp := make(map[miners.MinerType]int, len(c.hired))

	for k, v := range c.hired {
		tmp[k] = v
	}

	return tmp
}

func (c *Company) BuyEquipment(equip string) (equipment.EquipmentInfo, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var cost miners.Coal
	cost, ok := equipment.Equipments[equip]
	if !ok {
		return equipment.EquipmentInfo{}, ErrorUnknownEquipment
	}

	if c.boughtEquipment[equip] {
		return equipment.EquipmentInfo{}, ErrorEquipBought
	}

	if c.balance.Add(-int64(cost)) < 0 {
		c.balance.Add(int64(cost))
		return equipment.EquipmentInfo{}, ErrorNotEnoughCoal
	}

	c.boughtEquipment[equip] = true

	return equipment.EquipmentInfo{
		Equipment: equip,
		Cost:      cost,
		Bought:    true,
	}, nil
}

func (c *Company) GetInfoEquipment() map[string]bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	tmp := make(map[string]bool, 3)

	for k, v := range c.boughtEquipment {
		tmp[k] = v
	}

	return tmp
}

func (c *Company) Stats() Statistics {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	minersTMP := make(map[uuid.UUID]miners.MinerInfo, len(c.miners))
	for k, v := range c.miners {
		minersTMP[k] = v.Info()
	}

	hiredTMP := make(map[miners.MinerType]int, len(c.hired))
	for k, v := range c.hired {
		hiredTMP[k] = v
	}

	return Statistics{
		Balance:      c.balance.Load(),
		ActiveMiners: minersTMP,
		HiredMiners:  hiredTMP,
		GameTime:     time.Since(c.gameTime).Seconds(),
	}
}

func (c *Company) Run() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(1 * time.Second):
				c.balance.Add(1)
			}
		}
	}()
}

func (c *Company) Stop() Statistics {
	c.cancel()
	c.wg.Wait()
	return c.Stats()
}
