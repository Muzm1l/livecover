package agg

import (
	"context"

	"translator/internal/model"
)

type Engine interface {
	IngestHits(ctx context.Context, hits []model.HitEvent) error
	Snapshot(filter model.QueryFilter) model.Snapshot
	Subscribe(filter model.QueryFilter) (<-chan model.CoverageDelta, func())
}
