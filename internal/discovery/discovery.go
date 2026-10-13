package discovery

import (
	"context"

	"livecover/internal/model"
)

type Watcher interface {
	Start(ctx context.Context) (<-chan model.ProcessEvent, error)
}
