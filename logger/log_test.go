package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	l := getLog("test", true)

	//
	l.Info("just test")
	l.Warn("just test")
	l.Error("just test")

}
