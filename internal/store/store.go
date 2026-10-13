package store

import (
	"context"

	"translator/internal/model"
)

type Store interface {
	UpsertFunctions(ctx context.Context, functions []model.FunctionSymbol) error
	FlushBuckets(ctx context.Context, deltas []model.CoverageDelta) error
	UpdateLastSeen(ctx context.Context, seen []model.LastSeen) error
	QuerySnapshot(ctx context.Context, q model.QueryFilter) (model.Snapshot, error)
}
