package deadcode

import "translator/internal/model"

type Candidate struct {
	Service    string
	BuildID    string
	FuncID     uint64
	Score      float64
	Confidence string
	Reason     string
}

type Scorer interface {
	Evaluate(snapshot model.Snapshot) []Candidate
}
