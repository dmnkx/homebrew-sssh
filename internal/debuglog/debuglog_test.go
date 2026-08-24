package debuglog

import "testing"

func TestLogDoesNotPanic(t *testing.T) {
	Log("H0", "debuglog_test.go", "ok", map[string]any{"n": 1})
	Log("H0", "debuglog_test.go", "nil data", nil)
}

func TestLogSkipsUnmarshalableData(t *testing.T) {
	Log("H0", "debuglog_test.go", "bad", map[string]any{"fn": func() {}})
}
