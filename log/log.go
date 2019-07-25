package log

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
)

const LOG_FILE = "./go-proxy.log"

var log = new(logrus.Logger)

func Info(any interface{}) {
	log.Info(any)
}

func Warn(any interface{}) {
	log.Warn(any)
}

func Error(any interface{}) {
	log.Error(any)
}

func Fatal(any interface{}) {
	log.Fatal(any)
}


func Panic(any interface{}) {
	log.Panic(any)
}

func init() {
	log.Formatter = &prefixed.TextFormatter{
		TimestampFormat : "2006-01-02 15:04:05",
		FullTimestamp:true,
		ForceFormatting: true,
	}
	file, _ := os.OpenFile(LOG_FILE, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetLevel(logrus.DebugLevel)
	log.WithFields(logrus.Fields{
		"prefix": "main",
		"animal": "walrus",
		"number": 8,
	}).Debug("Started observing beach")

	log.WithFields(logrus.Fields{
		"prefix":      "sensor",
		"temperature": -4,
	}).Info("Temperature changes")
}
