package logger

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var info *log.Logger
var debug *log.Logger
var warn *log.Logger
var err *log.Logger

func Initialize() {
	info = log.New(os.Stdout, color.BlueString("[INFO] "), 0)
	debug = log.New(os.Stdout, color.GreenString("[DEBUG] "), 0)
	warn = log.New(os.Stdout, color.YellowString("[WARN] "), 0)
	err = log.New(os.Stderr, color.RedString("[ERROR] "), 0)
}

func Info(line string, params ...interface{}) {
	info.Printf(line+"\n", params...)
}

func Debug(line string, params ...interface{}) {
	debug.Printf(line+"\n", params...)
}

func Warn(line string, params ...interface{}) {
	warn.Printf(line+"\n", params...)
}

func Error(line string, params ...interface{}) {
	err.Printf(line+"\n", params...)
}

func ErrorAndExit(line string, params ...interface{}) {
	err.Printf(line+"\n", params...)
	os.Exit(1)
}
