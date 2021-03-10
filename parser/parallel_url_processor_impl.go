package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/uveaceslav/affise/model"
)

var (
	ErrDuringProcessing = errors.New("Error during processing")
)

type parallelURLProcessor struct {
	urlFetcher     URLFetcher
	workerPoolSize int
}

func NewParallelURLProcessor(urlFetcher URLFetcher, workerPoolSize int) URLProcessor {
	return &parallelURLProcessor{
		urlFetcher:     urlFetcher,
		workerPoolSize: workerPoolSize,
	}
}

func (up *parallelURLProcessor) Process(ctx context.Context, cancel context.CancelFunc, urls []string) ([]model.ProcessedEntity, error) {
	reqChan := make(chan string)
	respChan := make(chan model.ProcessedEntity)

	// producer
	go func(ctx context.Context, cancel context.CancelFunc, reqChan chan<- string, urls []string) {
		for _, uri := range urls {
			reqChan <- uri
		}
		close(reqChan)
	}(ctx, cancel, reqChan, urls)

	// worker pool
	go func(
		ctx context.Context,
		cancel context.CancelFunc,
		reqChan chan string,
		respChan chan model.ProcessedEntity,
	) {
		var wg sync.WaitGroup
		for i := 0; i < up.workerPoolSize; i++ {
			wg.Add(1)
			go up.worker(ctx, cancel, reqChan, respChan, &wg)
		}
		wg.Wait()
		close(respChan)
	}(ctx, cancel, reqChan, respChan)

	// consumer
	result := make([]model.ProcessedEntity, 0, 20)
	for entity := range respChan {
		result = append(result, entity)
	}

	if len(result) < len(urls) {
		return nil, ErrDuringProcessing
	}

	return result, nil
}

func (up *parallelURLProcessor) worker(
	ctx context.Context,
	cancel context.CancelFunc,
	reqChan chan string,
	respChan chan model.ProcessedEntity,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for uri := range reqChan {
		// Check if any error occurred in any other gorouties:
		select {
		case <-ctx.Done():
			return // Error somewhere, terminate
		default: // Default is must to avoid blocking
		}

		body, err := up.urlFetcher.Fetch(uri)
		if err != nil {
			cancel()
			return
		}

		var data interface{}
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&data); err != nil {
			cancel()
			return
		}

		respChan <- model.ProcessedEntity{
			URL:  uri,
			Data: data,
		}
	}
}
