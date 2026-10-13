package agentpipe

import (
	"context"

	"livecover/internal/model"
)

type Sender interface {
	Send(ctx context.Context, batch model.Batch) error
}
