package ilog_test

import (
	"testing"

	"github.com/xqk/good/pkg/ilog"
)

func Test_Info(t *testing.T) {
	ilog.Info("hello", ilog.Any("a", "b"))
}
