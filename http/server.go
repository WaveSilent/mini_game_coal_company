package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	handler *HTTPHandler
}

func NewHTTPServer(httpHandler *HTTPHandler) *HTTPServer {
	return &HTTPServer{
		handler: httpHandler,
	}
}

func (s *HTTPServer) Start() error {
	router := mux.NewRouter()

	router.Path("/miners").Methods("POST").HandlerFunc(s.handler.HandlerHireMiner)
	router.Path("/miners").Methods("GET").HandlerFunc(s.handler.HandlerGetMiners)
	router.Path("/miners/prices").Methods("GET").HandlerFunc(s.handler.HandlerGetPrices)

	router.Path("/equipments").Methods("PATCH").HandlerFunc(s.handler.HandlerBuyEquipment)
	router.Path("/equipments/prices").Methods("GET").HandlerFunc(s.handler.HandlerEquipmentsPrices)
	router.Path("/equipments").Methods("GET").HandlerFunc(s.handler.HandlerGetInfoEquipments)

	router.Path("/company").Methods("GET").HandlerFunc(s.handler.HandlerStatisticsGame)
	router.Path("/company/start").Methods("POST").HandlerFunc(s.handler.HandlerStartGame)
	router.Path("/company/stop").Methods("POST").HandlerFunc(s.handler.HandlerEndGame)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}
