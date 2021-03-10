package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/uveaceslav/affise/controller"
	"github.com/uveaceslav/affise/middleware"
	"github.com/uveaceslav/affise/model/validator"
	"github.com/uveaceslav/affise/parser"
	"github.com/uveaceslav/affise/server"
)

func main() {
	host := "0.0.0.0"                  // ToDo: Move to config
	port := 8080                       // ToDo: Move to config
	maxRate := 100                     // ToDo: Move to config
	workerPoolSize := 4                // ToDo: Move to config
	urlFetchTimeout := 1 * time.Second // ToDo: Move to config

	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	httpClient := &http.Client{
		Timeout: urlFetchTimeout,
	}

	urlFetcher := parser.NewURLFetcher(httpClient)
	// urlProcessor := parser.NewURLProcessor(urlFetcher)
	urlProcessor := parser.NewParallelURLProcessor(urlFetcher, workerPoolSize)
	requestEntityValidator := validator.NewRequestEntityValidator()
	mainController := controller.NewMainController(logger, requestEntityValidator, urlProcessor)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(maxRate)

	// Setup router
	router := http.NewServeMux()
	router.HandleFunc("/", rateLimitMiddleware.Limit(mainController.Index))

	app := server.NewServer(logger, router)
	app.Serve(fmt.Sprintf("%s:%d", host, port)) // config.App.ServerAddress
}
