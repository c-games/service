package service

import (
	"testing"
)

func TestLogger(t *testing.T) {

	l, _ := newLogger("logfile", "info")
	l.SetFormatter(FormatterTypeText, true)

	l.WithField(LevelInfo, "text", "1", "a", "b")
	l.WithField(LevelInfo, "", "", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "text", 1, "a", "b")
	_ = l.WithFieldFile(LevelInfo, "", "", "a", "b")

	l.SetFormatter(FormatterTypeJSON, true)
	l.WithField(LevelInfo, "json", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "json", 1, "a", "b")
	_ = l.WithFieldsFile(LevelInfo, Fields{"1": 1, "2": 3}, "a", "b")

}
