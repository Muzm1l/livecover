package symbols

import "translator/internal/model"

type Resolver interface {
	Resolve(exePath string, buildID string) ([]model.FunctionSymbol, error)
}
