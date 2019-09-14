package netpro

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func initViper() {
	// 确保viper最先被init

	// 设置默认值
	viper.SetDefault("runtime.goroutine_num", 30000)

	// 读取命令行参数
	workdir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	configPath := pflag.String("config_path", workdir, "Specify the path of the configuration file")
	configName := pflag.String("config_name", "config", "Configuration file names that do not contain suffix names. (Support JSON, TOML, YAML, HCL, envfile or Java properties format)")
	pflag.Parse()

	// 读取配置
	viper.SetConfigName(*configName)
	viper.AddConfigPath(*configPath)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("viper read in config error:", err)
	}
	logrus.Debugf("loaded %s in %s\n", *configName, *configPath)
}
