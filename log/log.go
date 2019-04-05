package loglib

import (
	"fmt"
	"os"

	logging "gopkg.in/op/go-logging.v1"
)

var Log *logging.Logger

func init() {
	Log = logging.MustGetLogger("")
	debugMode(false)
}

func LogInit(mod string, verbose bool) {
	logHeard := ""
	fmt.Sprintf(logHeard, "pscf-%s", mod)
	Log = logging.MustGetLogger(logHeard)
	debugMode(verbose)
}

func debugMode(verbose bool) {
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
	)
	var backend = logging.AddModuleLevel(
		logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

	if verbose {
		backend.SetLevel(logging.DEBUG, "")
	} else {
		backend.SetLevel(logging.ERROR, "")
	}

	logging.SetBackend(backend)
}
