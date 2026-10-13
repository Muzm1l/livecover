package discovery

import (
	"context"

	"translator/internal/model"
)

type Watcher interface {
	Start(ctx context.Context) (<-chan model.ProcessEvent, error)
}
