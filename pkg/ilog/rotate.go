package ilog

import (
	"io"

	"github.com/xqk/good/pkg/ilog/rotate"
)

func newRotate(config *Config) io.Writer {
	rotateLog := rotate.NewLogger()
	rotateLog.Filename = config.Filename()
	rotateLog.MaxSize = config.MaxSize // MB
	rotateLog.MaxAge = config.MaxAge   // days
	rotateLog.MaxBackups = config.MaxBackup
	rotateLog.Interval = config.Interval
	rotateLog.LocalTime = true
	rotateLog.Compress = false
	return rotateLog
}
