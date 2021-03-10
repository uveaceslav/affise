package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/uveaceslav/affise/controller"
	"github.com/uveaceslav/affise/middleware"
	"github.com/uveaceslav/affise/model/validator"
	"github.com/uveaceslav/affise/parser"
	"github.com/uveaceslav/affise/server"
)

//
// MOCKs
//

type mockURLFetcher struct {
	timeout int // seconds
}

func NewMockURLFetcher(timeout int) parser.URLFetcher {
	return &mockURLFetcher{
		timeout: timeout,
	}
}

func (muf *mockURLFetcher) Fetch(uri string) ([]byte, error) {
	if muf.timeout > 0 { // sleep a timeout * second per job to simulate an expensive task.
		time.Sleep(time.Duration(muf.timeout) * time.Second)
	}

	response := map[string]int64{
		"timestamp": time.Now().Unix(),
	}

	return json.Marshal(response)
}

//
// Test
//

func TestRateLimit(t *testing.T) {
	host := "127.0.0.1" // ToDo: Move to config
	port := 8080        // ToDo: Move to config
	maxRate := 2        // ToDo: Move to config
	workerPoolSize := 4 // ToDo: Move to config

	// Start server
	go func() {
		logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)
		urlFetcher := NewMockURLFetcher(5)
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
	}()

	time.Sleep(1 * time.Second)

	// Start client
	respChan := make(chan int)
	totalRequestCount := 15

	go func() {
		var wg sync.WaitGroup
		for i := 0; i < totalRequestCount; i++ {
			wg.Add(1)
			go func(i int, respChan chan int, wg *sync.WaitGroup) {
				defer wg.Done()

				resp, _ := http.Post(fmt.Sprintf("http://%s:%d", host, port), "application/json", createRequestBody(1))
				respChan <- resp.StatusCode
			}(i, respChan, &wg)
		}

		wg.Wait()

		close(respChan)
	}()

	counter := 0
	for statusCode := range respChan {
		if statusCode == http.StatusTooManyRequests {
			counter++
		}
	}

	if counter != (totalRequestCount - maxRate) {
		t.Errorf("expected: <%d> but was: <%d>", (totalRequestCount - maxRate), counter)
	}
}

func createRequestBody(countURLs int) io.Reader {
	urls := make([]string, 0, countURLs)
	for i := 0; i < countURLs; i++ {
		urls = append(urls, fmt.Sprintf("http://localhost/%d", i))
	}

	data := map[string]interface{}{"urls": urls}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return bytes.NewBuffer(jsonData)
}
