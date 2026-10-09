package http

import (
	"company_coal_mine/company/miners"
	"encoding/json"
	"errors"
	"time"
)

type MinerDTO struct {
	MinerType miners.MinerType `json:"minerType"`
}

func (m MinerDTO) ValidateForCreate() error {
	if m.MinerType == "" {
		return errors.New("minerType is empty")
	}

	return nil
}

type EquipmentDTO struct {
	Equip string `json:"equip"`
}

func (e EquipmentDTO) ValidateForCreate() error {
	if e.Equip == "" {
		return errors.New("equipment is empty")
	}

	return nil
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func NewErrorDTO(err error) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)

}
