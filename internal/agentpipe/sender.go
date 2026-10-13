package agentpipe

import (
	"context"

	"translator/internal/model"
)

type Sender interface {
	Send(ctx context.Context, batch model.Batch) error
}
