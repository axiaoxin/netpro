package netpro

import (
	"os"
	"path/filepath"
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
		logFilename := viper.GetString("log.filename")
		logFilepath, err := filepath.Abs(filepath.Dir(logFilename))
		if err != nil {
			logrus.Fatal("get logfilename abspath error:", err)
		}
		if err := os.MkdirAll(logFilepath, os.ModeDir); err != nil {
			logrus.Fatal("mkdir all for logfilepath error:", err)
		}
		file, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
