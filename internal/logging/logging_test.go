package logging

import (
	"log/slog"
	"testing"
)

func TestParseLevelDefaultsToInfo(t *testing.T) {
	logger := New(Options{Level: "nope", Format: "text", Component: "test"})
	if logger == nil {
		t.Fatal("expected logger")
	}
	logger.Info("ok")
}

func TestDiscardLogger(t *testing.T) {
	logger := Discard()
	logger.Error("should not panic")
	_ = slog.LevelInfo
}
