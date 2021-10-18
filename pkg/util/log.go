package util

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-kit/log/term"
)

func NewLogger() log.Logger {
	logger := term.NewLogger(os.Stdout, log.NewLogfmtLogger, colorFn)

	return logger
}

func colorFn(keyvals ...interface{}) term.FgBgColor {

	err := term.FgBgColor{Fg: term.Red}

	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(error); ok {
			return err
		}
	}

	for i := 0; i < len(keyvals)-1; i += 2 {
		if keyvals[i] != "level" {
			continue
		}
		switch keyvals[i+1] {
		case level.DebugValue():
			// return term.FgBgColor{Fg: term.DarkGray}
			return term.FgBgColor{Fg: term.Gray}
		case level.InfoValue():
			return term.FgBgColor{}
		case level.WarnValue():
			return term.FgBgColor{Fg: term.Yellow}
		case level.ErrorValue():
			return err
		default:
			fmt.Println("unknown")
			return term.FgBgColor{}
		}
	}
	return term.FgBgColor{}
}
