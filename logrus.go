package netpro

import (
	"os"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Logger logrus entry
var Logger *logrus.Entry

func initLogrus() {
	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		logrus.Error("logrus parse level error:", err)
	} else {
		logrus.SetLevel(level)
	}

	if strings.ToLower(viper.GetString("log.format")) == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	if strings.ToLower(viper.GetString("log.output")) == "file" {
		file, err := os.OpenFile(viper.GetString("log.filename"), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			logrus.Fatal("open log file error:", err)
		}
		logrus.SetOutput(file)
	} else {
		logrus.SetOutput(os.Stdout)
	}
	Logger = logrus.WithFields(logrus.Fields{
		"pid": syscall.Getpid(),
	})
}
