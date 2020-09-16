package log

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Log a logger instance
var Log = log.New()

func init() {
	var err error
	var logFile *os.File
	var logLevel = log.InfoLevel

	fn := "/var/log/pigeond.log"
	for _, arg := range os.Args {
		if arg == "--debug" {
			logLevel = log.DebugLevel
		}
	}

	if logFile, err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	Log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	Log.SetLevel(logLevel)
}
