package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/uveaceslav/affise/model"
)

var (
	ErrTerminated = errors.New("Process has been terminated")
)

type urlProcessor struct {
	urlFetcher URLFetcher
}

func NewURLProcessor(urlFetcher URLFetcher) URLProcessor {
	return &urlProcessor{
		urlFetcher: urlFetcher,
	}
}

func (up *urlProcessor) Process(ctx context.Context, _ context.CancelFunc, urls []string) ([]model.ProcessedEntity, error) {
	result := make([]model.ProcessedEntity, 0, 20)
	for _, uri := range urls {
		select {
		case <-ctx.Done():
			return nil, ErrTerminated
		default: // Default is must to avoid blocking
		}

		body, err := up.urlFetcher.Fetch(uri)
		if err != nil {
			return nil, err
		}

		var data interface{}
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&data); err != nil {
			return nil, err
		}

		result = append(result, model.ProcessedEntity{
			URL:  uri,
			Data: data,
		})
	}

	return result, nil
}
