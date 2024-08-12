package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[1;31m"
	Green       = "\033[1;32m"
	Yellow      = "\033[1;33m"
	Blue        = "\033[1;34m"
	Magenta     = "\033[1;35m"
	Cyan        = "\033[1;36m"
	White       = "\033[1;37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

func myFormat(pre, levelColor, text, reset string) string {
	timeHead := text[:23]
	levelBox := text[24:27]
	index := strings.Index(text[28:], "]")
	// 使用第一个 "]" 位置切割字符串
	caller := text[28:][:index+1]
	msg := strings.TrimSpace(text[28:][index+1:])
	levelColor = "\033[" + levelColor + "m"
	return fmt.Sprintf("%s%s%s %s%s%s %s%s%s %s%s%s\n",
		Cyan, timeHead, Reset,
		levelColor, levelBox, Reset,
		Green, caller, Reset,
		levelColor, msg, Reset)
}

// brush is a color join function
type brush func(string) string

// newBrush returns a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return myFormat(pre, color, text, reset)
	}
}

func emptyBrush(text string) string {
	return text
}

var colors = []brush{
	newBrush("0"),    // Trace              No Color
	newBrush("1;36"), // Debug              Light Cyan
	newBrush("1;34"), // Info 				Blue
	newBrush("1;33"), // Warn               Yellow
	newBrush("1;31"), // Error              Red
}

func colorBrushByLevel(level Level) brush {
	switch level {
	case TraceLevel:
		return colors[0]
	case DebugLevel:
		return colors[1]
	case InfoLevel:
		return colors[2]
	case WarnLevel:
		return colors[3]
	case ErrorLevel:
		return colors[4]
	default:
		return colors[2]
	}
}

var _ io.Writer = (*ConsoleWriter)(nil)

type ConsoleConfig struct {
	Colorful bool
}

type ConsoleWriter struct {
	w io.Writer
}

func NewConsoleWriter(cfg ConsoleConfig, out io.Writer) io.Writer {
	if !cfg.Colorful {
		return out
	}
	if out == nil {
		out = os.Stdout
	}
	return &ConsoleWriter{
		w: out,
	}
}

func (cw *ConsoleWriter) Write(p []byte) (n int, err error) {
	return cw.w.Write(p)
}

func (cw *ConsoleWriter) WriteLog(p []byte, level Level, when time.Time) (n int, err error) {
	return cw.w.Write([]byte(colorBrushByLevel(level)(string(p))))
}
