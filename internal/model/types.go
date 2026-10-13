package model

import "time"

type ProcessInfo struct {
	PID       int
	ExePath   string
	Cmdline   []string
	Service   string
	BuildID   string
	Namespace string
}

type ProcessEventType string

const (
	ProcessStarted  ProcessEventType = "started"
	ProcessExited   ProcessEventType = "exited"
	ProcessReplaced ProcessEventType = "replaced"
)

type ProcessEvent struct {
	Type    ProcessEventType
	Process ProcessInfo
}

type FunctionSymbol struct {
	FuncID     uint64
	Name       string
	Package    string
	File       string
	Line       int
	Address    uint64
	Confidence string
}

type ProbeSpec struct {
	PID       int
	BuildID   string
	Service   string
	Funcs     []FunctionSymbol
	MaxProbes int
}

type HitEvent struct {
	TsUnixNano uint64
	PID        uint32
	FuncID     uint64
	BuildID    string
	Service    string
	CPU        uint16
	NodeID     string
}

type Batch struct {
	NodeID string
	Hits   []HitEvent
}

type CoverageDelta struct {
	Service    string
	BuildID    string
	FuncID     uint64
	BucketUnix int64
	Hits       uint32
}

type LastSeen struct {
	Service    string
	BuildID    string
	FuncID     uint64
	LastSeenTs time.Time
	TotalHits  uint64
}

type QueryFilter struct {
	Service string
	BuildID string
	Window  time.Duration
}

type Snapshot struct {
	GeneratedAt time.Time
	Service     string
	BuildID     string
	Deltas      []CoverageDelta
	LastSeen    []LastSeen
}
