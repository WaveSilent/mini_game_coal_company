package http

import (
	"company_coal_mine/company"
	"company_coal_mine/company/equipment"
	"company_coal_mine/company/miners"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type HTTPHandler struct {
	company *company.Company
}

func NewHTTPHandler(company *company.Company) *HTTPHandler {
	return &HTTPHandler{
		company: company,
	}
}

func writeJSON(w http.ResponseWriter, val any) {
	b, err := json.MarshalIndent(val, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("error response body", err.Error())
		return
	}
}

/*
pattern: /miners
method: POST
info: JSON request body with need Miner

success:
response body: JSON body created Miner
status: 201 Created

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerHireMiner(w http.ResponseWriter, r *http.Request) {
	var minerDTO MinerDTO
	if err := json.NewDecoder(r.Body).Decode(&minerDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := minerDTO.ValidateForCreate(); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	id, err := h.company.Hire(minerDTO.MinerType)

	if err != nil {
		errDTO := NewErrorDTO(err)
		if errors.Is(err, company.ErrorClass) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else if errors.Is(err, company.ErrorNotEnoughCoal) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)

		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	writeJSON(w, id)
}

/*
pattern: /miners
method: GET
info: query param

success:
response body: JSON body represent all miners
status: 200 OK

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerGetMiners(w http.ResponseWriter, r *http.Request) {
	minerType := r.URL.Query().Get("class")

	var minerStr *miners.MinerType

	if minerType != "" {
		minerStr = (*miners.MinerType)(&minerType)
	}

	filtred := make(map[uuid.UUID]miners.MinerInfo, len(h.company.GetAllWorkedMiners()))

	for k, v := range h.company.GetAllWorkedMiners() {
		if minerStr != nil && *minerStr != v.MinerType {
			continue
		}

		filtred[k] = v
	}

	writeJSON(w, filtred)
}

/*
pattern: /miners/prices
method: GET
info: -

success:
response body: JSON body represent all miners prices
status: 200 OK

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerGetPrices(w http.ResponseWriter, r *http.Request) {

	b, err := json.MarshalIndent(miners.MinersSalaries, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("error response body", err.Error())
		return
	}
}

/*
pattern: /equipments
method: PATCH
info: JSON request body

success:
response body: JSON body represent bought equipment
status: 200 OK

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerBuyEquipment(w http.ResponseWriter, r *http.Request) {
	var equipmentDTO EquipmentDTO
	if err := json.NewDecoder(r.Body).Decode(&equipmentDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := equipmentDTO.ValidateForCreate(); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	equipInfo, err := h.company.BuyEquipment(equipmentDTO.Equip)
	if err != nil {
		errDTO := NewErrorDTO(err)
		if errors.Is(err, company.ErrorUnknownEquipment) || errors.Is(err, company.ErrorEquipBought) || errors.Is(err, company.ErrorNotEnoughCoal) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	writeJSON(w, equipInfo)
}

/*
pattern: /equipments/prices
method: GET
info: -

success:
response body: JSON body represent all prices equipment
status: 200 OK

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerEquipmentsPrices(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, equipment.Equipments)
}

/*
pattern: /equipments
method: GET
info: -

success:
response body: JSON body represent bought equipments
status: 200 OK

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerGetInfoEquipments(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.company.GetInfoEquipment())
}

/*
pattern: /company
method: GET
info: -

success:
response body: JSON body represent statistics game
status: 200 OK

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerStatisticsGame(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.company.Stats())
}

/*
pattern: /company/start
method: POST
info: -

success:
response body: JSON body represent finally statistics game
status: 200 OK

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerStartGame(w http.ResponseWriter, r *http.Request) {
	h.company.Run()

	if _, err := w.Write([]byte("Игра началась!")); err != nil {
		fmt.Println("error response HTTP body:", err.Error())
		return
	}
}

/*
pattern: /company/stop
method: POST
info: -

success:
response body: JSON body represent finally statistics game
status: 200 OK

failed:
response body: description error + with time
status: 400, 404, 409, 500, ...
*/
func (h *HTTPHandler) HandlerEndGame(w http.ResponseWriter, r *http.Request) {
	stats := h.company.Stop()

	writeJSON(w, stats)
}
