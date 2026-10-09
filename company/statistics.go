package company

import (
	"company_coal_mine/company/miners"

	"github.com/google/uuid"
)

type Statistics struct {
	Balance      int64                          `json:"balance"`
	ActiveMiners map[uuid.UUID]miners.MinerInfo `json:"active_miners"`
	HiredMiners  map[miners.MinerType]int       `json:"hired_miners"`
	GameTime     float64                        `json:"game_time"`
}
