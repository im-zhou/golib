package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Caller path involves line numbers, when the line number for logging in this function changes,
// it needs to be synchronized.
func TestLogWithCaller(t *testing.T) {
	require := require.New(t)

	testBuffer := bytes.NewBuffer(nil)
	logger := New(
		WithLevel(DebugLevel),
		WithCaller(true),
		WithOutput(testBuffer),
	)
	logger.Debug("test")

	require.Contains(testBuffer.String(), "log/log_test.go:24")
}

func testLogFormat(t time.Time, level Level, msg string) string {
	return fmt.Sprintf("%s %s %s\n", t.Format(logTimeFormat), level.LogPrefix(), msg)
}

func BenchmarkLog(b *testing.B) {
	logger := New(
		WithLevel(DebugLevel),
		WithCaller(true),
		WithOutput(io.Discard),
	)
	for i := 0; i < b.N; i++ {
		logger.Debugf("debug %d", i)
	}
}

func TestLogOffset(t *testing.T) {
	require := require.New(t)
	testBuffer := bytes.NewBuffer(nil)
	logger := New(
		WithLevel(DebugLevel),
		WithCaller(true),
		WithOutput(testBuffer),
	)
	wrapLog := func() {
		logger.Log(InfoLevel, 1, "test")
	}
	wrapLog()

	require.Contains(testBuffer.String(), "log/log_test.go:55")
}

func TestLogOutput(t *testing.T) {
	InitLogger("trace")

	LLogger.Infof("this is info")
	LLogger.Warnf("this is warn")
	LLogger.Errorf("this is error")
	LLogger.Tracef("this is trace")
	LLogger.Debugf("this is debug")
}

var LLogger *Logger

func init() {
	LLogger = New(
		WithCaller(true),
		AddCallerSkip(0),
		WithLevel(InfoLevel),
	)
}

func InitLogger(levelStr string) {
	var options []Option

	options = append(options,
		WithOutput(NewConsoleWriter(ConsoleConfig{
			Colorful: true,
		}, os.Stdout)))

	level, err := ParseLevel(levelStr)
	if err != nil {
		level = InfoLevel
	}
	options = append(options, WithLevel(level))
	LLogger = LLogger.WithOptions(options...)
}
