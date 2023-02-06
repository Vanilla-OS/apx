package log

import (
	"log"
	l "log"
	"strings"
)

const DefaultPrefix string = ""
const DefaultFlags int = log.Ldate

var AppName string = "orchid"

// Prefix configures std log package with a prefix
func Prefix(prefix string) {
	l.SetPrefix(prefix)
	AppName = strings.TrimSpace(prefix)
}

// Flags configures std log package with features
// like date, time, and file name
// see https://pkg.go.dev/log#SetFlags for options
func Flags(flags int) {
	l.SetFlags(flags)
}
