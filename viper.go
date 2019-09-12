package netpro

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initViper() {
	// 确保viper最先被init
	viper.SetConfigName("config")
	viper.SetDefault("runtime.goroutine_num", 30000)
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Error("get work dir error:", err)
	}
	viper.AddConfigPath(cwd)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Error("viper read in config error:", err)
	}
}
