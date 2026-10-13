package bpf

import (
	"context"

	"livecover/internal/model"
)

type Manager interface {
	Attach(ctx context.Context, spec model.ProbeSpec) error
	Detach(pid int) error
	Events() <-chan model.HitEvent
}
