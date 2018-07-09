package logger

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
)

var info *log.Logger
var debug *log.Logger
var warn *log.Logger
var err *log.Logger

// StubLogger sends all logs to /dev/null
func StubLogger() {
	info = log.New(ioutil.Discard, "", 0)
	debug = log.New(ioutil.Discard, "", 0)
	warn = log.New(ioutil.Discard, "", 0)
	err = log.New(ioutil.Discard, "", 0)
}

// TestLogger applies a test logger to allow fetching of logoutput during tests
func TestLogger() *LogOutput {
	logOutput := &LogOutput{}
	info = log.New(logOutput, "INFO: ", 0)
	debug = log.New(logOutput, "DEBUG: ", 0)
	warn = log.New(logOutput, "WARN: ", 0)
	err = log.New(logOutput, "ERR: ", 0)
	return logOutput
}

// LogOutput store the log output as an array of lines
type LogOutput struct {
	Lines []string
}

func (logoutput *LogOutput) Write(p []byte) (int, error) {
	logoutput.Lines = append(logoutput.Lines, string(p))
	return len(p), nil
}

// Initialize loggers for each level
func Initialize() {
	info = log.New(os.Stdout, color.BlueString("[INFO] "), 0)
	debug = log.New(os.Stdout, color.GreenString("[DEBUG] "), 0)
	warn = log.New(os.Stdout, color.YellowString("[WARN] "), 0)
	err = log.New(os.Stderr, color.RedString("[ERROR] "), 0)
}

// Info log to info level
func Info(line string, params ...interface{}) {
	info.Printf(line+"\n", params...)
}

// Debug log to debug level
func Debug(line string, params ...interface{}) {
	debug.Printf(line+"\n", params...)
}

// Warn log to warn level
func Warn(line string, params ...interface{}) {
	warn.Printf(line+"\n", params...)
}

// Error log to error level
func Error(line string, params ...interface{}) {
	err.Printf(line+"\n", params...)
}

func ErrorAndExit(line string, params ...interface{}) {
	err.Printf(line+"\n", params...)
	os.Exit(1)
}
