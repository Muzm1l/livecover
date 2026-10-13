package symbols

import "livecover/internal/model"

type Resolver interface {
	Resolve(exePath string, buildID string) ([]model.FunctionSymbol, error)
}
