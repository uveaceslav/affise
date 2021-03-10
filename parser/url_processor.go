package parser

import (
	"context"

	"github.com/uveaceslav/affise/model"
)

type URLProcessor interface {
	Process(ctx context.Context, cancel context.CancelFunc, urls []string) ([]model.ProcessedEntity, error)
}
