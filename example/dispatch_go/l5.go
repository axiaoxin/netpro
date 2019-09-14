package main

import (
	"log"
	"os"
	"strings"
	"time"

	"selfgit/components/l5"
	"github.com/axiaoxin/netpro"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var L5Client *l5.Api

func InitL5() {
	// 调整MaxPacketSize为默认值的2倍避免bad packet错误
	option := l5.Option{
		Host:                    "127.0.0.1",                                                                               //默认agent ip
		Port:                    8888,                                                                                      //默认agent port
		Timeout:                 time.Second,                                                                               //操作超时
		MaxPacketSize:           20480,                                                                                     //agent通信允许最大包
		MaxConn:                 5,                                                                                         //与agent最大连接数
		StaticNameFiles:         []string{"/data/L5Backup/name2sid.backup", "/data/L5Backup/name2sid.cache.bin"},           //默认domain静态文件
		StaticNameReload:        30 * time.Second,                                                                          //静态domain重载时间
		StaticRouteFiles:        []string{"/data/L5Backup/current_route.backup", "/data/L5Backup/current_route_v2.backup"}, //默认server静态文件
		StatErrorReportInterval: time.Second,                                                                               //错误上报间隔
		StatReportInterval:      5 * time.Second,                                                                           //正常上报间隔
		StatMaxErrorCount:       16,                                                                                        //最大错误数
		StatMaxErrorRate:        0.2,                                                                                       //最大错误比例
		BalancerFunc:            l5.NewWeightedRoundRobinBalancer,
		Logger:                  log.New(os.Stdout, "[l5]", log.Ldate|log.Ltime),
	}
	var err error
	L5Client, err = l5.NewApi(&option)
	if err != nil {
		logrus.Fatal("init L5Client NewApi error:", err)
	}
}

func GetLogaccServer() (l5.Server, error) {
	var mod int32
	var cmd int32
	if strings.ToLower(viper.GetString("env")) == "dev" {
		mod = viper.GetInt32("l5.logacc.dev.mod")
		cmd = viper.GetInt32("l5.logacc.dev.cmd")
	}
	if strings.ToLower(viper.GetString("env")) == "default_cluster" {
		mod = viper.GetInt32("l5.logacc.default_cluster.mod")
		cmd = viper.GetInt32("l5.logacc.default_cluster.cmd")
	}
	netpro.Logger.Debugf("logacc l5 sid = %d:%d", mod, cmd)
	srv, err := L5Client.GetServerBySid(mod, cmd)
	return srv, errors.Wrap(err, "l5 get server by sid error")
}
