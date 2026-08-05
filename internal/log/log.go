//go:build windows

// Package log provides leveled logging: timestamp + level tag, written to a
// pluggable sink (the GUI injects a sink that appends to the log box).
package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var (
	logMu        sync.Mutex
	sink         io.Writer
	fatalHandler func(text string)
)

// SetSink sets the destination for log lines (e.g. the UI log box).
func SetSink(w io.Writer) {
	logMu.Lock()
	sink = w
	logMu.Unlock()
}

// SetFatalHandler sets the callback invoked by Fatal (e.g. show a message box).
func SetFatalHandler(h func(text string)) {
	logMu.Lock()
	fatalHandler = h
	logMu.Unlock()
}

func line(level, format string, args ...interface{}) {
	logMu.Lock()
	defer logMu.Unlock()
	msg := fmt.Sprintf(format, args...)
	l := time.Now().Format("2006-01-02 15:04:05") + "  " + levelPad(level) + " " + msg + "\r\n"
	if sink != nil {
		sink.Write([]byte(l))
	}
}

// levelPad left-pads the level tag to 7 columns so lines align.
func levelPad(l string) string {
	for len(l) < 7 {
		l += " "
	}
	return l
}

// Info logs general information (INFO).
func Info(format string, args ...interface{}) { line("INFO", format, args...) }

// Event logs system/hardware events (EVENT, e.g. lid close/open, wake).
func Event(format string, args ...interface{}) { line("EVENT", format, args...) }

// Action logs program actions (ACTION, e.g. scheduling/executing a sleep).
func Action(format string, args ...interface{}) { line("ACTION", format, args...) }

// Error logs an error (ERROR). Long messages may use \n; put details on a
// "Reason:" line.
func Error(format string, args ...interface{}) { line("ERROR", format, args...) }

// Print logs an INFO line, preserving the caller's original formatting.
func Print(format string, args ...interface{}) { line("INFO", format, args...) }

// Fatal runs the fatal handler (if set), otherwise prints to stderr, then exits.
func Fatal(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	logMu.Lock()
	h := fatalHandler
	logMu.Unlock()
	if h != nil {
		h(text)
	} else {
		fmt.Fprintln(os.Stderr, text)
	}
	os.Exit(1)
}
