package main

import (
	"company_coal_mine/company"
	"company_coal_mine/http"
	"context"
	"fmt"
	"sync"
)

var mtx sync.Mutex

func main() {
	company := company.NewCompany(context.Background())

	HTTPHandler := http.NewHTTPHandler(company)

	HTTPServer := http.NewHTTPServer(HTTPHandler)

	if err := HTTPServer.Start(); err != nil {
		fmt.Println("error to start game")
	}
}
